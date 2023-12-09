package auth

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

func createAdmin(name, password string) (*Admin, *AuthError) {
	if _, adm := getAdminWithAdminName(name); adm != nil {
		return nil, alreadyExistsErr("admin name")
	}
	time := now()
	id := uuid.New()
	admin := Admin{
		Id:        id.String(),
		Name:      name,
		Password:  password,
		UpdatedAt: time,
		CreatedAt: time,
	}
	p := *profs.Admins
	p = append(p, admin)
	return &admin, save(nil, nil, &p)
}

func getAdmin(id string) (int, *Admin) {
	for i, v := range *profs.Admins {
		if v.Id == id {
			return i, &v
		}
	}
	return -1, nil
}

func getAdminWithAdminName(adminName string) (int, *Admin) {
	for i, v := range *profs.Admins {
		if v.Name == adminName {
			return i, &v
		}
	}
	return -1, nil
}

// Updates the admin with given id.
// If admin not exists nil and an error.
// If something goes wrong during saving admin returns nil and the error.
// If no change has made returns nil,nil.
func updateAdmin(id string, password, name *string) (*Admin, *AuthError) {
	changed := false
	i, adm := getAdmin(id)
	if adm == nil {
		return nil, notFoundErr()
	}
	if password != nil {
		if adm.Password != *password {
			adm.Password = *password
			changed = true
		}
	}
	if name != nil {
		if adm.Name != *name {
			if _, anotherAdm := getAdminWithAdminName(*name); anotherAdm != nil {
				return nil, alreadyExistsErr("admin name")
			}
			adm.Name = *name
			changed = true
		}
	}
	if changed {
		adm.UpdatedAt = now()
		p := *profs.Admins
		p[i] = *adm
		return adm, save(nil, nil, &p)

	}
	return nil, nil

}

// func deleteAdmin(id string) *AuthError {
// 	changed := false
// 	i, adm := getAdmin(id)
// 	if adm == nil {
// 		return notFoundErr()
// 	}
// 	if changed {
// 		p := *profs.Admins
// 		p = append(p[:i], p[i+1:]...)
// 		p[i] = *adm
// 		return save(nil, nil, &p)
// 	}
// 	return nil
// }
