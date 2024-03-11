package profiles

import (
	"sync"
	"time"

	"github.com/kmg7/fson/internal/state"
)

type GroupProfiles struct {
	mutex     *sync.Mutex
	filePath  string
	UpdatedAt time.Time `json:"uat"`
	CreatedAt time.Time `json:"cat"`
	Groups    []Group   `json:"groups"`
}

func (p *GroupProfiles) FilePath() string {
	return p.filePath
}

// If any error occurs returns nil group and a state.
// Otherwise returns new Resource and a nil state.
func (p *GroupProfiles) Create(groupCreate GroupCreate) (*Group, state.Error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, adm := p.GetWithName(*groupCreate.Name); adm != nil {
		return nil, &state.ErrAlreadyExists{
			Resource: "name",
			Subject:  *groupCreate.Name,
		}
	}

	newGroup := groupCreate.ToGroup()
	newGroups := append(p.Groups, newGroup)

	newProfiles := *p
	newProfiles.Groups = newGroups

	return &newGroup, nil
}

func (p *GroupProfiles) Get(id string) (int, *Group) {
	for i, v := range p.Groups {
		if v.Id == id {
			return i, &v
		}
	}
	return -1, nil
}

func (p *GroupProfiles) GetWithName(groupName string) (int, *Group) {
	for i, v := range p.Groups {
		if v.Name == groupName {
			return i, &v
		}
	}
	return -1, nil
}

func (p *GroupProfiles) Update(id string, groupUpdate *GroupUpdate) (*Group, state.Error) {
	i, group := p.Get(id)
	if group == nil {
		return nil, &state.ErrNotFound{
			Resource: "group",
			Subject:  id,
		}
	}

	if groupUpdate.Name != nil {
		if _, anotherAdm := p.GetWithName(*groupUpdate.Name); anotherAdm != nil {
			return nil, &state.ErrAlreadyExists{
				Resource: "name",
				Subject:  *groupUpdate.Name,
			}
		}
	}

	newGroups := p.Groups
	updatedGroup := group.UpdateFrom(groupUpdate)
	newGroups[i] = updatedGroup

	return &updatedGroup, nil
}
