package profiles

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/kmg7/fson/internal/adapter"
	"github.com/kmg7/fson/internal/config"
)

type Profiles struct {
	init     bool
	filePath string
	fa       adapter.FileAdapter
	users    *UserProfiles
	admins   *AdminProfiles
	groups   *GroupProfiles
}

var (
	si              *Profiles
	instanciateOnce sync.Once
	groupMutex      sync.Mutex
	adminMutex      sync.Mutex
	userMutex       sync.Mutex
)

func Instance() *Profiles {
	instanciateOnce.Do(func() {
		i := &Profiles{}
		if err := i.initialize(); err != nil {
			log.Fatal(err.Error())
		}
		si = i
	})
	return si
}

func Setup(id, pwd string) error {
	var err error
	instanciateOnce.Do(func() {
		i := &Profiles{}
		err = i.setup(id, pwd)
		si = i
	})
	return err

}

func (p *Profiles) setup(id, pwd string) error {
	if p.init {
		return fmt.Errorf("profiles init already")
	}
	cfg := config.Instance()
	p.filePath = cfg.JoinConfigDir("profiles")
	p.fa = &adapter.File{
		Parse:   json.Marshal,
		Unparse: json.Unmarshal,
	}
	t := time.Now()
	if err := os.MkdirAll(p.filePath, 0700); err != nil {
		return err
	}

	p.admins = &AdminProfiles{
		filePath:  path.Join(p.filePath, "admin.cfg"),
		UpdatedAt: t,
		CreatedAt: t,
		Admins:    []Admin{{Id: id, Name: "admin", Password: pwd, UpdatedAt: t, CreatedAt: t}},
	}

	p.groups = &GroupProfiles{
		filePath:  path.Join(p.filePath, "groups.cfg"),
		UpdatedAt: t,
		CreatedAt: t,
	}

	p.users = &UserProfiles{
		filePath:  path.Join(p.filePath, "users.cfg"),
		UpdatedAt: t,
		CreatedAt: t,
	}

	if err := p.saveAdmins(); err != nil {
		return err
	}

	if err := p.saveGroups(); err != nil {
		return err
	}

	if err := p.saveUsers(); err != nil {
		return err
	}

	p.init = true
	return nil

}

func (p *Profiles) initialize() error {
	if p.init {
		return fmt.Errorf("profiles init already")
	}
	cfg := config.Instance()
	p.filePath = cfg.JoinConfigDir("profiles")
	p.fa = &adapter.File{
		Parse:   json.Marshal,
		Unparse: json.Unmarshal,
	}

	// Admin User Group initialization
	p.admins = &AdminProfiles{
		filePath: path.Join(p.filePath, "admin.cfg"),
		mutex:    &adminMutex,
	}
	p.groups = &GroupProfiles{
		filePath: path.Join(p.filePath, "groups.cfg"),
		mutex:    &groupMutex,
	}

	p.users = &UserProfiles{
		filePath: path.Join(p.filePath, "users.cfg"),
		mutex:    &userMutex,
	}

	if err := p.readAdmins(); err != nil {
		return err
	}

	if err := p.readGroups(); err != nil {
		return err
	}

	if err := p.readUsers(); err != nil {
		return err
	}

	p.init = true
	return nil
}
