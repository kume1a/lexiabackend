package translate

import (
	"errors"
	"fmt"
)

// Translation service specific errors
var (
	ErrEmptyText               = errors.New("text cannot be empty")
	ErrSameLanguage            = errors.New("source and target languages cannot be the same")
	ErrUnsupportedLanguage     = errors.New("unsupported language")
	ErrTranslationFailed       = errors.New("translation failed")
	ErrLanguageDetectionFailed = errors.New("language detection failed")
	ErrNoTranslationResults    = errors.New("no translation results returned")
	ErrGoogleCredentials       = errors.New("Google Cloud credentials not configured")
)

// TranslationError represents a structured error for translation operations
type TranslationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *TranslationError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewTranslationError creates a new TranslationError
func NewTranslationError(code, message, details string) *TranslationError {
	return &TranslationError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common error constructors
func NewUnsupportedLanguageError(lang string) *TranslationError {
	return NewTranslationError(
		"UNSUPPORTED_LANGUAGE",
		"Language not supported",
		fmt.Sprintf("Language '%s' is not supported for translation", lang),
	)
}

func NewTranslationFailedError(details string) *TranslationError {
	return NewTranslationError(
		"TRANSLATION_FAILED",
		"Translation service failed",
		details,
	)
}

func NewCredentialsError() *TranslationError {
	return NewTranslationError(
		"CREDENTIALS_ERROR",
		"Google Cloud credentials not configured",
		"Please set GOOGLE_APPLICATION_CREDENTIALS environment variable or configure Application Default Credentials",
	)
}

func NewTextTooLongError(length int) *TranslationError {
	return NewTranslationError(
		"TEXT_TOO_LONG",
		"Text exceeds maximum length",
		fmt.Sprintf("Text length %d exceeds the maximum allowed length of 5000 characters", length),
	)
}

func NewNoDetectionError() *TranslationError {
	return NewTranslationError(
		"NO_DETECTION",
		"No language could be detected",
		"The input text may be too short or contain unsupported characters",
	)
}
