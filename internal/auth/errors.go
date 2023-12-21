package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kmg7/fson/internal/err"
)

type AuthError = err.AppError

const (
	ErrTokenInvalid     = "invalid-token"
	ErrTokenExpired     = "expired-token"
	ErrAlreadyExists    = "already-exists"
	ErrInternal         = "internal-error"
	ErrNotFound         = "not-found"
	ErrNotAuthenticated = "not-authenticated"
)

func notFoundErr() *AuthError {
	return &AuthError{
		Code: ErrNotFound,
	}
}
func notAuthenticated() *AuthError {
	return &AuthError{
		Code: ErrNotAuthenticated,
	}
}

func alreadyExistsErr(whatEverExist string) *AuthError {
	return &AuthError{
		Code:     ErrAlreadyExists,
		Messages: []string{whatEverExist + " already taken"},
	}
}

func internalErr(err error) *AuthError {
	return &AuthError{
		Code:     ErrInternal,
		Internal: true,
		Messages: []string{"Unexpected error happened"},
		Err:      err,
	}
}
func invalidTokenErr() *AuthError {
	return &AuthError{
		Code: ErrTokenInvalid,
	}
}

func expiredTokenErr() *AuthError {
	return &AuthError{
		Code: ErrTokenExpired,
	}
}

func parseLibError(e error) *AuthError {
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
