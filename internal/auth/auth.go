package auth

import (
	"fmt"
	"log"
	"sync"

	"github.com/kmg7/fson/internal/config"
	"github.com/kmg7/fson/internal/crypt"
	"github.com/kmg7/fson/internal/profiles"
)

type Auth struct {
	init     bool
	Admin    *AdminAuth
	UserAuth *UserAuth
	UserJWT  *JWT
	AdminJWT *JWT
}

var (
	si              *Auth
	instanciateOnce sync.Once
)

func Instance() *Auth {
	instanciateOnce.Do(func() {
		i := &Auth{}
		if err := i.initialize(); err != nil {
			log.Fatal(err.Error())
		}
		si = i
	})
	return si
}

func (a *Auth) initialize() error {
	if a.init {
		return fmt.Errorf("auth init already")
	}
	cfg := config.Instance()
	acfg := cfg.GetAuth()
	profiles := profiles.Instance()

	adminCrypt := crypt.Instance(crypt.Options{
		BcryptCost: 8,
	})

	userCrypt := crypt.Instance(crypt.Options{
		// BcryptCost: 4,
	})

	a.UserJWT = &JWT{
		expireTolerant: acfg.TokenExpireTolerant,
		expire:         acfg.TokensExpiresAfter,
		secret:         []byte(acfg.Secret),
	}

	a.AdminJWT = &JWT{
		expireTolerant: acfg.TokenExpireTolerant,
		expire:         acfg.TokensExpiresAfter,
		secret:         []byte(acfg.AdminSecret),
	}

	a.Admin = &AdminAuth{
		profiles: profiles,
		hash:     adminCrypt.Bcrypt,
		compare:  adminCrypt.BcryptCompare}

	a.UserAuth = &UserAuth{
		profiles: profiles,
		hash:     userCrypt.Bcrypt,
		compare:  userCrypt.BcryptCompare}

	return nil
}
