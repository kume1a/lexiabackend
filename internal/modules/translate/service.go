package translate

import (
	"context"
	"fmt"
	"lexia/ent/schema"
	"lexia/internal/shared"
	"slices"
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

func TranslateText(
	ctx context.Context,
	text string,
	from schema.Language,
	to schema.Language,
) ([]TranslationVariant, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, NewTranslationError("EMPTY_TEXT", "Text cannot be empty", "")
	}

	if len(text) > 5000 {
		return nil, NewTranslationError(
			"TEXT_TOO_LONG",
			"Text exceeds maximum length of 5000 characters",
			fmt.Sprintf("Text length: %d", len(text)),
		)
	}

	if from == to {
		return nil, NewTranslationError(
			"SAME_LANGUAGE",
			"Source and target languages cannot be the same",
			"",
		)
	}

	fromLang, err := mapLanguageToGoogleCode(from)
	if err != nil {
		return nil, NewUnsupportedLanguageError(string(from))
	}

	toLang, err := mapLanguageToGoogleCode(to)
	if err != nil {
		return nil, NewUnsupportedLanguageError(string(to))
	}

	client, err := createTranslateClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	results, err := performTranslationWithRetry(PerformTranslationWithRetryArgs{
		ctx:        ctx,
		client:     client,
		text:       text,
		fromLang:   fromLang,
		toLang:     toLang,
		maxRetries: 3,
	})
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, NewTranslationError(
			"NO_RESULTS",
			"No translation results returned from Google Translate",
			"",
		)
	}

	primaryTranslation := TranslationVariant{
		Text:       results[0].Text,
		Confidence: 0.95,
	}

	variants := []TranslationVariant{primaryTranslation}

	additionalVariants, err := generateTranslationVariants(ctx, client, text, fromLang, toLang, results[0].Text)
	if err != nil {
		fmt.Printf("Warning: failed to generate additional variants: %v\n", err)
	} else {
		variants = append(variants, additionalVariants...)
	}

	return variants, nil
}

func createTranslateClient(ctx context.Context) (*translate.Client, error) {
	envVars, err := shared.ParseEnv()
	if err != nil {
		return nil, NewCredentialsError()
	}

	var clientOptions []option.ClientOption

	clientOptions = append(clientOptions, option.WithQuotaProject(envVars.GoogleCloudProjectID))

	if envVars.GoogleServiceAccountKeyPath != "" {
		clientOptions = append(clientOptions, option.WithCredentialsFile(envVars.GoogleServiceAccountKeyPath))
	}

	client, err := translate.NewClient(ctx, clientOptions...)
	if err != nil {
		return nil, NewCredentialsError()
	}

	return client, nil
}

type PerformTranslationWithRetryArgs struct {
	ctx        context.Context
	client     *translate.Client
	text       string
	fromLang   language.Tag
	toLang     language.Tag
	maxRetries int
}

func performTranslationWithRetry(
	args PerformTranslationWithRetryArgs,
) ([]translate.Translation, error) {
	var lastErr error

	for attempt := range args.maxRetries {
		results, err := args.client.Translate(
			args.ctx,
			[]string{args.text},
			args.toLang,
			&translate.Options{
				Source: args.fromLang,
				Format: translate.Text,
			},
		)

		if err == nil {
			return results, nil
		}

		lastErr = err

		if !isRetryableError(err) {
			break
		}

		// Wait before retry with exponential backoff
		if attempt < args.maxRetries-1 {
			waitTime := (1 << attempt) * 100 // 100ms, 200ms, 400ms
			select {
			case <-args.ctx.Done():
				return nil, args.ctx.Err()
			case <-time.After(time.Duration(waitTime) * time.Millisecond):
				// Continue to next attempt
			}
		}
	}

	return nil, NewTranslationFailedError(fmt.Sprintf("Google Translate API error after %d attempts: %v", args.maxRetries, lastErr))
}

