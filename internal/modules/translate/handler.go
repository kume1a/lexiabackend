package translate

import (
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
)

func handleTranslate(_ *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		var body TranslateRequestDTO
		if validationErr := shared.BindAndValidate(c, &body); validationErr != nil {
			shared.ResValidationError(c, validationErr)
			return
		}

		if body.LanguageFrom == body.LanguageTo {
			shared.ResBadRequest(c, "Source and target languages must be different")
			return
		}

		translations, err := TranslateText(c.Request.Context(), body.Text, body.LanguageFrom, body.LanguageTo)
		if err != nil {
			// Handle custom translation errors
			if translationErr, ok := err.(*TranslationError); ok {
				switch translationErr.Code {
				case "EMPTY_TEXT", "SAME_LANGUAGE", "TEXT_TOO_LONG":
					shared.ResBadRequest(c, translationErr.Message)
					return
				case "UNSUPPORTED_LANGUAGE":
					shared.ResBadRequest(c, translationErr.Message)
					return
				case "CREDENTIALS_ERROR":
					shared.ResInternalServerError(c, "Translation service configuration error")
					return
				default:
					shared.ResInternalServerError(c, "Translation failed")
					return
				}
			}
			shared.ResInternalServerError(c, "Translation service unavailable")
			return
		}

		response := TranslateResponseDTO{
			OriginalText: body.Text,
			LanguageFrom: body.LanguageFrom,
			LanguageTo:   body.LanguageTo,
			Translations: translations,
		}

		shared.ResOK(c, response)
	}
}

func handleDetectLanguage(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		var body DetectLanguageRequestDTO
		if validationErr := shared.BindAndValidate(c, &body); validationErr != nil {
			shared.ResValidationError(c, validationErr)
			return
		}

		// Detect language
		detectedLang, confidence, err := DetectLanguage(c.Request.Context(), body.Text)
		if err != nil {
			// Handle custom translation errors
			if translationErr, ok := err.(*TranslationError); ok {
				switch translationErr.Code {
				case "EMPTY_TEXT", "TEXT_TOO_LONG":
					shared.ResBadRequest(c, translationErr.Message)
					return
				case "UNSUPPORTED_DETECTED_LANGUAGE":
					shared.ResBadRequest(c, translationErr.Message)
					return
				case "CREDENTIALS_ERROR":
					shared.ResInternalServerError(c, "Language detection service configuration error")
					return
				default:
					shared.ResInternalServerError(c, "Language detection failed")
					return
				}
			}
			shared.ResInternalServerError(c, "Language detection service unavailable")
			return
		}

		response := DetectLanguageResponseDTO{
			DetectedLanguage: detectedLang,
			Confidence:       confidence,
			Text:             body.Text,
		}

		shared.ResOK(c, response)
	}
}

func handleGetSupportedLanguages(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		languages := GetSupportedLanguages()
		response := SupportedLanguagesResponseDTO{
			Languages: languages,
		}

		shared.ResOK(c, response)
	}
}
