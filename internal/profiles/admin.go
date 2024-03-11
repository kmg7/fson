package profiles

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	Id          string    `json:"uid"`
	Name        string    `json:"name"`
	Password    string    `json:"pwd"`
	LastLoginAt time.Time `json:"lat"`
	UpdatedAt   time.Time `json:"uat"`
	CreatedAt   time.Time `json:"cat"`
}

type AdminUpdate struct {
	Id       *string
	Name     *string
	Password *string
}

type AdminCreate struct {
	Name     *string
	Password *string
}

func (a *Admin) UpdateFrom(adminUpdate *AdminUpdate) Admin {
	new := *a
	if adminUpdate.Password != nil {
		new.Password = *adminUpdate.Password
	}

	if adminUpdate.Name != nil {
		new.Name = *adminUpdate.Name
	}
	new.UpdatedAt = time.Now()
	return new
}

func (adminCreate *AdminCreate) ToAdmin() Admin {
	time := time.Now()
	id := uuid.New()
	return Admin{
		Id:        id.String(),
		Name:      *adminCreate.Name,
		Password:  *adminCreate.Password,
		UpdatedAt: time,
		CreatedAt: time,
	}
}
