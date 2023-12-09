package auth

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

func createUser(name, password string) (*User, *AuthError) {
	time := now()
	id := uuid.New()
	if _, usr := getUserWithUsername(name); usr != nil {
		return nil, alreadyExistsErr("username")
	}
	user := User{
		Id:        id.String(),
		Name:      name,
		Password:  password,
		UpdatedAt: time,
		CreatedAt: time,
	}
	p := *profs.Users
	p = append(p, user)
	return &user, save(&p, nil, nil)
}

func getUser(id string) (int, *User) {
	for i, v := range *profs.Users {
		if v.Id == id {
			return i, &v
		}
	}
	return -1, nil
}

func getUserWithUsername(username string) (int, *User) {
	for i, v := range *profs.Users {
		if v.Name == username {
			return i, &v
		}
	}
	return -1, nil
}

// Updates the user with given id.
// If user not exists nil and an error.
// If something goes wrong during saving user returns nil and the error.
// If no change has made returns nil,nil.
func updateUser(id string, password, name *string) (*User, *AuthError) {
	changed := false
	i, usr := getUser(id)
	if usr == nil {
		return nil, notFoundErr()
	}

	if password != nil {
		if usr.Password != *password {
			usr.Password = *password
			changed = true
		}
	}
	if name != nil {
		if usr.Name != *name {
			if _, usr := getUserWithUsername(*name); usr != nil {
				return nil, alreadyExistsErr("username")
			}
			usr.Name = *name
			changed = true
		}
	}
	if changed {
		usr.UpdatedAt = now()
		p := *profs.Users
		p[i] = *usr
		return usr, save(&p, nil, nil)
	}
	return usr, nil

}

func deleteUser(id string) *AuthError {
	changed := false
	i, usr := getUser(id)
	if usr == nil {
		return notFoundErr()
	}
	if changed {
		p := *profs.Users
		if len(p) == 1 {
			p = []User{}
		} else {
			p = append(p[:i], p[i+1:]...)
		}
		p[i] = *usr
		return save(&p, nil, nil)
	}
	return nil
}
