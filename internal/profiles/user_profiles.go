package profiles

import (
	"sync"
	"time"

	"github.com/kmg7/fson/internal/state"
)

type UserProfiles struct {
	filePath  string
	mutex     *sync.Mutex
	UpdatedAt time.Time `json:"uat"`
	CreatedAt time.Time `json:"cat"`
	Users     []User    `json:"users"`
}

func (p *UserProfiles) FilePath() string {
	return p.filePath
}

// If any error occurs returns nil user and a state.
// Otherwise returns new Resource and a nil state.

func (p *UserProfiles) Create(userCreate UserCreate) (*User, state.Error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, usr := p.GetWithName(*userCreate.Name); usr != nil {
		return nil, &state.ErrAlreadyExists{
			Resource: "name",
			Subject:  *userCreate.Name,
		}
	}

	newUser := userCreate.ToUser()
	newUsers := append(p.Users, newUser)
	newProfiles := *p
	newProfiles.Users = newUsers

	return &newUser, nil
}

func (p *UserProfiles) Get(id string) (int, *User) {
	for i, v := range p.Users {
		if v.Id == id {
			return i, &v
		}
	}
	return -1, nil
}

func (p *UserProfiles) GetWithName(userName string) (int, *User) {
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
func (p *UserProfiles) Update(id string, userUpdate *UserUpdate) (*User, state.Error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	i, user := p.Get(id)
	if user == nil {
		return nil, &state.ErrAlreadyExists{
			Resource: "id",
			Subject:  id,
		}
	}

	if userUpdate.Name != nil {
		if _, anotherUsr := p.GetWithName(*userUpdate.Name); anotherUsr != nil {
			return nil, &state.ErrAlreadyExists{
				Resource: "name",
				Subject:  *userUpdate.Name,
			}
		}
	}

	newUsers := p.Users
	updatedUser := user.UpdateFrom(userUpdate)
	newUsers[i] = updatedUser

	return &updatedUser, nil
}
