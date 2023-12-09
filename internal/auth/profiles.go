package auth

import (
	"encoding/json"
	"os"
	"time"

	"github.com/kmg7/fson/internal/logger"
	"github.com/kmg7/fson/pkg/fsutils"
)

type Profiles struct {
	UpdatedAt time.Time `json:"uat"`
	CreatedAt time.Time `json:"cat"`
	Users     *[]User   `json:"usrs"`
	Groups    *[]Group  `json:"grps"`
	Admins    *[]Admin  `json:"adms"`
}

var profs Profiles

func save(usrs *[]User, grps *[]Group, admins *[]Admin) *AuthError {
	changed := false
	prf := profs
	if usrs != nil {
		changed = true
		prf.Users = usrs
	}
	if grps != nil {
		changed = true
		prf.Groups = grps
	}
	if admins != nil {
		changed = true
		prf.Admins = admins
	}
	if changed {

		file, err := os.Create(profsPath)
		if err != nil {
			logger.Fatal("Cannot create profs file", err)
			return internalErr(err)
		}
		defer file.Close()
		prf.UpdatedAt = now()
		data, err := json.MarshalIndent(prf, "", " ")
		if err != nil {
			logger.Fatal("Cannot parse profs", err)
			return internalErr(err)
		}

		if _, err := file.Write(data); err != nil {
			logger.Error("Cannot write profs file", err)
			return internalErr(err)
		}
		profs = prf
	}

	return nil
}

func read() *AuthError {
	ex, err := fsutils.Exists(profsPath)
	if err != nil {
		logger.Fatal("Cannot verify if profs exists", err)
		return internalErr(err)
	}
	if !ex {
		logger.Info("Not found any profile")
		profs = Profiles{
			CreatedAt: now(),
			UpdatedAt: now(),
			Users:     &[]User{},
			Groups:    &[]Group{},
			Admins:    &[]Admin{},
		}
		pwd, err := generateHash("admin")
		if err != nil {
			logger.Fatal("Something went wrong while creating first super user", err)
			return internalErr(err)
		}
		createAdmin("admin", string(pwd))
	}
	file, err := os.Open(profsPath)
	if err != nil {
		logger.Fatal("Cannot open profs file", err)
		return internalErr(err)
	}

	defer file.Close()
	data, err := os.ReadFile(profsPath)
	if err != nil {
		logger.Fatal("Cannot read profs file", err)
		return internalErr(err)
	}

	err = json.Unmarshal(data, &profs)
	if err != nil {
		logger.Fatal("Cannot parse profs file", err)
		return internalErr(err)
	}
	logger.Info("Profiles read")
	return nil

}
