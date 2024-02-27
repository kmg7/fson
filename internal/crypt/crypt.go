package crypt

import "golang.org/x/crypto/bcrypt"

type Options struct {
	BcryptCost int
}

func Instance(opt Options) *crypt {
	if opt.BcryptCost <= 0 {
		opt.BcryptCost = 4
	}

	return &crypt{
		bcryptCost: opt.BcryptCost,
	}
}

type crypt struct {
	bcryptCost int
}

func (c *crypt) Bcrypt(text []byte) (*[]byte, error) {
	p, err := bcrypt.GenerateFromPassword(text, c.bcryptCost)
	if err != nil {
		return nil, err
	}

	return &p, err
}

func (c *crypt) BcryptCompare(hash, text []byte) error {
	return bcrypt.CompareHashAndPassword(hash, text)
}
