package profiles

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        string    `json:"uid"`
	Name      string    `json:"name"`
	Password  string    `json:"pwd"`
	UpdatedAt time.Time `json:"uat"`
	CreatedAt time.Time `json:"cat"`
}

type UserUpdate struct {
	Id       *string
	Name     *string
	Password *string
}

type UserCreate struct {
	Name     *string
	Password *string
}

func (u *User) UpdateFrom(userUpdate *UserUpdate) User {
	new := *u
	if userUpdate.Password != nil {
		new.Password = *userUpdate.Password
	}

	if userUpdate.Name != nil {
		new.Name = *userUpdate.Name
	}
	new.UpdatedAt = time.Now()

	return new
}

func (uc *UserCreate) ToUser() User {
	time := time.Now()
	id := uuid.New()

	return User{
		Id:        id.String(),
		Name:      *uc.Name,
		Password:  *uc.Password,
		UpdatedAt: time,
		CreatedAt: time,
	}
}
