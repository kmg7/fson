package config

import (
	"fmt"
	"time"
)

type TransferPath struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type TransferConfig struct {
	filePath  string
	UpdatedAt time.Time      `json:"updatedAt"`
	CreatedAt time.Time      `json:"createdAt"`
	AutoStart bool           `json:"autoStart"`
	UploadDir string         `json:"uploadDir"`
	TempDir   string         `json:"tempDir"`
	Transfer  []TransferPath `json:"transfer"`
}

type TransferConfigUpdate struct {
	AutoStart *bool           `json:"autoStart"`
	UploadDir *string         `json:"uploadDir"`
	TempDir   *string         `json:"tempDir"`
	Transfer  *[]TransferPath `json:"transfer"`
}

func (tc *TransferConfig) FilePath() string {
	return tc.filePath
}

func (tc *TransferConfig) UpdateFrom(uc *TransferConfigUpdate) *TransferConfig {
	if uc == nil {
		return nil
	}
	nc := *tc
	updated := false

	if uc.AutoStart != nil {
		if nc.AutoStart != *uc.AutoStart {
			nc.AutoStart = *uc.AutoStart
			updated = true
		}
	}

	if uc.TempDir != nil {
		if nc.TempDir != *uc.TempDir {
			nc.TempDir = *uc.TempDir
			updated = true
		}
	}

	if uc.UploadDir != nil {
		if nc.UploadDir != *uc.UploadDir {
			nc.UploadDir = *uc.UploadDir
			updated = true
		}
	}

	if uc.Transfer != nil {
		nc.Transfer = *uc.Transfer
		updated = true
	}

	if updated {
		nc.UpdatedAt = time.Now()
		return &nc
	}
	return nil
}

// Reads config from file.
func (c *Config) readTranfer() error {
	read := &TransferConfig{
		filePath: c.tcfgPath,
	}

	if err := c.fa.ReadAndParse(read); err != nil {
		return fmt.Errorf("reading transfer config failed err: %w", err)
	}
	c.tcfg = read

	return nil
}

// Setup transfer config with defaults.
func (c *Config) setupTransfer() error {
	if !c.init {
		return ErrConfigNotInit
	}

	t := time.Now()
	nc := &TransferConfig{
		CreatedAt: t,
		UpdatedAt: t,
		filePath:  c.tcfgPath,
	}
	if err := c.saveTransfer(nc); err != nil {
		return fmt.Errorf("while transfer setup: %w", err)
	}
	return nil
}

// Save current transfer config. If any error happens tcfg will not be nil for sure.
func (c *Config) saveTransfer(u *TransferConfig) error {
	if err := c.fa.ParseAndSave(u); err != nil {
		return fmt.Errorf("saving transfer config failed err: %w", err)
	}
	c.tcfg = u
	return nil
}

// Get Transfer Config.
func (c *Config) GetTransfer() *ConfigState {
	return &ConfigState{
		Status: StatusSuccess,
		Data:   c.tcfg,
	}
}

// If pointer to update is nil return State with StatusFailUpdate.
func (c *Config) UpdateTransfer(u *TransferConfigUpdate) *ConfigState {
	update := c.tcfg.UpdateFrom(u)
	if update != nil {
		if err := c.saveTransfer(update); err != nil {
			return &ConfigState{
				Status:      StatusFailInternal,
				ErrInternal: true,
				Err:         err,
			}
		}
		return &ConfigState{
			Status: StatusUpdate,
			Data:   update,
		}
	}
	return &ConfigState{
		Status: StatusFailUpdate,
	}
}
