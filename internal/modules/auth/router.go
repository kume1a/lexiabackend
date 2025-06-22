package auth

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/auth").Subrouter()

	r.HandleFunc("/status", handleGetAuthStatus()).Methods("GET")

	r.HandleFunc("/googleSignIn", handleGoogleAuth(apiCfg)).Methods("POST")
	r.HandleFunc("/emailSignIn", handleEmailAuth(apiCfg)).Methods("POST")

	return r
}
