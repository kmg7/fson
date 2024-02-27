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

// If any error occurs returns nil admin and a state.
// Otherwise returns new Resource and a nil state.
func (p *Profiles) CreateAdmin(adminCreate AdminCreate) (*Admin, *State) {
	if _, adm := p.getAdminWithName(*adminCreate.Name); adm != nil {
		return nil, stateAlreadyExists("name")
	}

	newAdmin := adminCreate.ToAdmin()
	newAdmins := append(p.Admins, newAdmin)
	newProfiles := p.UpdateFrom(&ProfilesUpdate{Admins: &newAdmins})

	if err := p.saveProfiles(newProfiles); err != nil {
		return nil, stateInternalErr(err)
	}
	return &newAdmin, nil
}

func (p *Profiles) getAdmin(id string) (int, *Admin) {
	for i, v := range p.Admins {
		if v.Id == id {
			return i, &v
		}
	}
	return -1, nil
}

func (p *Profiles) getAdminWithName(adminName string) (int, *Admin) {
	for i, v := range p.Admins {
		if v.Name == adminName {
			return i, &v
		}
	}
	return -1, nil
}

func (p *Profiles) UpdateAdmin(id string, adminUpdate *AdminUpdate) (*Admin, *State) {
	i, admin := p.getAdmin(id)
	if admin == nil {
		return nil, stateNotFound()
	}

	if adminUpdate.Name != nil {
		if _, anotherAdm := p.getAdminWithName(*adminUpdate.Name); anotherAdm != nil {
			return nil, stateAlreadyExists("name")
		}
	}

	newAdmins := p.Admins
	updatedAdmin := admin.UpdateFrom(adminUpdate)
	newAdmins[i] = updatedAdmin

	newProfiles := p.UpdateFrom(&ProfilesUpdate{Admins: &newAdmins})
	if err := p.saveProfiles(newProfiles); err != nil {
		return nil, stateInternalErr(err)
	}

	return &updatedAdmin, nil
}
