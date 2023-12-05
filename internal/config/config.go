package config

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kmg7/fson/internal/logger"
	"github.com/kmg7/fson/pkg/fsutils"
)

const debugMode bool = true

var appCfg AppConfig
var transferCfg TransferConfig

func Init() {
	initializeLogger()
	appCfg = initAppConfig()
	transferCfg = initTransferConfig()
}

func GetAppCfg() *AppConfig {
	return &appCfg
}

func GetTransferCfg() *TransferConfig {
	return &transferCfg
}

func UpdateAppConfig(autoStart bool, tempDir, uploadDir string) *ConfigError {
	cfg := appCfg
	cfg.AutoStart = autoStart
	cfg.TempDir = tempDir
	cfg.UploadDir = uploadDir
	cfg.UpdatedAt = time.Now()
	if err := setAppConfig(cfg); err != nil {
		return &ConfigError{
			Internal: true, Code: "cannot-save", Messages: []string{"Cannot save app config."}, Err: err}
	}
	appCfg = cfg
	return nil
}

func AddTransferPath(path string) *ConfigError {
	cfg := transferCfg
	path = strings.TrimSuffix(strings.TrimSpace(path), "/")
	for _, tp := range cfg.Transfer {
		if tp.Path == path {
			return &ConfigError{
				Code: "already-exists", Messages: []string{"Given transfer path already exists."}}
		}
	}
	if ex, _ := fsutils.Exists(path); !ex {
		return &ConfigError{Code: "broken-path", Messages: []string{"Given trasnfer path broken or not exists in sytem."}}
	}

	id, err := uuid.NewRandom()

	if err != nil {
		return &ConfigError{Internal: true, Code: "unexpected-internal", Messages: []string{"Unexpected error happened try again later."}}
	}

	cfg.Transfer = append(cfg.Transfer, TransferPath{Id: id.String(), Path: path})
	cfg.UpdatedAt = time.Now()
	if err := setTransferConfig(cfg); err != nil {
		return &ConfigError{Internal: true, Code: "cannot-save", Messages: []string{"Cannot save transfer config"}, Err: err}
	}
	transferCfg = cfg
	return nil
}

func UpdateTransferPath(id, path string) *ConfigError {
	cfg := transferCfg
	found := false
	for i, tp := range cfg.Transfer {
		if tp.Id == id {
			path = strings.TrimSuffix(strings.TrimSpace(path), "/")
			if tp.Path == path {
				return nil
			}
			if ex, err := fsutils.Exists(path); err != nil {
				return &ConfigError{Internal: true, Code: "unexpected-internal", Messages: []string{"Unexpected error happened try again later."}}
			} else {
				if ex {
					cfg.Transfer[i] = TransferPath{Id: id, Path: path}
				} else {
					return &ConfigError{Code: "path-not-exist", Messages: []string{"Given transfer path not exists"}}
				}
			}
			found = true
			break
		}
	}

	if !found {
		return &ConfigError{Code: "not-exists", Messages: []string{"Transfer path not not found"}}
	}

	cfg.UpdatedAt = time.Now()
	if err := setTransferConfig(cfg); err != nil {
		return &ConfigError{Internal: true, Code: "cannot-save", Messages: []string{"Cannot save transfer config"}, Err: err}
	}
	transferCfg = cfg
	return nil

}

func DeleteTransferPath(id string) *ConfigError {
	cfg := transferCfg
	deleted := false
	for i, tp := range cfg.Transfer {
		if tp.Id == id {
			cfg.Transfer = append(cfg.Transfer[:i], cfg.Transfer[i+1:]...)
			deleted = true
			break
		}
	}
	if !deleted {
		return &ConfigError{Code: "not-exists", Messages: []string{"Transfer path not not found"}}
	}
	cfg.UpdatedAt = time.Now()
	if err := setTransferConfig(cfg); err != nil {
		return &ConfigError{Internal: true, Code: "cannot-save", Messages: []string{"Cannot save transfer config"}, Err: err}
	}
	transferCfg = cfg
	return nil
}

func initializeLogger() {
	opt := logger.Options{Development: debugMode}
	logDir, err := appLogsDir()
	if err != nil {
		//*logger failed
		return
	}
	now := time.Now()
	file := fmt.Sprintf("%d-%d-%d-%d.log", now.Day(), now.Month(), now.Year(), now.Hour())
	fileP := path.Join(logDir, file)

	if err := fsutils.TouchAll(fileP); err != nil {
		fmt.Println("Failed to create log files")
	} else {
		if f, err := os.OpenFile(fileP, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			opt.Output = f
		}
	}
	logger.InitLogger(opt)

}

type ConfigError struct {
	Internal bool
	Code     string
	Messages []string
	Err      error
}
