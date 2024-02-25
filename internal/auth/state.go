package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kmg7/fson/internal/state"
)

type AuthState = state.AppState

const (
	StatusTokenInvalid     = "invalid-token"
	StatusTokenExpired     = "expired-token"
	StatusInternal         = "internal-auth-error"
	StatusNotAuthenticated = "not-authenticated"
)

func notAuthenticated() *AuthState {
	return &AuthState{
		Status: StatusNotAuthenticated,
	}
}

func internalErr(err error) *AuthState {
	return &AuthState{
		Status:      StatusInternal,
		ErrInternal: true,
		Meta:        []string{"Unexpected error happened"},
		Err:         err,
	}
}

func invalidTokenErr() *AuthState {
	return &AuthState{
		Status: StatusTokenInvalid,
	}
}

func expiredTokenErr() *AuthState {
	return &AuthState{
		Status: StatusTokenExpired,
	}
}

func parseLibErr(e error) *AuthState {
	is := func(target error) bool {
		return errors.Is(e, target)
	}
	if is(jwt.ErrHashUnavailable) {
		return internalErr(e)
	}
	if is(jwt.ErrTokenExpired) {
		return expiredTokenErr()
	}
	return notAuthenticated()
}
