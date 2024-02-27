package profiles

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kmg7/fson/internal/adapter"
	"github.com/kmg7/fson/internal/config"
)

type Profiles struct {
	init      bool
	filePath  string
	fa        adapter.FileAdapter
	UpdatedAt time.Time `json:"uat"`
	CreatedAt time.Time `json:"cat"`
	Users     []User    `json:"users"`
	Groups    []Group   `json:"groups"`
	Admins    []Admin   `json:"admins"`
}

type ProfilesUpdate struct {
	Users  *[]User  `json:"users"`
	Groups *[]Group `json:"groups"`
	Admins *[]Admin `json:"admins"`
}

var (
	si              *Profiles
	instanciateOnce sync.Once
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

func (p *Profiles) initialize() error {
	if p.init {
		return fmt.Errorf("profiles init already")
	}
	cfg := config.Instance()
	p.filePath = cfg.JoinConfigDir("profiles.cfg")
	p.fa = &adapter.File{
		Parse:   json.Marshal,
		Unparse: json.Unmarshal,
	}
	return nil
}

func (p *Profiles) FilePath() string {
	return p.filePath
}

func (old *Profiles) UpdateFrom(update *ProfilesUpdate) *Profiles {
	if update == nil {
		return nil
	}
	new := *old
	updated := false

	if update.Users != nil {
		new.Users = *update.Users
		updated = true
	}

	if update.Admins != nil {
		new.Admins = *update.Admins
		updated = true
	}

	if update.Groups != nil {
		new.Groups = *update.Groups
		updated = true
	}

	if updated {
		return &new
	}
	return old
}

func (p *Profiles) saveProfiles(newP *Profiles) error {
	if err := p.fa.ParseAndSave(newP); err != nil {
		return fmt.Errorf("saving profiles failed err: %w", err)
	}
	p = newP
	return nil
}

func (p *Profiles) readProfiles() error {
	read := &Profiles{
		init:     p.init,
		filePath: p.filePath,
		fa:       p.fa,
	}
	if err := p.fa.ReadAndParse(read); err != nil {
		return fmt.Errorf("reading profiles failed err: %w", err)
	}
	p = read
	return nil

}
