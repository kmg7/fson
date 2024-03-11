package config

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kmg7/fson/internal/state"
)

type AuthConfig struct {
	filePath            string        `json:"-"`
	UpdatedAt           time.Time     `json:"updatedAt"`
	CreatedAt           time.Time     `json:"createdAt"`
	TokensExpiresAfter  time.Duration `json:"tokensExpiresAfter"`
	TokenExpireTolerant time.Duration `json:"tokenExpireTolerant"`
	Secret              string        `json:"secret"`
	AdminSecret         string        `json:"adminSecret"`
}

type AuthConfigUpdate struct {
	Expire         *time.Duration
	ExpireTolerant *time.Duration
}

func (ac *AuthConfig) FilePath() string {
	return ac.filePath
}

// Updates current config with given config update.
// Updates config if values are not same.
func (ac *AuthConfig) UpdateFrom(uc *AuthConfigUpdate) *AuthConfig {
	if uc == nil {
		return nil
	}
	nc := *ac
	updated := false

	if uc.Expire != nil {
		if nc.TokensExpiresAfter != *uc.Expire {
			nc.TokensExpiresAfter = *uc.Expire
			updated = true
		}
	}

	if uc.ExpireTolerant != nil {
		if nc.TokenExpireTolerant != *uc.ExpireTolerant {
			nc.TokenExpireTolerant = *uc.ExpireTolerant
			updated = true
		}
	}

	if updated {
		nc.UpdatedAt = time.Now()
		return &nc
	}
	return nil
}

// Reads config from file.
func (c *Config) readAuth() error {
	read := &AuthConfig{
		filePath: c.acfgPath,
	}

	if err := c.fa.ReadAndParse(read); err != nil {
		return fmt.Errorf("reading auth config failed err: %w", err)
	}
	c.acfg = read

	return nil
}

// Setup auth config with defaults.
func (c *Config) setupAuth() error {
	if !c.init {
		return ErrConfigNotInit
	}

	secret := "CHANGE_THIS_ASAP" // TODO ref
	adminSecret := secret
	if rand, err := uuid.NewRandom(); err == nil {
		secret = rand.String()
	}
	if rand, err := uuid.NewRandom(); err == nil {
		adminSecret = rand.String()
	}

	t := time.Now()

	nc := &AuthConfig{
		CreatedAt:           t,
		UpdatedAt:           t,
		filePath:            c.acfgPath,
		TokensExpiresAfter:  time.Hour,
		TokenExpireTolerant: time.Second * 30,
		Secret:              secret,
		AdminSecret:         adminSecret,
	}

	if err := c.saveAuth(nc); err != nil {
		return fmt.Errorf("while auth setup: %w", err)
	}

	return nil
}

// Save current auth config. If any error happens acfg will not be nil for sure.
func (c *Config) saveAuth(u *AuthConfig) error {
	if err := c.fa.ParseAndSave(u); err != nil {
		return fmt.Errorf("saving auth config failed err: %w", err)
	}
	c.acfg = u
	return nil
}

// Get Auth Config.
func (c *Config) GetAuth() *AuthConfig {
	r := *c.acfg
	return &r
}

// If pointer to update is nil return State with StatusFailUpdate.
func (c *Config) UpdateAuth(u *AuthConfigUpdate) (*AuthConfig, state.Error) {
	c.acfgMutex.Lock()
	defer c.acfgMutex.Unlock()

	update := c.acfg.UpdateFrom(u)
	if update != nil {
		if err := c.saveAuth(update); err != nil {
			return nil, &state.ErrInternal{
				While: "updating auth config",
				Err:   err,
			}
		}
	}

	r := *c.acfg
	return &r, nil
}
