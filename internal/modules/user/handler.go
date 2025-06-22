package user

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleUpdateUser(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*updateUserDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		user, err := UpdateUser(r.Context(), apiCfg.DB, &database.UpdateUserParams{
			Name: sql.NullString{String: body.Name, Valid: body.Name != ""},
			ID:   tokenPayload.UserID,
		})

		if err != nil {
			log.Println(err)
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, UserEntityToDto(user))
	}
}

func handleGetAuthUser(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		user, err := GetUserByID(r.Context(), apiCfg.DB, tokenPayload.UserID)
		if err != nil {
			shared.ResNotFound(w, shared.ErrUserNotFound)
			return
		}

		shared.ResOK(w, UserEntityToDto(user))
	}
}
