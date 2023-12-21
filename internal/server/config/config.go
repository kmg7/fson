package config

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kmg7/fson/internal/auth"
	"github.com/kmg7/fson/internal/logger"
	mw "github.com/kmg7/fson/internal/server/middleware"
	"github.com/kmg7/fson/internal/server/utils"
	"github.com/kmg7/fson/internal/validator"
)

func StartConfigServer(address string) {
	auth.Init()
	if err := validator.Instantiate(); err != nil {
		logger.Fatal(err.Error())
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.AllowContentType(utils.ContentTypeJson))
	r.Use(middleware.CleanPath)
	r.Post("/login", adminSignIn)

	r.Route("/", func(r chi.Router) {
		r.Use(mw.AuthJwt)
		r.Use(mw.AuthenticateAdmin)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/renew", adminRenewToken)
			r.Post("/update", updateAdmin)
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/", getAllUsers)
		})

		r.Route("/groups", func(r chi.Router) {
			r.Get("/", getAllGroups)
			r.Post("/", createGroup)
			r.Route("/{id}", func(r chi.Router) {
				r.Patch("/", updateGroup)
				r.Delete("/", deleteGroup)
			})
		})

		r.Route("/cfg", func(r chi.Router) {
			r.Route("/app", func(r chi.Router) {
				r.Get("/", getAppConfig)
				r.Patch("/", patchAppConfig)
			})

			r.Route("/transfer", func(r chi.Router) {
				r.Get("/", getTransferConfig)
				r.Post("/", postTransferPath)
				r.Patch("/{id}", patchTransferPath)
				r.Delete("/{id}", deleteTransferPath)
			})

			r.Route("/auth", func(r chi.Router) {
				r.Get("/", getAuthConfig)
			})
		})

	})
	logger.Info("Config server starting on ", address)
	if err := http.ListenAndServe(address, r); err != nil {
		logger.Fatal(err.Error())
	}
}
