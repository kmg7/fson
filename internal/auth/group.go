package auth

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	Id        string    `json:"uid"`
	Name      string    `json:"name"`
	Users     *[]string `json:"users"`
	UpdatedAt time.Time `json:"uat"`
	CreatedAt time.Time `json:"cat"`
}

func createGroup(name *string, users *[]string) (*Group, *AuthError) {
	time := now()
	id := uuid.New()
	group := Group{
		Id:        id.String(),
		Name:      *name,
		Users:     users,
		UpdatedAt: time,
		CreatedAt: time,
	}
	p := *profs.Groups
	p = append(p, group)
	return &group, save(nil, &p, nil)

}

func getGroup(id string) (int, *Group) {
	for i, v := range *profs.Groups {
		if v.Id == id {
			return i, &v
		}
	}
	return -1, nil
}

// func getGroupWithName(name string) (int, *Group) {
// 	for i, v := range *profs.Groups {
// 		if v.Name == name {
// 			return i, &v
// 		}
// 	}
// 	return -1, nil
// }

// Updates the group with given id.
// If group not exists nil and an error.
// If something goes wrong during saving group returns nil and the error.
// If no change has made returns nil,nil.
func updateGroup(id, name *string, users *[]string) (*Group, *AuthError) {
	changed := false
	i, grp := getGroup(*id)
	if grp == nil {
		return nil, notFoundErr()
	}
	if name != nil {
		if grp.Name != *name {
			grp.Name = *name
			changed = true
		}
	}
	if users != nil {
		grp.Users = users
		changed = true
	}
	if changed {
		grp.UpdatedAt = now()
		p := *profs.Groups
		p[i] = *grp
		return grp, save(nil, &p, nil)
	}
	return grp, nil

}

func deleteGroup(id *string) *AuthError {
	i, grp := getGroup(*id)
	if grp == nil {
		return notFoundErr()
	}
	p := *profs.Groups
	if len(p) == 1 {
		p = []Group{}
	} else {
		p = append(p[:i], p[i+1:]...)
	}
	return save(nil, &p, nil)
}
