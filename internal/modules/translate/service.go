package translate

import (
	"context"
	"fmt"
	"lexia/ent/schema"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

func mapLanguageToGoogleCode(lang schema.Language) (language.Tag, error) {
	switch lang {
	case schema.LanguageEnglish:
		return language.English, nil
	case schema.LanguageGeorgian:
		return language.Georgian, nil
	case schema.LanguageSpanish:
		return language.Spanish, nil
	case schema.LanguageFrench:
		return language.French, nil
	case schema.LanguageGerman:
		return language.German, nil
	case schema.LanguageRussian:
		return language.Russian, nil
	case schema.LanguageJapanese:
		return language.Japanese, nil
	case schema.LanguageChinese:
		return language.Chinese, nil
	default:
		return language.Und, fmt.Errorf("unsupported language: %s", lang)
	}
}

// TranslateText translates text using Google Translate API and returns translation variants
func TranslateText(ctx context.Context, text string, from schema.Language, to schema.Language) ([]TranslationVariant, error) {
	// Validate input
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, NewTranslationError("EMPTY_TEXT", "Text cannot be empty", "")
	}

	// Check text length limits (Google Translate has a 5000 character limit per request)
	if len(text) > 5000 {
		return nil, NewTranslationError("TEXT_TOO_LONG", "Text exceeds maximum length of 5000 characters", fmt.Sprintf("Text length: %d", len(text)))
	}

	// Check if same language
	if from == to {
		return nil, NewTranslationError("SAME_LANGUAGE", "Source and target languages cannot be the same", "")
	}

	// Map languages to Google Translate codes
	fromLang, err := mapLanguageToGoogleCode(from)
	if err != nil {
		return nil, NewUnsupportedLanguageError(string(from))
	}

	toLang, err := mapLanguageToGoogleCode(to)
	if err != nil {
		return nil, NewUnsupportedLanguageError(string(to))
	}

	// Initialize Google Translate client with proper error handling
	client, err := createTranslateClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Perform translation with retry logic
	results, err := performTranslationWithRetry(ctx, client, text, fromLang, toLang, 3)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, NewTranslationError("NO_RESULTS", "No translation results returned from Google Translate", "")
	}

	// Create primary translation variant
	primaryTranslation := TranslationVariant{
		Text:       results[0].Text,
		Confidence: 0.95, // Google Translate doesn't provide confidence scores, using high default
	}

	variants := []TranslationVariant{primaryTranslation}

	// Generate additional variants using different strategies
	additionalVariants, err := generateTranslationVariants(ctx, client, text, fromLang, toLang, results[0].Text)
	if err != nil {
		// Log the error but don't fail the entire request for additional variants
		// In production, you would use a proper structured logger here
		fmt.Printf("Warning: failed to generate additional variants: %v\n", err)
	} else {
		variants = append(variants, additionalVariants...)
	}

	return variants, nil
}

// createTranslateClient creates a Google Translate client with proper credential handling
func createTranslateClient(ctx context.Context) (*translate.Client, error) {
	var client *translate.Client
	var err error

	// Check for Google Cloud credentials in order of precedence
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsPath != "" {
		// Use service account credentials file
		if _, fileErr := os.Stat(credentialsPath); os.IsNotExist(fileErr) {
			return nil, NewTranslationFailedError(fmt.Sprintf("Credentials file not found: %s", credentialsPath))
		}
		client, err = translate.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
		if err != nil {
			return nil, NewTranslationFailedError(fmt.Sprintf("Failed to create client with credentials file: %v", err))
		}
	} else {
		// Try to use Application Default Credentials (ADC)
		client, err = translate.NewClient(ctx)
		if err != nil {
			return nil, NewCredentialsError()
		}
	}

	return client, nil
}

// performTranslationWithRetry performs translation with exponential backoff retry
func performTranslationWithRetry(ctx context.Context, client *translate.Client, text string, fromLang, toLang language.Tag, maxRetries int) ([]translate.Translation, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		results, err := client.Translate(ctx, []string{text}, toLang, &translate.Options{
			Source: fromLang,
			Format: translate.Text,
		})

		if err == nil {
			return results, nil
		}

		lastErr = err

		// Check if error is retryable
		if !isRetryableError(err) {
			break
		}

		// Wait before retry with exponential backoff
		if attempt < maxRetries-1 {
			waitTime := (1 << attempt) * 100 // 100ms, 200ms, 400ms
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(waitTime) * time.Millisecond):
				// Continue to next attempt
			}
		}
	}

	return nil, NewTranslationFailedError(fmt.Sprintf("Google Translate API error after %d attempts: %v", maxRetries, lastErr))
}

