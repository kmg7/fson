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

// If any error occurs returns nil user and a state.
// Otherwise returns new Resource and a nil state.
func (p *Profiles) CreateUser(userCreate UserCreate) (*User, *State) {
	if _, usr := p.getUserWithName(*userCreate.Name); usr != nil {
		return nil, stateAlreadyExists("name")
	}

	newUser := userCreate.ToUser()
	newUsers := append(p.Users, newUser)
	newProfiles := p.UpdateFrom(&ProfilesUpdate{Users: &newUsers})

	if err := p.saveProfiles(newProfiles); err != nil {
		return nil, stateInternalErr(err)
	}
	return &newUser, nil
}

func (p *Profiles) getUser(id string) (int, *User) {
	for i, v := range p.Users {
		if v.Id == id {
			return i, &v
		}
	}
	return -1, nil
}

func (p *Profiles) getUserWithName(userName string) (int, *User) {
	for i, v := range p.Users {
		if v.Name == userName {
			return i, &v
		}
	}
	return -1, nil
}

// Updates the user with given id.
// If user not exists nil and an error.
// If something goes wrong during saving user returns nil and the error.
// If no change has made returns nil,nil.
func (p *Profiles) UpdateUser(id string, userUpdate *UserUpdate) (*User, *State) {
	i, user := p.getUser(id)
	if user == nil {
		return nil, stateNotFound()
	}

	if userUpdate.Name != nil {
		if _, anotherUsr := p.getUserWithName(*userUpdate.Name); anotherUsr != nil {
			return nil, stateAlreadyExists("name")
		}
	}

	newUsers := p.Users
	updatedUser := user.UpdateFrom(userUpdate)
	newUsers[i] = updatedUser

	newProfiles := p.UpdateFrom(&ProfilesUpdate{Users: &newUsers})
	if err := p.saveProfiles(newProfiles); err != nil {
		return nil, stateInternalErr(err)
	}

	return &updatedUser, nil
}
