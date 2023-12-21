package config

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/kmg7/fson/internal/logger"
	"github.com/kmg7/fson/pkg/fsutils"
)

// For configuring application
type AppConfig struct {
	Password  string    `json:"password"`
	AutoStart bool      `json:"autoStart"`
	UploadDir string    `json:"uploadDir"`
	TempDir   string    `json:"tempDir"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// The path of Application Config File.
var appCfgPath string

// Changes app config with provided config.
// Then attempts to save it.
// Return an error if something happens during the process.
func setAppConfig(newCfg AppConfig) error {
	if err := writeAppConfig(&newCfg); err != nil {
		return err
	}
	return nil
}

// Initializes applications config file.
func initAppConfig() *AppConfig {
	p := appConfigFile()
	appCfgPath = p
	ex, err := fsutils.Exists(appCfgPath)
	if err != nil {
		logger.Error(err)
	}
	if !ex {
		logger.Info("Not found any config file")
		confd := defaultAppConfig()
		writeAppConfig(confd)
		return confd
	}
	conf, err := readAppConfig()
	if err != nil {
		confd := defaultAppConfig()
		logger.Info("Overriding broken config")
		writeAppConfig(confd)
		return confd
	}
	return conf
}

// Reads app config file then returns config and an error.
// If any error occurs on reading logs and returns it, otherwise returns nil
// and logs an info message
func readAppConfig() (*AppConfig, error) {
	conf := &AppConfig{}
	file, err := os.Open(appCfgPath)
	if err != nil {
		logger.Error("Cannot open config gile")
		return conf, err
	}

	defer file.Close()
	data, err := os.ReadFile(appCfgPath)
	if err != nil {
		logger.Error("Cannot read config file", err)
		return conf, err
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		logger.Error("Cannot parse config file", err)
	}
	logger.Info("App config read")
	return conf, err
}

// Writes app config file then logs and returns if any error occurs.
// If not returns nil and logs an info message.
func writeAppConfig(config *AppConfig) error {
	file, err := os.Create(appCfgPath)
	if err != nil {
		logger.Error("Cannot create config file", err)
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		logger.Error("Cannot parse config", err)
		return err
	}

	if _, err := file.Write(data); err != nil {
		logger.Error("Cannot write config file", err)
		return err
	}
	logger.Info("App config wrote")
	return nil
}

// Returns default config and logs an info message.
// If any error occurs on reading appLibDir()
// logs the error and panics.
func defaultAppConfig() *AppConfig {
	lib := appLibDir()
	updir := path.Join(lib, "uploaded")
	tmpDir := path.Join(lib, "temp")
	if err := os.MkdirAll(updir, 0755); err != nil {
		logger.Fatal("Cannot create uploaded dir", err.Error())
	}
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		logger.Fatal("Cannot create uploaded dir", err.Error())
	}
	logger.Info("Creating new default app config")
	return &AppConfig{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Password:  "admin",
		AutoStart: false,
		UploadDir: path.Join(lib, "uploaded"),
		TempDir:   path.Join(lib, "temp"),
	}
}

// Returns application library path in the user home directory.
// If the path doesn't exists creates with 0755.
func appLibDir() string {
	var homeDir string
	if !debugMode {
		userHome, err := os.UserHomeDir()
		if err != nil {
			logger.Fatal("Cannot read app lib", err.Error())
			return homeDir
		}
		homeDir = userHome
	} else {
		homeDir = path.Join(debugTempDir(), "home")
	}
	p := path.Join(homeDir, "fson")
	if err := os.MkdirAll(p, 0755); err != nil {
		logger.Fatal("Cannot create go")
	}

	return p
}

// Returns application library path in the user home directory.
// If the path doesn't exists creates with 0755.
func appConfigDir() string {
	var configDir string
	if !debugMode {
		userConfig, err := os.UserConfigDir()
		if err != nil {
			logger.Fatal("User config directory not reachable ", err.Error())
		}
		configDir = userConfig

	} else {
		configDir = path.Join(debugTempDir(), "config")
	}
	p := path.Join(configDir, "fson")
	if err := os.MkdirAll(p, 0755); err != nil {
		logger.Fatal("Cannot create app config dir ", err.Error())
	}
	return p
}

// Returns directory of the logs creates if it not exists.
func appLogsDir() string {
	var logsDir string
	appLib := appLibDir()
	logsDir = path.Join(appLib, "logs")

	if err := os.MkdirAll(logsDir, 0755); err != nil {
		logger.Fatal("Cannot create logs dir ", err.Error())
	}
	return logsDir
}

// Config file name is .APP_CFG.
func appConfigFile() string {
	appCfg := appConfigDir()

	return path.Join(appCfg, ".APP_CFG")
}

// Returns temp directory of current working process.
// Useful when observing changes during applications development.
// Also protects developers configurations
func debugTempDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(wd, "tmp")
}
