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

var authCfgPath string

// Initalizes auth config on first run
func initAuthConfig() *AuthConfig {
	p := authConfigFilePath()
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

// Reads auth config if an error happens returns
func readAuthConfig() (*AuthConfig, error) {
	conf := &AuthConfig{}
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

	err = json.Unmarshal(data, conf)
	if err != nil {
		logger.Error("Cannot parse auth config", err)
	}
	logger.Info("Auth config read")
	return conf, err
}

func writeAuthConfig(config *AuthConfig) error {
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

func authConfigFilePath() string {
	appCfg := appConfigDir()
	return path.Join(appCfg, ".AUTH_CFG")
}

func AuthProfilesFilePath() string {
	appCfg := appConfigDir()
	return path.Join(appCfg, ".PROFILES")
}

func defaultAuthConfig() *AuthConfig {
	secret := "CHANGE_THIS_ASAP" //TODO ref
	adminSecret := secret
	if rand, err := uuid.NewRandom(); err == nil {
		secret = rand.String()
	}
	if rand, err := uuid.NewRandom(); err == nil {
		adminSecret = rand.String()
	}

	return &AuthConfig{
		TokensExpiresAfter:  time.Hour,
		TokenExpireTolerant: time.Second * 30,
		Secret:              secret,
		AdminSecret:         adminSecret,
	}
}