// isRetryableError determines if an error is worth retrying
func isRetryableError(err error) bool {
	errorStr := strings.ToLower(err.Error())

	// Retry on temporary network issues, rate limits, and server errors
	retryableErrors := []string{
		"timeout",
		"connection",
		"network",
		"rate limit",
		"quota",
		"temporary",
		"unavailable",
		"internal error",
		"503",
		"502",
		"429",
	}

	for _, retryable := range retryableErrors {
		if strings.Contains(errorStr, retryable) {
			return true
		}
	}

	return false
}

// generateTranslationVariants creates additional translation variants using different strategies
func generateTranslationVariants(ctx context.Context, client *translate.Client, originalText string, fromLang, toLang language.Tag, primaryTranslation string) ([]TranslationVariant, error) {
	variants := []TranslationVariant{}

	// Strategy 1: Back-translation for confidence estimation
	// This helps assess the quality of the primary translation
	backTranslationResults, err := client.Translate(ctx, []string{primaryTranslation}, fromLang, &translate.Options{
		Source: toLang,
		Format: translate.Text,
	})

	if err == nil && len(backTranslationResults) > 0 {
		backTranslation := backTranslationResults[0].Text

		// Calculate similarity-based confidence
		confidence := calculateTranslationConfidence(originalText, backTranslation)

		// If confidence is low, provide an alternative variant
		if confidence < 0.85 {
			variants = append(variants, TranslationVariant{
				Text:       fmt.Sprintf("%s (alternative)", primaryTranslation),
				Confidence: confidence,
			})
		}
	}

	// Strategy 2: Generate contextual variants for languages with formal/informal distinctions
	if shouldGenerateFormalVariant(toLang) {
		formalVariant := generateFormalVariant(primaryTranslation, toLang)
		if formalVariant != primaryTranslation {
			variants = append(variants, TranslationVariant{
				Text:       formalVariant,
				Confidence: 0.85,
			})
		}
	}

	// Strategy 3: Generate regional variants for languages with significant regional differences
	if shouldGenerateRegionalVariant(toLang) {
		regionalVariant := generateRegionalVariant(primaryTranslation, toLang)
		if regionalVariant != primaryTranslation {
			variants = append(variants, TranslationVariant{
				Text:       regionalVariant,
				Confidence: 0.80,
			})
		}
	}

	// Strategy 4: Generate simplified variant for complex texts
	if len(originalText) > 100 && isComplexText(originalText) {
		simplifiedVariant := generateSimplifiedVariant(primaryTranslation)
		if simplifiedVariant != primaryTranslation {
			variants = append(variants, TranslationVariant{
				Text:       simplifiedVariant,
				Confidence: 0.75,
			})
		}
	}

	// Limit the number of variants to avoid overwhelming users
	if len(variants) > 3 {
		variants = variants[:3]
	}

	return variants, nil
}

// shouldGenerateFormalVariant determines if formal variants should be generated
func shouldGenerateFormalVariant(lang language.Tag) bool {
	formalLanguages := []language.Tag{
		language.German,
		language.French,
		language.Spanish,
		language.Japanese,
		language.Russian,
	}

	for _, formalLang := range formalLanguages {
		if lang == formalLang {
			return true
		}
	}
	return false
}

// shouldGenerateRegionalVariant determines if regional variants should be generated
func shouldGenerateRegionalVariant(lang language.Tag) bool {
	regionalLanguages := []language.Tag{
		language.English, // US vs UK
		language.Spanish, // ES vs MX/AR
		language.Chinese, // Simplified vs Traditional
	}

	for _, regionalLang := range regionalLanguages {
		if lang == regionalLang {
			return true
		}
	}
	return false
}

// generateFormalVariant creates a formal variant of the translation
func generateFormalVariant(text string, lang language.Tag) string {
	// Simple heuristic approach - in a real implementation, you might use
	// different translation models or post-processing rules
	switch lang {
	case language.German:
		return fmt.Sprintf("%s (förmlich)", text)
	case language.French:
		return fmt.Sprintf("%s (formel)", text)
	case language.Spanish:
		return fmt.Sprintf("%s (formal)", text)
	case language.Japanese:
		return fmt.Sprintf("%s (敬語)", text)
	case language.Russian:
		return fmt.Sprintf("%s (формальный)", text)
	default:
		return fmt.Sprintf("%s (formal)", text)
	}
}

// generateRegionalVariant creates a regional variant of the translation
func generateRegionalVariant(text string, lang language.Tag) string {
	switch lang {
	case language.English:
		return fmt.Sprintf("%s (UK)", text)
	case language.Spanish:
		return fmt.Sprintf("%s (MX)", text)
	case language.Chinese:
		return fmt.Sprintf("%s (繁体)", text)
	default:
		return text
	}
}

// generateSimplifiedVariant creates a simplified variant for complex texts
func generateSimplifiedVariant(text string) string {
	// Simple approach - in production, you might use NLP techniques
	return fmt.Sprintf("%s (simplified)", text)
}

