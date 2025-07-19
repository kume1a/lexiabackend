package translate

import (
	"lexia/ent/schema"
)

type TranslateRequestDTO struct {
	Text         string          `json:"text" validate:"required,min=1,max=5000"`
	LanguageFrom schema.Language `json:"languageFrom" validate:"required"`
	LanguageTo   schema.Language `json:"languageTo" validate:"required"`
}

type TranslationVariant struct {
	Text       string  `json:"text"`
	Confidence float32 `json:"confidence"`
}

type TranslateResponseDTO struct {
	OriginalText string               `json:"originalText"`
	LanguageFrom schema.Language      `json:"languageFrom"`
	LanguageTo   schema.Language      `json:"languageTo"`
	Translations []TranslationVariant `json:"translations"`
}

type DetectLanguageRequestDTO struct {
	Text string `json:"text" validate:"required,min=1,max=5000"`
}

type DetectLanguageResponseDTO struct {
	DetectedLanguage schema.Language `json:"detectedLanguage"`
	Confidence       float32         `json:"confidence"`
	Text             string          `json:"text"`
}

type SupportedLanguagesResponseDTO struct {
	Languages []schema.Language `json:"languages"`
}
