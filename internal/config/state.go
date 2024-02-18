package config

import (
	"errors"

	"github.com/kmg7/fson/internal/state"
)

type ConfigState = state.AppState

var (
	ErrConfigNotInit = errors.New("config not initialized yet")
)

const (
	StatusSuccess      = "success"
	StatusFailInternal = "internal-fail"
	StatusUpdate       = "success-update"
	StatusFailUpdate   = "update-fail"
)
