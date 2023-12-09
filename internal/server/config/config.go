package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kmg7/fson/internal/auth"
	"github.com/kmg7/fson/internal/config"
	"github.com/kmg7/fson/internal/logger"
	"github.com/kmg7/fson/internal/validator"
)

func StartConfigServer(a string) {
	auth.Init()
	if err := validator.Instantiate(); err != nil {
		logger.Fatal(err.Error())
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", SignIn)
		r.Post("/renew", RenewToken)
		r.Post("/update", UpdateAdmin)
	})
	r.Route("/users", func(r chi.Router) {
		r.Get("/", GetAllUsers)

	})
	r.Route("/groups", func(r chi.Router) {
		r.Get("/", GetAllGroups)
		r.Post("/", CreateGroup)
		r.Route("/{id}", func(r chi.Router) {
			r.Patch("/", UpdateGroup)
			r.Delete("/", DeleteGroup)
		})
	})
	r.Route("/cfg", func(r chi.Router) {
		r.Route("/app", func(r chi.Router) {
			r.Get("/", GetAppConfig)
			r.Patch("/", PatchAppConfig)
		})
		r.Route("/transfer", func(r chi.Router) {
			r.Get("/", GetTransferConfig)
			r.Post("/", PostTransferPath)
			r.Patch("/{id}", PatchTransferPath)
			r.Delete("/{id}", DeleteTransferPath)
		})
		r.Route("/auth", func(r chi.Router) {
			r.Get("/", GetAuthConfig)
		})
	})
	logger.Info("Config server starting on ", a)
	if err := http.ListenAndServe(a, r); err != nil {
		logger.Fatal(err.Error())
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	dto := AdminSignInDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.Validate(dto); err != nil {
		ErrorResponse(w, r, false, "invalid-body", err.String())
		return
	}
	tkn, err := auth.SuperUserSignIn(*dto.Username, *dto.Password)
	if err != nil {
		AuthErrorResponse(w, r, err.IsInternal, err.Code, err.Messages)
		return
	}
	b, encErr := json.Marshal(tkn)
	if encErr != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func RenewToken(w http.ResponseWriter, r *http.Request) {
	const bearer = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, bearer) {
		w.WriteHeader(401)
		return
	}
	tkn := strings.TrimPrefix(authHeader, bearer)
	newToken, err := auth.RenewAdminToken(tkn)
	if err != nil {
		if err != nil {
			if err.IsInternal {
				logger.Error(err)
			}
			AuthErrorResponse(w, r, err.IsInternal, err.Code, err.Messages)
			return
		}
	}
	b, encErr := json.Marshal(newToken)
	if encErr != nil {

		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	dto := AdminUpdateDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.Validate(dto); err != nil {
		ErrorResponse(w, r, false, "invalid-body", err.String())
		return
	}
	const bearer = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, bearer) {
		w.WriteHeader(401)
		return
	}
	tkn := strings.TrimPrefix(authHeader, bearer)

	err := auth.SuperUserUpdate(&tkn, dto.Password, dto.NewPassword, dto.Username)
	if err != nil {
		AuthErrorResponse(w, r, err.IsInternal, err.Code, err.Messages)
		return
	}
	w.WriteHeader(200)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	const bearer = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, bearer) {
		w.WriteHeader(401)
		return
	}
	tkn := strings.TrimPrefix(authHeader, bearer)
	usrs, err := auth.GetAllUsers(&tkn)
	if err != nil {
		if err != nil {
			if err.IsInternal {
				logger.Error(err)
			}
			AuthErrorResponse(w, r, err.IsInternal, err.Code, err.Messages)
			return
		}
	}
	b, encErr := json.Marshal(usrs)
	if encErr != nil {

		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	const bearer = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, bearer) {
		w.WriteHeader(401)
		return
	}
	tkn := strings.TrimPrefix(authHeader, bearer)
	usrs, err := auth.GetAllGroups(&tkn)
	if err != nil {
		if err != nil {
			if err.IsInternal {
				logger.Error(err)
			}
			AuthErrorResponse(w, r, err.IsInternal, err.Code, err.Messages)
			return
		}
	}
	b, encErr := json.Marshal(usrs)
	if encErr != nil {

		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {

	const bearer = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, bearer) {
		w.WriteHeader(401)
		return
	}
	tkn := strings.TrimPrefix(authHeader, bearer)

	dto := GroupCreateDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.Validate(dto); err != nil {
		ErrorResponse(w, r, false, "invalid-body", err.String())
		return
	}

	grp, err := auth.CreateGroup(&tkn, dto.Name, dto.Users)
	if err != nil {
		if err != nil {
			if err.IsInternal {
				logger.Error(err)
			}
			AuthErrorResponse(w, r, err.IsInternal, err.Code, err.Messages)
			return
		}
	}
	b, encErr := json.Marshal(grp)
	if encErr != nil {

		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {

	const bearer = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, bearer) {
		w.WriteHeader(401)
		return
	}
	tkn := strings.TrimPrefix(authHeader, bearer)

	id := chi.URLParam(r, "id")
	dto := GroupUpdateDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.Validate(dto); err != nil {
		ErrorResponse(w, r, false, "invalid-body", err.String())
		return
	}

	grp, err := auth.UpdateGroup(&tkn, &id, dto.Name, dto.Users)
	if err != nil {
		if err != nil {
			if err.IsInternal {
				logger.Error(err)
			}
			AuthErrorResponse(w, r, err.IsInternal, err.Code, err.Messages)
			return
		}
	}
	b, encErr := json.Marshal(grp)
	if encErr != nil {

		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {

	const bearer = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, bearer) {
		w.WriteHeader(401)
		return
	}
	tkn := strings.TrimPrefix(authHeader, bearer)

	id := chi.URLParam(r, "id")
	err := auth.DeleteGroup(&tkn, &id)
	if err != nil {
		if err != nil {
			if err.IsInternal {
				logger.Error(err)
			}
			AuthErrorResponse(w, r, err.IsInternal, err.Code, err.Messages)
			return
		}
	}
	w.WriteHeader(200)
}

func GetAuthConfig(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetAuthConfig()
	b, err := json.Marshal(AuthConfigDTO{
		TokensExpiresAfter:  cfg.TokensExpiresAfter,
		TokenExpireTolerant: cfg.TokenExpireTolerant})
	if err != nil {
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func GetAppConfig(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(config.GetAppCfg())
	if err != nil {
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func PatchAppConfig(w http.ResponseWriter, r *http.Request) {
	dto := AppConfigDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(400)
		return
	}

	if err := validator.Validate(dto); err != nil {
		ErrorResponse(w, r, false, "invalid-body", err.String())
		return
	}
	if err := config.UpdateAppConfig(dto.AutoStart, dto.TempDir, dto.UploadDir); err != nil {
		if err.Internal {
			logger.Error(err.Err.Error())
		}
		ErrorResponse(w, r, err.Internal, err.Code, err.Messages)
		return
	}
	b, err := json.Marshal(*config.GetAppCfg())
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(b)

}

func GetTransferConfig(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(config.GetTransferCfg())
	if err != nil {
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func PostTransferPath(w http.ResponseWriter, r *http.Request) {
	cfg := TransferPathDTO{}
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.Validate(cfg); err != nil {
		ErrorResponse(w, r, false, "invalid-body", err.String())
		return
	}
	if err := config.AddTransferPath(cfg.Path); err != nil {
		if err.Internal {
			logger.Error(err.Err.Error())
		}
		ErrorResponse(w, r, err.Internal, err.Code, err.Messages)
		return
	}
	b, err := json.Marshal(*config.GetTransferCfg())
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(b)
}

func PatchTransferPath(w http.ResponseWriter, r *http.Request) {
	cfg := TransferPathDTO{}
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		w.WriteHeader(400)
		return
	}
	if err := validator.Validate(cfg); err != nil {
		ErrorResponse(w, r, false, "invalid-body", err.String())
		return
	}
	id := chi.URLParam(r, "id")

	if err := config.UpdateTransferPath(id, cfg.Path); err != nil {
		if err.Internal {
			logger.Error(err.Err.Error())
		}
		ErrorResponse(w, r, err.Internal, err.Code, err.Messages)
		return
	}
	b, err := json.Marshal(*config.GetTransferCfg())
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(b)
}

func DeleteTransferPath(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := config.DeleteTransferPath(id); err != nil {
		if err.Internal {
			logger.Error(err.Err.Error())
		}
		ErrorResponse(w, r, err.Internal, err.Code, err.Messages)
		return
	}
	b, err := json.Marshal(*config.GetTransferCfg())
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(b)
}

func AuthErrorResponse(w http.ResponseWriter, r *http.Request, isInternal bool, code string, msgs []string) {
	b, _ := json.Marshal(ErrorResponseBody{Code: code, Err: msgs})
	if isInternal {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(401)
	}
	w.Write(b)

}

func ErrorResponse(w http.ResponseWriter, r *http.Request, isInternal bool, code string, msgs []string) {
	b, _ := json.Marshal(ErrorResponseBody{Code: code, Err: msgs})
	if isInternal {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(400)
	}
	w.Write(b)

}

func UnexpectedInternalErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	b, _ := json.Marshal(ErrorResponseBody{Code: "unexpected-internal", Err: []string{}})
	w.WriteHeader(500)
	logger.Error(err.Error())
	w.Write(b)
}

type ErrorResponseBody struct {
	Code string   `json:"code"`
	Err  []string `json:"errors"`
}
