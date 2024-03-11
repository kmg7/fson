package auth

import (
	"github.com/kmg7/fson/internal/profiles"
	"github.com/kmg7/fson/internal/state"
)

type AdminAuth struct {
	profiles *profiles.Profiles
	hash     func(text []byte) (*[]byte, error)
	compare  func(hash, text []byte) error
}

func (a *AdminAuth) Signin(name, password *string) (state.Error, *string) {
	serr, admin := a.profiles.GetAdminName(name)
	if serr != nil {
		return serr, nil
	}

	err := a.compare([]byte(admin.Password), []byte(*password))
	if err != nil {
		return &state.ErrNotFound{
			Resource: "admin",
		}, nil
	}

	return nil, &admin.Id
}
