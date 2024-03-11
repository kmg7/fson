package profiles

import (
	"sync"
	"time"

	"github.com/kmg7/fson/internal/state"
)

type AdminProfiles struct {
	filePath  string //`json:"-"`
	mutex     *sync.Mutex
	UpdatedAt time.Time `json:"uat"`
	CreatedAt time.Time `json:"cat"`
	Admins    []Admin   `json:"admins"`
}

func (ap *AdminProfiles) FilePath() string {
	return ap.filePath
}

// If any error occurs returns nil admin and a state.
// Otherwise returns new Resource and a nil state.
func (p *AdminProfiles) Create(create AdminCreate) (*Admin, state.Error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, adm := p.GetWithName(*create.Name); adm != nil {
		return nil, &state.ErrAlreadyExists{
			Resource: "name",
			Subject:  *create.Name,
		}
	}

	newAdmin := create.ToAdmin()
	newAdmins := append(p.Admins, newAdmin)

	newProfiles := *p
	newProfiles.Admins = newAdmins
	// if err := p.saveProfiles(newProfiles); err != nil {
	// return nil, stateInternalErr(err)
	// }
	return &newAdmin, nil
}

// Returns index and the resource. When resource nil index will be -1
func (p *AdminProfiles) Get(id string) (int, *Admin) {
	for i, v := range p.Admins {
		if v.Id == id {
			return i, &v
		}
	}
	return -1, nil
}

// Returns index and the resource. When resource nil index will be -1
func (p *AdminProfiles) GetWithName(name string) (int, *Admin) {
	for i, v := range p.Admins {
		if v.Name == name {
			return i, &v
		}
	}
	return -1, nil
}

func (p *AdminProfiles) Update(id string, update *AdminUpdate) (*Admin, state.Error) {
	i, admin := p.Get(id)
	if admin == nil {
		return nil, &state.ErrNotFound{
			Resource: "id",
			Subject:  id,
		}
	}

	if update.Name != nil {
		if _, anotherAdm := p.GetWithName(*update.Name); anotherAdm != nil {
			return nil, &state.ErrAlreadyExists{
				Resource: "name",
				Subject:  *update.Name,
			}
		}
	}

	newAdmins := p.Admins
	updatedAdmin := admin.UpdateFrom(update)
	newAdmins[i] = updatedAdmin

	return &updatedAdmin, nil
}
