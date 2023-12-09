package auth

import "golang.org/x/crypto/bcrypt"

func generateHash(password string) (string, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(p), err
}
func comapereHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
