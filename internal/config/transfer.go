package config

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/kmg7/fson/internal/logger"
	"github.com/kmg7/fson/pkg/fsutils"
)

type TransferPath struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type TransferConfig struct {
	UpdatedAt time.Time      `json:"updatedAt"`
	CreatedAt time.Time      `json:"createdAt"`
	Transfer  []TransferPath `json:"transfer"`
}

var transferCfgPath string

func setTransferConfig(newCfg TransferConfig) error {
	if err := writeTransferConfig(&newCfg); err != nil {
		return err
	}
	return nil
}

func initTransferConfig() *TransferConfig {
	p := transferConfigFile()
	transferCfgPath = p
	ex, err := fsutils.Exists(transferCfgPath)
	if err != nil {
		logger.Error("Cannot check if any config file exists", err)
	}
	if !ex {
		confd := defaultTransferConfig()
		writeTransferConfig(confd)
		return confd
	}
	conf, err := readTransferConfig()
	if err != nil {
		confd := defaultTransferConfig()
		logger.Warn("Overriding broken config", err)
		writeTransferConfig(confd)
		return confd
	}
	return conf
}

func readTransferConfig() (*TransferConfig, error) {
	conf := &TransferConfig{}
	file, err := os.Open(transferCfgPath)
	if err != nil {
		logger.Error("Cannot open transfer config", err)
		return conf, err
	}

	defer file.Close()
	data, err := os.ReadFile(transferCfgPath)
	if err != nil {
		logger.Error("Cannot read transfer config", err)
		return conf, err
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		logger.Error("Cannot parse transfer config", err)
	}
	logger.Info("Transfer config read")
	return conf, err
}

func writeTransferConfig(config *TransferConfig) error {
	file, err := os.Create(transferCfgPath)
	if err != nil {
		logger.Error("Cannot create transfer config", err)
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		logger.Error("Cannot parse transfer config", err)
		return err
	}

	if _, err := file.Write(data); err != nil {
		logger.Error("Cannot write transfer config", err)
		return err
	}
	logger.Info("Transfer config wrote")
	return nil
}

func defaultTransferConfig() *TransferConfig {
	return &TransferConfig{
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		Transfer:  []TransferPath{},
	}

}

func transferConfigFile() string {
	appCfg := appConfigDir()
	return path.Join(appCfg, ".SERVE_CFG")
}
