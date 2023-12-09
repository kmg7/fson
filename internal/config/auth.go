package config

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/kmg7/fson/internal/logger"
	"github.com/kmg7/fson/pkg/fsutils"
)

type AuthConfig struct {
	TokensExpiresAfter  time.Duration
	TokenExpireTolerant time.Duration
	Secret              string
	AdminSecret         string
}
type Admin struct {
	Username string
	Password string
}

var authCfgPath string

func initAuthConfig() AuthConfig {
	p, err := authConfigFilePath()
	if err != nil {
		logger.Error(err)
	}
	authCfgPath = p
	ex, err := fsutils.Exists(authCfgPath)
	if err != nil {
		logger.Error(err)
	}
	if !ex {
		logger.Info("Not found any auth config")
		confd := defaultAuthConfig()
		writeAuthConfig(confd)
		return confd
	}
	conf, err := readAuthConfig()
	if err != nil {
		confd := defaultAuthConfig()
		logger.Info("Overriding broken auth config")
		writeAuthConfig(confd)
		return confd
	}
	return conf
}

func readAuthConfig() (AuthConfig, error) {
	var conf AuthConfig
	file, err := os.Open(authCfgPath)
	if err != nil {
		logger.Error("Cannot open auth config")
		return conf, err
	}

	defer file.Close()
	data, err := os.ReadFile(authCfgPath)
	if err != nil {
		logger.Error("Cannot read auth config", err)
		return conf, err
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		logger.Error("Cannot parse auth config", err)
	}
	logger.Info("Auth config read")
	return conf, err
}

func writeAuthConfig(config AuthConfig) error {
	file, err := os.Create(authCfgPath)
	if err != nil {
		logger.Error("Cannot create auth config", err)
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		logger.Error("Cannot parse auth config", err)
		return err
	}

	if _, err := file.Write(data); err != nil {
		logger.Error("Cannot write auth config", err)
		return err
	}
	logger.Info("Auth config wrote")
	return nil
}

func authConfigFilePath() (string, error) {
	appCfg, err := appConfigDir()
	if err != nil {
		logger.Error("Cannot read app config dir", err.Error())
		return "", err
	}
	return path.Join(appCfg, ".AUTH_CFG"), nil
}

func AuthProfilesFilePath() string {
	appCfg, err := appConfigDir()
	if err != nil {
		logger.Fatal("Cannot read app config dir", err.Error())
	}
	return path.Join(appCfg, ".PROFILES")
}

func defaultAuthConfig() AuthConfig {
	secret := "CHANGE_THIS_ASAP" //TODO ref
	adminSecret := secret
	if rand, err := uuid.NewRandom(); err == nil {
		secret = rand.String()
	}
	if rand, err := uuid.NewRandom(); err == nil {
		adminSecret = rand.String()
	}

	return AuthConfig{
		TokensExpiresAfter:  time.Hour,
		TokenExpireTolerant: time.Second * 30,
		Secret:              secret,
		AdminSecret:         adminSecret,
	}
}
