package translate

import (
	"fmt"
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

		supportedLanguages := GetSupportedLanguages()
		isFromSupported := false
		isToSupported := false

		for _, lang := range supportedLanguages {
			if lang == body.LanguageFrom {
				isFromSupported = true
			}
			if lang == body.LanguageTo {
				isToSupported = true
			}
		}

		if !isFromSupported {
			shared.ResBadRequest(c, fmt.Sprintf("Source language '%s' is not supported", body.LanguageFrom))
			return
		}

		if !isToSupported {
			shared.ResBadRequest(c, fmt.Sprintf("Target language '%s' is not supported", body.LanguageTo))
			return
		}

		translations, err := TranslateText(c.Request.Context(), body.Text, body.LanguageFrom, body.LanguageTo)
		if err != nil {
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

func handleDetectLanguage(_ *shared.ApiConfig) gin.HandlerFunc {
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

		detectedLang, confidence, err := DetectLanguage(c.Request.Context(), body.Text)
		if err != nil {
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

func handleGetSupportedLanguages(_ *shared.ApiConfig) gin.HandlerFunc {
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
