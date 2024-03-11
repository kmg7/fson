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
