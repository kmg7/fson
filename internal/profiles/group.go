package profiles

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	Id        string    `json:"gid"`
	Name      string    `json:"name"`
	Users     []string  `json:"users"`
	UpdatedAt time.Time `json:"uat"`
	CreatedAt time.Time `json:"cat"`
}

type GroupUpdate struct {
	Id    *string
	Name  *string
	Users *[]string
}

type GroupCreate struct {
	Name     *string
	Password *string
}

func (g *Group) UpdateFrom(groupUpdate *GroupUpdate) Group {
	new := *g
	if groupUpdate.Users != nil {
		new.Users = *groupUpdate.Users
	}

	if groupUpdate.Name != nil {
		new.Name = *groupUpdate.Name
	}
	new.UpdatedAt = time.Now()
	return new
}

func (groupCreate *GroupCreate) ToGroup() Group {
	time := time.Now()
	id := uuid.New()
	return Group{
		Id:        id.String(),
		Name:      *groupCreate.Name,
		UpdatedAt: time,
		CreatedAt: time,
	}
}

// If any error occurs returns nil group and a state.
// Otherwise returns new Resource and a nil state.
func (p *Profiles) CreateGroup(groupCreate GroupCreate) (*Group, *State) {
	if _, adm := p.getGroupWithName(*groupCreate.Name); adm != nil {
		return nil, stateAlreadyExists("name")
	}

	newGroup := groupCreate.ToGroup()
	newGroups := append(p.Groups, newGroup)
	newProfiles := p.UpdateFrom(&ProfilesUpdate{Groups: &newGroups})

	if err := p.saveProfiles(newProfiles); err != nil {
		return nil, stateInternalErr(err)
	}
	return &newGroup, nil
}

func (p *Profiles) getGroup(id string) (int, *Group) {
	for i, v := range p.Groups {
		if v.Id == id {
			return i, &v
		}
	}
	return -1, nil
}

func (p *Profiles) getGroupWithName(groupName string) (int, *Group) {
	for i, v := range p.Groups {
		if v.Name == groupName {
			return i, &v
		}
	}
	return -1, nil
}

func (p *Profiles) UpdateGroup(id string, groupUpdate *GroupUpdate) (*Group, *State) {
	i, group := p.getGroup(id)
	if group == nil {
		return nil, stateNotFound()
	}

	if groupUpdate.Name != nil {
		if _, anotherAdm := p.getGroupWithName(*groupUpdate.Name); anotherAdm != nil {
			return nil, stateAlreadyExists("name")
		}
	}

	newGroups := p.Groups
	updatedGroup := group.UpdateFrom(groupUpdate)
	newGroups[i] = updatedGroup

	newProfiles := p.UpdateFrom(&ProfilesUpdate{Groups: &newGroups})
	if err := p.saveProfiles(newProfiles); err != nil {
		return nil, stateInternalErr(err)
	}

	return &updatedGroup, nil
}
