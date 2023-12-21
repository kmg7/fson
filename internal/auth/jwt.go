package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UID string
	jwt.RegisteredClaims
}

func Validate(token *string) (*Claims, *AuthError) {
	claims := &Claims{}
	_, err := parseAndValidate(claims, token, keyFunc)
	return claims, err
}

func ValidateAdmin(token *string) (*Claims, *AuthError) {
	claims := &Claims{}
	_, err := parseAndValidate(claims, token, keyFuncAdmin)
	return claims, err

}

func parseAndValidate(claims *Claims, token *string, keyFunc jwt.Keyfunc) (*jwt.Token, *AuthError) {
	tkn, err := jwt.ParseWithClaims(*token, claims, keyFunc)
	if err != nil {
		return nil, parseLibError(err)
	}
	if !tkn.Valid {
		return nil, invalidTokenErr()
	}
	if time.Until(claims.ExpiresAt.Time) < cfg.TokenExpireTolerant {
		return nil, expiredTokenErr()
	}
	return tkn, nil

}

func GetToken(user *User) (*string, *AuthError) {
	claims := &Claims{
		UID: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, internalErr(err)
	}
	return &tokenString, nil

}

func RenewToken(token *string) (*string, *AuthError) {
	claims := &Claims{}
	if _, err := parseAndValidate(claims, token, keyFunc); err != nil {
		return nil, err
	}
	claims.ExpiresAt = jwt.NewNumericDate(expTime())

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenStr, err := newToken.SignedString(secret)
	if err != nil {
		return nil, internalErr(err)
	}
	return &newTokenStr, nil
}

func RenewAdminToken(token *string) (*string, *AuthError) {
	claims := &Claims{}
	if _, err := parseAndValidate(claims, token, keyFuncAdmin); err != nil {
		return nil, err
	}
	claims.ExpiresAt = jwt.NewNumericDate(expTime())

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenStr, err := newToken.SignedString(secret)
	if err != nil {
		return nil, internalErr(err)
	}
	return &newTokenStr, nil
}

func GetTokenAdmin(admin *Admin) (*string, *AuthError) {
	claims := &Claims{
		UID: admin.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(admSecret)
	if err != nil {
		return nil, internalErr(err)
	}
	return &tokenString, nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	return secret, nil
}

func keyFuncAdmin(token *jwt.Token) (interface{}, error) {
	return admSecret, nil
}

func expTime() time.Time {
	return time.Now().Add(cfg.TokensExpiresAfter)
}
