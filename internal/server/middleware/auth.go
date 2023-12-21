package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kmg7/fson/internal/auth"
	"github.com/kmg7/fson/internal/server/utils"
)

const (
	TokenJwt CtxKey = "JWT_TOKEN"
)

type CtxKey string

// Parses jwt token to context
func AuthJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const bearer = "Bearer "
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, bearer) {
			w.WriteHeader(401)
			return
		}
		tkn := strings.TrimPrefix(authHeader, bearer)

		ctx := context.WithValue(r.Context(), TokenJwt, &tkn)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthenticateAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tkn := r.Context().Value(TokenJwt).(*string)
		if _, err := auth.ValidateAdmin(tkn); err != nil {
			utils.ErrorResponse(w, r, 401, err)
			return
		}
		next.ServeHTTP(w, r)

	})
}
