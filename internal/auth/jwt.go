package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UID string
	jwt.RegisteredClaims
}

type JWT struct {
	expireTolerant time.Duration
	expire         time.Duration
	secret         []byte
}

func (j *JWT) keyFunc(token *jwt.Token) (interface{}, error) {
	return j.secret, nil
}

func (j *JWT) expTime() time.Time {
	return time.Now().Add(j.expire)
}

func (j *JWT) Validate(token *string) *AuthState {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(*token, claims, j.keyFunc)
	if err != nil {
		return parseLibErr(err)
	}
	if !tkn.Valid {
		return invalidTokenErr()
	}
	if time.Until(claims.ExpiresAt.Time) < j.expireTolerant {
		return expiredTokenErr()
	}
	return nil
}

func (j *JWT) GetToken(uid string) (*string, *AuthState) {
	claims := &Claims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(j.expTime()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return nil, internalErr(err)
	}
	return &tokenString, nil

}

func (j *JWT) RenewToken(token *string) (*string, *AuthState) {
	claims := &Claims{}
	if err := j.Validate(token); err != nil {
		return nil, err
	}
	claims.ExpiresAt = jwt.NewNumericDate(j.expTime())

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenStr, err := newToken.SignedString(j.secret)
	if err != nil {
		return nil, internalErr(err)
	}
	return &newTokenStr, nil
}
