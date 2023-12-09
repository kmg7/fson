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
			IsInternal: true,
			Code:       ErrInternal,
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

func UserUpdate(username string, newName *string, newPassword *string) (*string, *AuthError) {
	_, usr := getUserWithUsername(username)
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

func SuperUserUpdate(token, password, newPassword, name *string) *AuthError {
	clms, err := ValidateAdmin(token)
	if err != nil {
		return err
	}
	_, adm := getAdminWithAdminName(clms.Username)
	if adm == nil {
		return notAuthenticated()
	}
	if err := comapereHash(adm.Password, *password); err != nil {
		return notAuthenticated()
	}
	_, err = updateAdmin(adm.Id, newPassword, name)
	return err
}

func GetAllUsers(token *string) (*[]User, *AuthError) {
	_, err := ValidateAdmin(token)
	if err != nil {
		return nil, err
	}
	return profs.Users, nil
}

func DeleteUser(token *string, id *string) *AuthError {
	_, err := ValidateAdmin(token)
	if err != nil {
		return err
	}
	return deleteUser(*id)
}

func GetAllGroups(token *string) (*[]Group, *AuthError) {
	_, err := ValidateAdmin(token)
	if err != nil {
		return nil, err
	}
	return profs.Groups, nil
}

func CreateGroup(token, name *string, users *[]string) (*Group, *AuthError) {
	_, err := ValidateAdmin(token)
	if err != nil {
		return nil, err
	}
	return createGroup(name, users)
}

func UpdateGroup(token, id, name *string, users *[]string) (*Group, *AuthError) {
	_, err := ValidateAdmin(token)
	if err != nil {
		return nil, err
	}
	return updateGroup(id, name, users)
}

func DeleteGroup(token, id *string) *AuthError {
	_, err := ValidateAdmin(token)
	if err != nil {
		return err
	}
	return deleteGroup(id)
}

func now() time.Time {
	return time.Now()
}
