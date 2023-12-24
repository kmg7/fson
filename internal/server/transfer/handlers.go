package transfer

import (
	"encoding/json"
	"net/http"

	"github.com/kmg7/fson/internal/auth"
	"github.com/kmg7/fson/internal/server/middleware"
	"github.com/kmg7/fson/internal/server/utils"
	"github.com/kmg7/fson/internal/validator"
)

func getExplore(w http.ResponseWriter, r *http.Request) {

}
func userSignIn(w http.ResponseWriter, r *http.Request) {
	dto := signIn{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.ValidateStruct(dto); err != nil {
		utils.InvalidBodyErrorResponse(w, r, err)
		return
	}
	tkn, err := auth.UserSignIn(*dto.Username, *dto.Password)
	if err != nil {
		utils.ErrorResponse(w, r, 401, err)
		return
	}
	b, encErr := json.Marshal(tkn)
	if encErr != nil {
		w.WriteHeader(500)
		return
	}

	w.Write(b)
}

func renewToken(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value(middleware.TokenJwt).(*string)
	newToken, err := auth.RenewToken(tkn)
	if err != nil {
		if err != nil {
			utils.ErrorResponse(w, r, 400, err)
			return
		}
	}
	b, encErr := json.Marshal(newToken)
	if encErr != nil {
		w.WriteHeader(500)
		return
	}

	w.Write(b)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	dto := userUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.ValidateStruct(dto); err != nil {
		utils.InvalidBodyErrorResponse(w, r, err)
		return
	}
	tkn := r.Context().Value(middleware.TokenJwt).(*string)

	err := auth.UserUpdate(tkn, dto.Password, dto.NewPassword, dto.Username)
	if err != nil {
		utils.ErrorResponse(w, r, 400, err)
		return
	}
	w.WriteHeader(200)
}
