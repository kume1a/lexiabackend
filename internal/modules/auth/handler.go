package auth

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGetAuthStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shared.ResOK(w, shared.OkDTO{Ok: true})
	}
}

func handleGoogleAuth(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequestBody[*googleSignInDTO](r)

		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		tokenPayload, httpErr := AuthWithGoogle(*apiCfg, r.Context(), body.Token)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		shared.ResOK(w, tokenPayload)
	}
}

func handleEmailAuth(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequestBody[*emailSignInDTO](r)

		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		tokenPayload, httpErr := AuthWithEmail(*apiCfg, r.Context(), body.Email, body.Password)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		shared.ResOK(w, tokenPayload)
	}
}
