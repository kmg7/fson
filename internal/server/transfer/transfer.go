/*
some settings wont need to change applications current state thus
for every settings a handler needed
any serving changes needs to be just in time
also stopping and scheduling starting of transfer services needs to be just in time
*/
package transfer

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kmg7/fson/internal/logger"
)

func Start(a string) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Post("/login", userSignIn)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/renew", renewToken)
		r.Post("/update", updateUser)
	})

	r.Route("/explore", func(r chi.Router) {
		r.Get("/", getExplore)
	})
	logger.Info("Transfer server starting on ", a)
	if err := http.ListenAndServe(a, r); err != nil {
		logger.Fatal(err.Error())
	}
	http.ListenAndServe(":8081", nil)
}
