package profiles

import "github.com/kmg7/fson/internal/state"

func (p *Profiles) GetAdmins() (state.Error, []Admin) {
	return nil, p.admins.Admins
}

func (p *Profiles) GetAdmin(id *string) (state.Error, *Admin) {
	var adm Admin
	_, admin := p.admins.Get(*id)
	if admin == nil {
		return &state.ErrNotFound{
			Resource: "admin",
			Subject:  *id,
		}, nil
	}
	adm = *admin
	return nil, &adm
}

func (p *Profiles) GetAdminName(name *string) (state.Error, *Admin) {
	var adm Admin
	_, admin := p.admins.GetWithName(*name)
	if admin == nil {
		return &state.ErrNotFound{
			Resource: "admin",
			Subject:  *name,
		}, nil
	}
	adm = *admin
	return nil, &adm
}

func (p *Profiles) CreateAdmin(create AdminCreate) (state.Error, *Admin) {
	backUp := *p.admins
	created, serr := p.admins.Create(create)
	if serr != nil {
		return serr, nil
	}

	p.admins.mutex.Lock()
	defer p.admins.mutex.Unlock()

	if err := p.saveAdmins(); err != nil {
		p.admins = &backUp // TODO 1:check if file corrupt better rollback strategy
		return &state.ErrInternal{
			While: "creating admin and saving profiles",
			Err:   err,
		}, nil
	}
	return nil, created
}

func (p *Profiles) UpdateAdmin(id string, update *AdminUpdate) (state.Error, *Admin) {
	backUp := *p.admins

	updated, serr := p.admins.Update(id, update)
	if serr != nil {
		return serr, nil
	}

	p.admins.mutex.Lock()
	defer p.admins.mutex.Unlock()

	if err := p.saveAdmins(); err != nil {
		p.admins = &backUp // TODO 1
		return &state.ErrInternal{
			While: "updating admin and saving profiles",
			Err:   err,
		}, nil
	}

	return nil, updated
}

func (p *Profiles) GetUsers() (state.Error, []User) {
	return nil, p.users.Users
}

func (p *Profiles) GetUser(id *string) (state.Error, *User) {
	var usr User
	_, user := p.users.Get(*id)
	if user == nil {
		return &state.ErrNotFound{
			Resource: "user",
			Subject:  *id,
		}, nil
	}
	usr = *user
	return nil, &usr
}

func (p *Profiles) GetUserName(name *string) (state.Error, *User) {
	var usr User
	_, user := p.users.GetWithName(*name)
	if user == nil {
		return &state.ErrNotFound{
			Resource: "user",
			Subject:  *name,
		}, nil
	}
	usr = *user
	return nil, &usr
}

func (p *Profiles) CreateUser(create UserCreate) (state.Error, *User) {
	backUp := *p.users
	created, serr := p.users.Create(create)
	if serr != nil {
		return serr, nil
	}

	p.users.mutex.Lock()
	defer p.users.mutex.Unlock()

	if err := p.saveUsers(); err != nil {
		p.users = &backUp // TODO 1
		return &state.ErrInternal{
			While: "creating user and saving profiles",
			Err:   err,
		}, nil
	}
	return nil, created
}

func (p *Profiles) UpdateUser(id string, update *UserUpdate) (state.Error, *User) {
	backUp := *p.users
	updated, serr := p.users.Update(id, update)
	if serr != nil {
		return serr, nil
	}

	p.users.mutex.Lock()
	defer p.users.mutex.Unlock()

	if err := p.saveUsers(); err != nil {
		p.users = &backUp // TODO 1
		return &state.ErrInternal{
			While: "updating user and saving profiles",
			Err:   err,
		}, nil
	}
	return nil, updated
}

func (p *Profiles) GetGroups() (state.Error, []Group) {
	return nil, p.groups.Groups
}

func (p *Profiles) GetGroup(id *string) (state.Error, *Group) {
	var grp Group
	_, group := p.groups.Get(*id)
	if group == nil {
		return &state.ErrNotFound{
			Resource: "group",
			Subject:  *id,
		}, nil
	}
	grp = *group
	return nil, &grp
}

func (p *Profiles) GetGroupName(name *string) (state.Error, *Group) {
	var grp Group
	_, group := p.groups.GetWithName(*name)
	if group == nil {
		return &state.ErrNotFound{
			Resource: "group",
			Subject:  *name,
		}, nil
	}
	grp = *group
	return nil, &grp
}

func (p *Profiles) CreateGroup(create GroupCreate) (state.Error, *Group) {
	backUp := *p.groups
	created, serr := p.groups.Create(create)
	if serr != nil {
		return serr, nil
	}

	p.groups.mutex.Lock()
	defer p.groups.mutex.Unlock()

	if err := p.saveGroups(); err != nil {
		p.groups = &backUp // TODO 1
		return &state.ErrInternal{
			While: "creating group and saving profiles",
			Err:   err,
		}, nil
	}
	return nil, created
}

func (p *Profiles) UpdateGroup(id string, update *GroupUpdate) (state.Error, *Group) {
	backUp := *p.groups
	updated, serr := p.groups.Update(id, update)
	if serr != nil {
		return serr, nil
	}

	p.groups.mutex.Lock()
	defer p.groups.mutex.Unlock()

	if err := p.saveGroups(); err != nil {
		p.groups = &backUp // TODO 1
		return &state.ErrInternal{
			While: "updating group and saving profiles",
			Err:   err,
		}, nil
	}
	return nil, updated
}