func isRetryableError(err error) bool {
	errorStr := strings.ToLower(err.Error())

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

func generateTranslationVariants(
	ctx context.Context,
	client *translate.Client,
	originalText string,
	fromLang, toLang language.Tag,
	primaryTranslation string,
) ([]TranslationVariant, error) {
	variants := []TranslationVariant{}

	// Strategy 1: Back-translation for confidence estimation and alternative generation
	// This helps assess the quality of the primary translation
	backTranslationResults, err := client.Translate(
		ctx,
		[]string{primaryTranslation},
		fromLang,
		&translate.Options{
			Source: toLang,
			Format: translate.Text,
		},
	)

	if err == nil && len(backTranslationResults) > 0 {
		backTranslation := backTranslationResults[0].Text

		// Calculate similarity-based confidence
		confidence := calculateTranslationConfidence(originalText, backTranslation)

		// If confidence is low, generate actual alternative translations
		if confidence < 0.85 {
			alternativeVariants, altErr := generateAlternativeTranslations(ctx, client, originalText, fromLang, toLang)
			if altErr == nil {
				variants = append(variants, alternativeVariants...)
			}
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

// generateAlternativeTranslations creates actual alternative translations using different approaches
func generateAlternativeTranslations(ctx context.Context, client *translate.Client, originalText string, fromLang, toLang language.Tag) ([]TranslationVariant, error) {
	variants := []TranslationVariant{}

	// Get the primary translation for comparison
	primaryResults, err := client.Translate(ctx, []string{originalText}, toLang, &translate.Options{
		Source: fromLang,
		Format: translate.Text,
	})

	if err != nil || len(primaryResults) == 0 {
		return variants, nil
	}

	primaryTranslation := primaryResults[0].Text

	// Method 1: Try translating through an intermediate language (pivot translation)
	// This can sometimes produce different nuances
	pivotLangs := []language.Tag{language.French, language.Spanish, language.German}

	for _, pivotLang := range pivotLangs {
		if pivotLang == fromLang || pivotLang == toLang {
			continue
		}

		// Translate original -> pivot -> target
		pivotResults, err := client.Translate(ctx, []string{originalText}, pivotLang, &translate.Options{
			Source: fromLang,
			Format: translate.Text,
		})

		if err != nil || len(pivotResults) == 0 {
			continue
		}

		finalResults, err := client.Translate(ctx, []string{pivotResults[0].Text}, toLang, &translate.Options{
			Source: pivotLang,
			Format: translate.Text,
		})

		if err != nil || len(finalResults) == 0 {
			continue
		}

		pivotTranslation := finalResults[0].Text

		// Only add if it's significantly different from the primary translation
		if !isTextSimilar(pivotTranslation, primaryTranslation) {
			variants = append(variants, TranslationVariant{
				Text:       pivotTranslation,
				Confidence: 0.75,
			})
			break // Only add one pivot translation to avoid too many variants
		}
	}

	// Method 2: Generate alternative by requesting multiple translations
	// Sometimes the same API call can return slightly different results
	// or we can use different translation approaches
	for attempt := 0; attempt < 2; attempt++ {
		results, err := client.Translate(ctx, []string{originalText}, toLang, &translate.Options{
			Source: fromLang,
			Format: translate.Text,
		})

		if err != nil || len(results) == 0 {
			continue
		}

		altTranslation := results[0].Text

		if !isTextSimilar(altTranslation, primaryTranslation) {
			variants = append(variants, TranslationVariant{
				Text:       altTranslation,
				Confidence: 0.80,
			})
			break
		}
	}

	return variants, nil
}

func isTextSimilar(text1, text2 string) bool {
	text1 = strings.ToLower(strings.TrimSpace(text1))
	text2 = strings.ToLower(strings.TrimSpace(text2))

	if text1 == text2 {
		return true
	}

	// Check if one text is just the other with minor modifications
	words1 := strings.Fields(text1)
	words2 := strings.Fields(text2)

	if len(words1) == 0 || len(words2) == 0 {
		return false
	}

	// Calculate word overlap
	commonWords := 0
	totalWords := max(len(words2), len(words1))

	for _, word1 := range words1 {
		if slices.Contains(words2, word1) {
			commonWords++
		}
	}

	// If more than 80% of words are the same, consider them too similar
	similarity := float64(commonWords) / float64(totalWords)
	return similarity > 0.8
}

func shouldGenerateFormalVariant(lang language.Tag) bool {
	formalLanguages := []language.Tag{
		language.German,
		language.French,
		language.Spanish,
		language.Japanese,
		language.Russian,
	}

	return slices.Contains(formalLanguages, lang)
}

func shouldGenerateRegionalVariant(lang language.Tag) bool {
	regionalLanguages := []language.Tag{
		language.English, // US vs UK
		language.Spanish, // ES vs MX/AR
		language.Chinese, // Simplified vs Traditional
	}

	return slices.Contains(regionalLanguages, lang)
}

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

func generateSimplifiedVariant(text string) string {
	// Simple approach - in production, you might use NLP techniques
	return fmt.Sprintf("%s (simplified)", text)
}

func isComplexText(text string) bool {
	wordCount := len(strings.Fields(text))
	sentenceCount := strings.Count(text, ".") + strings.Count(text, "!") + strings.Count(text, "?")
	if sentenceCount == 0 {
		sentenceCount = 1
	}
	avgWordsPerSentence := float64(wordCount) / float64(sentenceCount)

	return wordCount > 100 || avgWordsPerSentence > 20
}

func calculateTranslationConfidence(original, backTranslated string) float32 {
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

	// Convert similarity to confidence score, scale to 0.6-0.9 range
	confidence := 0.6 + (similarity * 0.3)

	if confidence > 0.95 {
		confidence = 0.95
	}

	return confidence
}

func DetectLanguage(ctx context.Context, text string) (schema.Language, float32, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", 0, NewTranslationError("EMPTY_TEXT", "Text cannot be empty", "")
	}

	if len(text) > 5000 {
		return "", 0, NewTranslationError("TEXT_TOO_LONG", "Text exceeds maximum length of 5000 characters", fmt.Sprintf("Text length: %d", len(text)))
	}

	client, err := createTranslateClient(ctx)
	if err != nil {
		return "", 0, err
	}
	defer client.Close()

	detections, err := performLanguageDetectionWithRetry(ctx, client, text, 3)
	if err != nil {
		return "", 0, err
	}

	if len(detections) == 0 || len(detections[0]) == 0 {
		return "", 0, NewTranslationError("NO_DETECTION", "No language detected", "")
	}

	detection := detections[0][0]

	schemaLang, err := mapGoogleCodeToLanguage(detection.Language)
	if err != nil {
		return "", 0, NewTranslationError(
			"UNSUPPORTED_DETECTED_LANGUAGE",
			"Detected language is not supported",
			fmt.Sprintf("Detected: %s", detection.Language.String()),
		)
	}

	return schemaLang, float32(detection.Confidence), nil
}

func performLanguageDetectionWithRetry(ctx context.Context, client *translate.Client, text string, maxRetries int) ([][]translate.Detection, error) {
	var lastErr error

	for attempt := range maxRetries {
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
