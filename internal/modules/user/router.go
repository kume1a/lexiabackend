package user

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/users").Subrouter()

	r.HandleFunc("/updateMe", shared.AuthMW(handleUpdateUser(apiCfg))).Methods("PATCH")
	r.HandleFunc("/authUser", shared.AuthMW(handleGetAuthUser(apiCfg))).Methods("GET")

	return r
}
