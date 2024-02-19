package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthConfig(t *testing.T) {
	started := time.Now()
	acfg := AuthConfig{
		filePath:            "path",
		CreatedAt:           started,
		UpdatedAt:           started,
		TokensExpiresAfter:  time.Hour,
		TokenExpireTolerant: time.Minute,
	}

	t.Run("MethodFilePath", func(t *testing.T) {
		assert.Equal(t, acfg.filePath, acfg.FilePath())
	})

	t.Run("MethodUpdateFrom", func(t *testing.T) {
		uExpire := time.Hour * 2
		uExpireTolerant := time.Minute * 2
		update := AuthConfigUpdate{
			Expire: &uExpire,
			ExpireTolerant: &uExpireTolerant,
		}

		t.Run("UpdateWithChanges", func(t *testing.T) {
			got := acfg.UpdateFrom(&update)
			assert.NotNil(t, got)
			assert.NotEqual(t, acfg.UpdatedAt, got.UpdatedAt)
			assert.Equal(t, acfg.CreatedAt, got.CreatedAt)
		})

		t.Run("UpdateWithNoChanges", func(t *testing.T) {
			got := acfg.UpdateFrom(nil)
			assert.Nil(t, got)
		})

		t.Run("UpdateWithPartialChanges", func(t *testing.T) {
			got := acfg.UpdateFrom(&AuthConfigUpdate{
		Expire: &uExpire,	
			})
			assert.NotNil(t, got)
			assert.NotEqual(t, acfg.TokensExpiresAfter, got.TokensExpiresAfter)
			assert.Equal(t, acfg.TokenExpireTolerant, got.TokenExpireTolerant)
		})
	})
}

func TestConfigAuthIntegr(t *testing.T) {
	notInitCfg := &Config{init: false}
	os.Setenv("FSON_DEBUG", "1")
	cfg := Instance()
	defer func() {
		if err := os.RemoveAll(cfg.debugDir); err != nil {
			t.Logf("Removing test dir failed err: %v path: %v", err.Error(), cfg.debugDir)
		}
	}()

	t.Run("NotInitConfig", func(t *testing.T) {
		assert.Nil(t, notInitCfg.acfg)
		assert.NoFileExists(t, cfg.acfgPath)
	})

	t.Run("NotExistingRead", func(t *testing.T) {
		// On first run read will raise error
		err := cfg.readAuth()
		assert.NotNil(t, err)
		assert.Nil(t, cfg.acfg)
	})

	t.Run("Setup", func(t *testing.T) {
		err := cfg.setupAuth()
		assert.Nil(t, err)
		assert.NotNil(t, cfg.acfg)
	})

	t.Run("EqualPaths", func(t *testing.T) {
		assert.Equal(t, cfg.acfgPath, cfg.acfg.FilePath())
	})

	t.Run("ExistingRead", func(t *testing.T) {
		cfg.acfg = nil
		err := cfg.readAuth()
		assert.Nil(t, err)
		assert.NotNil(t, cfg.acfg)
	})

	t.Run("Saving", func(t *testing.T) {
		t.Run("GibberishPath", func(t *testing.T) {
			save := *cfg.acfg
			save.filePath = "/someGibberishPath/file.cfg"
			err := cfg.saveAuth(&save)
			assert.NotNil(t, err)
			assert.NotNil(t, cfg.acfg)
		})

		oldField := cfg.acfg.TokensExpiresAfter
		newField := oldField * 2
		cfg.acfg.TokensExpiresAfter= newField

		t.Run("FinePath", func(t *testing.T) {
			cfg.acfg.filePath = cfg.acfgPath
			save := *cfg.acfg
			err := cfg.saveAuth(&save)
			assert.Nil(t, err)
			assert.NotNil(t, cfg.acfg)
			assert.NotEqual(t, oldField, cfg.acfg.TokensExpiresAfter)
		})

		t.Run("ReadAfterSaving", func(t *testing.T) {
			err := cfg.readAuth()
			assert.Nil(t, err)
			assert.NotNil(t, cfg.acfg)
			assert.Equal(t, newField, cfg.acfg.TokensExpiresAfter)
		})
	})
	t.Run("Get", func(t *testing.T) {
		state := cfg.GetAuth()
		assert.NotNil(t, state)
		assert.Equal(t, StatusSuccess, state.Status)
		assert.Equal(t, cfg.acfg, state.Data)
	})
	t.Run("UpdateEmpty", func(t *testing.T) {
		state := cfg.UpdateAuth(nil)
		assert.NotNil(t, state)
		assert.Equal(t, StatusFailUpdate, state.Status)
	})
	t.Run("Update", func(t *testing.T) {
		uExpire := time.Hour* 3
		update := AuthConfigUpdate{
		Expire: &uExpire,
		ExpireTolerant: nil,
		}
		t.Run("InternalError", func(t *testing.T) {
			cfg.acfg.filePath = ""
			state := cfg.UpdateAuth(&update)
			assert.NotNil(t, state)
			assert.NotNil(t, state.Err)
			assert.Equal(t, true, state.ErrInternal)
			assert.Equal(t, StatusFailInternal, state.Status)
		})
		cfg.acfg.filePath = cfg.acfgPath
		state := cfg.UpdateAuth(&update)
		assert.NotNil(t, state)

		assert.Equal(t, StatusUpdate, state.Status)
		assert.Equal(t, state.Data, cfg.acfg)
	})
}
