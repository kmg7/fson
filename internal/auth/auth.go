package auth

import (
	"fmt"
	"log"
	"sync"

	"github.com/kmg7/fson/internal/config"
)

type Auth struct {
	init  bool
	User  *JWT
	Admin *JWT
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
	acfg := cfg.AuthConfig()
	a.User = &JWT{
		expireTolerant: acfg.TokenExpireTolerant,
		expire:         acfg.TokensExpiresAfter,
		secret:         []byte(acfg.Secret),
	}
	a.Admin = &JWT{
		expireTolerant: acfg.TokenExpireTolerant,
		expire:         acfg.TokensExpiresAfter,
		secret:         []byte(acfg.AdminSecret),
	}
	return nil
}
