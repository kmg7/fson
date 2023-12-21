package config

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kmg7/fson/internal/auth"
	"github.com/kmg7/fson/internal/config"
	"github.com/kmg7/fson/internal/logger"
	"github.com/kmg7/fson/internal/server/middleware"
	"github.com/kmg7/fson/internal/server/utils"
	"github.com/kmg7/fson/internal/validator"
)

func adminSignIn(w http.ResponseWriter, r *http.Request) {
	dto := signIn{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.ValidateStruct(dto); err != nil {
		utils.InvalidBodyErrorResponse(w, r, err)
		return
	}
	tkn, err := auth.SuperUserSignIn(*dto.Username, *dto.Password)
	if err != nil {
		utils.ErrorResponse(w, r, 401, err)
		return
	}
	b, encErr := json.Marshal(tkn)
	if encErr != nil {
		w.WriteHeader(500)
		return
	}

	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}

func adminRenewToken(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value(middleware.TokenJwt).(*string)
	newToken, err := auth.RenewAdminToken(tkn)
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

	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}

func updateAdmin(w http.ResponseWriter, r *http.Request) {
	dto := adminUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.ValidateStruct(dto); err != nil {
		utils.InvalidBodyErrorResponse(w, r, err)
		return
	}
	tkn := r.Context().Value(middleware.TokenJwt).(*string)

	err := auth.SuperUserUpdate(tkn, dto.Password, dto.NewPassword, dto.Username)
	if err != nil {
		utils.ErrorResponse(w, r, 400, err)
		return
	}
	w.WriteHeader(200)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	usrs, err := auth.GetAllUsers()
	if err != nil {
		if err != nil {
			utils.ErrorResponse(w, r, 400, err)
			return
		}
	}
	b, encErr := json.Marshal(usrs)
	if encErr != nil {
		w.WriteHeader(500)
		return
	}

	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}

func getAllGroups(w http.ResponseWriter, r *http.Request) {
	usrs, err := auth.GetAllGroups()
	if err != nil {
		if err != nil {
			utils.ErrorResponse(w, r, 400, err)
			return
		}
	}
	b, encErr := json.Marshal(usrs)
	if encErr != nil {
		w.WriteHeader(500)
		return
	}

	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	dto := groupCreate{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.ValidateStruct(dto); err != nil {
		utils.InvalidBodyErrorResponse(w, r, err)
		return
	}

	grp, err := auth.CreateGroup(dto.Name, dto.Users)
	if err != nil {
		if err != nil {
			utils.ErrorResponse(w, r, 400, err)
			return
		}
	}
	b, encErr := json.Marshal(grp)
	if encErr != nil {

		w.WriteHeader(500)
		return
	}

	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}

func updateGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	dto := groupUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.ValidateStruct(dto); err != nil {
		utils.InvalidBodyErrorResponse(w, r, err)
		return
	}

	grp, err := auth.UpdateGroup(&id, dto.Name, dto.Users)
	if err != nil {
		if err != nil {
			utils.ErrorResponse(w, r, 400, err)
			return
		}
	}
	b, encErr := json.Marshal(grp)
	if encErr != nil {

		w.WriteHeader(500)
		return
	}

	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}

func deleteGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := auth.DeleteGroup(&id)
	if err != nil {
		if err != nil {
			utils.ErrorResponse(w, r, 400, err)
			return
		}
	}
	w.WriteHeader(200)
}

func getAuthConfig(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetAuthConfig()
	b, err := json.Marshal(&authConfig{
		TokensExpiresAfter:  &cfg.TokensExpiresAfter,
		TokenExpireTolerant: &cfg.TokenExpireTolerant})
	if err != nil {
		w.WriteHeader(500)
	}

	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}

func getAppConfig(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(config.GetAppCfg())
	if err != nil {
		w.WriteHeader(500)
	}

	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}

func patchAppConfig(w http.ResponseWriter, r *http.Request) {
	dto := appConfig{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}

	if err := validator.ValidateStruct(dto); err != nil {
		utils.InvalidBodyErrorResponse(w, r, err)
		return
	}
	if err := config.UpdateAppConfig(dto.AutoStart, dto.TempDir, dto.UploadDir); err != nil {
		if err.Internal {
			logger.Error(err.Err.Error())
		}
		utils.ErrorResponse(w, r, 400, err)
		return
	}
	b, err := json.Marshal(*config.GetAppCfg())
	if err != nil {
		w.WriteHeader(500)
		return
	}
	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)

}

func getTransferConfig(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(config.GetTransferCfg())
	if err != nil {
		w.WriteHeader(500)
	}
	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}

func postTransferPath(w http.ResponseWriter, r *http.Request) {
	cfg := transferPath{}
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.ValidateStruct(cfg); err != nil {
		utils.InvalidBodyErrorResponse(w, r, err)
		return
	}
	if err := config.AddTransferPath(cfg.Path); err != nil {
		if err.Internal {
			logger.Error(err.Err.Error())
		}
		utils.ErrorResponse(w, r, 400, err)
		return
	}
	b, err := json.Marshal(*config.GetTransferCfg())
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}

func patchTransferPath(w http.ResponseWriter, r *http.Request) {
	cfg := transferPath{}
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.ValidateStruct(cfg); err != nil {
		utils.InvalidBodyErrorResponse(w, r, err)
		return
	}
	id := chi.URLParam(r, "id")

	if err := config.UpdateTransferPath(&id, cfg.Path); err != nil {
		if err.Internal {
			logger.Error(err.Err.Error())
		}
		utils.ErrorResponse(w, r, 400, err)
		return
	}
	b, err := json.Marshal(*config.GetTransferCfg())
	if err != nil {
		w.WriteHeader(500)
		return
	}
	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}

func deleteTransferPath(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := config.DeleteTransferPath(&id); err != nil {
		if err.Internal {
			logger.Error(err.Err.Error())
		}
		utils.ErrorResponse(w, r, 400, err)
		return
	}
	b, err := json.Marshal(*config.GetTransferCfg())
	if err != nil {
		w.WriteHeader(500)
		return
	}
	utils.SetContentType(w, utils.ContentTypeJson)
	w.Write(b)
}