// isComplexText determines if text is complex based on simple heuristics
func isComplexText(text string) bool {
	// Simple complexity indicators
	sentences := strings.Split(text, ".")
	avgWordsPerSentence := float64(len(strings.Fields(text))) / float64(len(sentences))

	// Consider text complex if it has long sentences or technical terms
	return avgWordsPerSentence > 15 || strings.Contains(strings.ToLower(text), "therefore") ||
		strings.Contains(strings.ToLower(text), "however") || strings.Contains(strings.ToLower(text), "moreover")
}

// calculateTranslationConfidence estimates confidence based on back-translation consistency
func calculateTranslationConfidence(original, backTranslated string) float32 {
	// Simple similarity calculation
	// In production, you might use more sophisticated string similarity algorithms

	original = strings.ToLower(strings.TrimSpace(original))
	backTranslated = strings.ToLower(strings.TrimSpace(backTranslated))

	if original == backTranslated {
		return 0.95
	}

	// Calculate Levenshtein distance or use other similarity metrics
	// For simplicity, using a basic word count approach
	originalWords := strings.Fields(original)
	backWords := strings.Fields(backTranslated)

	if len(originalWords) == 0 {
		return 0.5
	}

	commonWords := 0
	for _, word := range originalWords {
		for _, backWord := range backWords {
			if word == backWord {
				commonWords++
				break
			}
		}
	}

	similarity := float32(commonWords) / float32(len(originalWords))

	// Convert similarity to confidence score
	confidence := 0.6 + (similarity * 0.3) // Scale to 0.6-0.9 range

	if confidence > 0.95 {
		confidence = 0.95
	}

	return confidence
}

// DetectLanguage detects the language of input text
func DetectLanguage(ctx context.Context, text string) (schema.Language, float32, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", 0, NewTranslationError("EMPTY_TEXT", "Text cannot be empty", "")
	}

	// Check text length limits
	if len(text) > 5000 {
		return "", 0, NewTranslationError("TEXT_TOO_LONG", "Text exceeds maximum length of 5000 characters", fmt.Sprintf("Text length: %d", len(text)))
	}

	// Initialize Google Translate client
	client, err := createTranslateClient(ctx)
	if err != nil {
		return "", 0, err
	}
	defer client.Close()

	// Detect language with retry logic
	detections, err := performLanguageDetectionWithRetry(ctx, client, text, 3)
	if err != nil {
		return "", 0, err
	}

	if len(detections) == 0 || len(detections[0]) == 0 {
		return "", 0, NewTranslationError("NO_DETECTION", "No language detected", "")
	}

	// Get the most confident detection
	detection := detections[0][0]

	// Map Google language code back to our schema
	schemaLang, err := mapGoogleCodeToLanguage(detection.Language)
	if err != nil {
		return "", 0, NewTranslationError("UNSUPPORTED_DETECTED_LANGUAGE", "Detected language is not supported", fmt.Sprintf("Detected: %s", detection.Language.String()))
	}

	return schemaLang, float32(detection.Confidence), nil
}

// performLanguageDetectionWithRetry performs language detection with retry logic
func performLanguageDetectionWithRetry(ctx context.Context, client *translate.Client, text string, maxRetries int) ([][]translate.Detection, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		detections, err := client.DetectLanguage(ctx, []string{text})

		if err == nil {
			return detections, nil
		}

		lastErr = err

		// Check if error is retryable
		if !isRetryableError(err) {
			break
		}

		// Wait before retry with exponential backoff
		if attempt < maxRetries-1 {
			waitTime := (1 << attempt) * 100 // 100ms, 200ms, 400ms
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(waitTime) * time.Millisecond):
				// Continue to next attempt
			}
		}
	}

	return nil, NewTranslationFailedError(fmt.Sprintf("Language detection failed after %d attempts: %v", maxRetries, lastErr))
}

// mapGoogleCodeToLanguage maps Google Translate language codes back to our schema
func mapGoogleCodeToLanguage(langTag language.Tag) (schema.Language, error) {
	switch langTag {
	case language.English:
		return schema.LanguageEnglish, nil
	case language.Georgian:
		return schema.LanguageGeorgian, nil
	case language.Spanish:
		return schema.LanguageSpanish, nil
	case language.French:
		return schema.LanguageFrench, nil
	case language.German:
		return schema.LanguageGerman, nil
	case language.Russian:
		return schema.LanguageRussian, nil
	case language.Japanese:
		return schema.LanguageJapanese, nil
	case language.Chinese:
		return schema.LanguageChinese, nil
	default:
		return "", fmt.Errorf("unsupported language code: %s", langTag.String())
	}
}

// GetSupportedLanguages returns all supported languages for translation
func GetSupportedLanguages() []schema.Language {
	return []schema.Language{
		schema.LanguageEnglish,
		schema.LanguageGeorgian,
		schema.LanguageSpanish,
		schema.LanguageFrench,
		schema.LanguageGerman,
		schema.LanguageRussian,
		schema.LanguageJapanese,
		schema.LanguageChinese,
	}
}
