package profiles

import (
	"fmt"
	"path"
)

func (p *Profiles) saveAdmins() error {
	if err := p.fa.ParseAndSave(p.admins); err != nil {
		return fmt.Errorf("saving admin profiles failed err: %w", err)
	}
	return nil
}

func (p *Profiles) readAdmins() error {
	read := &AdminProfiles{
		filePath: path.Join(p.filePath, "admin.cfg"),
		mutex:    &adminMutex,
	}

	if err := p.fa.ReadAndParse(read); err != nil {
		return fmt.Errorf("reading admin profiles failed err: %w", err)
	}
	return nil

}

func (p *Profiles) saveUsers() error {
	if err := p.fa.ParseAndSave(p.users); err != nil {
		return fmt.Errorf("saving user profiles failed err: %w", err)
	}
	return nil
}

func (p *Profiles) readUsers() error {
	read := &UserProfiles{
		filePath: path.Join(p.filePath, "users.cfg"),
		mutex:    &userMutex,
	}

	if err := p.fa.ReadAndParse(read); err != nil {
		return fmt.Errorf("reading user profiles failed err: %w", err)
	}
	return nil

}

func (p *Profiles) saveGroups() error {
	if err := p.fa.ParseAndSave(p.groups); err != nil {
		return fmt.Errorf("saving group profiles failed err: %w", err)
	}
	return nil
}

func (p *Profiles) readGroups() error {
	read := &GroupProfiles{
		filePath: path.Join(p.filePath, "groups.cfg"),
		mutex:    &groupMutex,
	}

	if err := p.fa.ReadAndParse(read); err != nil {
		return fmt.Errorf("reading group profiles failed err: %w", err)
	}
	return nil

}
