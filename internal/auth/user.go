package auth

import (
	"github.com/kmg7/fson/internal/profiles"
	"github.com/kmg7/fson/internal/state"
)

type UserAuth struct {
	profiles *profiles.Profiles
	hash     func(text []byte) (*[]byte, error)
	compare  func(hash, text []byte) error
}

func (a *UserAuth) Signin(name, password *string) (state.Error, *string) {
	serr, user := a.profiles.GetUserName(name)
	if serr != nil {
		return serr, nil
	}

	err := a.compare([]byte(user.Password), []byte(*password))
	if err != nil {
		return &state.ErrNotFound{
			Resource: "user",
		}, nil
	}

	return nil, &user.Id
}

func (a *UserAuth) Signup(name, password *string) (state.Error, *profiles.User) {
	hashed, err := a.hash([]byte(*password))
	if err != nil {
		return &state.ErrInternal{
			While: "hashing password",
			Err:   err}, nil
	}

	pwd := string(*hashed)
	create := profiles.UserCreate{
		Name:     name,
		Password: &pwd,
	}

	serr, user := a.profiles.CreateUser(create)
	if serr != nil {
		return serr, nil
	}
	return nil, user
}
