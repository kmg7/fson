package auth

import (
	"time"

	"github.com/kmg7/fson/internal/config"
	"github.com/kmg7/fson/internal/logger"
)

var secret []byte
var admSecret []byte

var cfg *config.AuthConfig
var profsPath string

func Init() {
	cfg = config.GetAuthConfig()
	secret = []byte(cfg.Secret)
	admSecret = []byte(cfg.AdminSecret)
	profsPath = config.AuthProfilesFilePath()
	if err := read(); err != nil {
		logger.Fatal("Profiles read failed")
	}
}

func UserSignUp(username, password string) *AuthError {
	hashedPassword, err := generateHash(password)
	if err != nil {
		return &AuthError{
			Internal: true,
			Code:     ErrInternal,
		}
	}
	_, authErr := createUser(username, string(hashedPassword))
	return authErr
}

func UserSignIn(username, password string) (*string, *AuthError) {
	_, usr := getUser(username)
	if usr == nil {
		return nil, notFoundErr()
	}
	if err := comapereHash(usr.Password, password); err != nil {
		return nil, notAuthenticated()
	}
	return GetToken(usr)
}

func UserUpdate(uid, password, newName, newPassword *string) (*string, *AuthError) {
	_, usr := getUser(*uid)
	if usr == nil {
		return nil, notFoundErr()
	}
	usr, err := updateUser(usr.Id, newName, newPassword)
	if err != nil {
		return nil, err
	}
	return GetToken(usr)
}

func SuperUserSignIn(name, password string) (*string, *AuthError) {
	_, adm := getAdminWithAdminName(name)
	if adm == nil {
		return nil, notFoundErr()
	}
	if err := comapereHash(adm.Password, password); err != nil {
		return nil, notAuthenticated()
	}
	return GetTokenAdmin(adm)
}

func SuperUserUpdate(admin, password, newPassword, name *string) *AuthError {
	_, adm := getAdmin(*admin)
	if adm == nil {
		return notAuthenticated()
	}
	if err := comapereHash(adm.Password, *password); err != nil {
		return notAuthenticated()
	}
	_, err := updateAdmin(adm.Id, newPassword, name)
	return err
}

func GetAllUsers() (*[]User, *AuthError) {
	return profs.Users, nil
}

func DeleteUser(id *string) *AuthError {
	return deleteUser(*id)
}

func GetAllGroups() (*[]Group, *AuthError) {
	return profs.Groups, nil
}

func CreateGroup(name *string, users *[]string) (*Group, *AuthError) {
	return createGroup(name, users)
}

func UpdateGroup(id, name *string, users *[]string) (*Group, *AuthError) {
	return updateGroup(id, name, users)
}

func DeleteGroup(id *string) *AuthError {
	return deleteGroup(id)
}

func now() time.Time {
	return time.Now()
}
