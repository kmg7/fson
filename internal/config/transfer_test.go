package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransferConfig(t *testing.T) {
	started := time.Now()
	tcfg := TransferConfig{
		filePath:  "path",
		AutoStart: true,
		CreatedAt: started,
		UpdatedAt: started,
		UploadDir: "oldUploadDir",
		TempDir:   "oldTempDir",
		Transfer: []TransferPath{
			{Id: "old1", Path: "oldPath1"},
			{Id: "old2", Path: "oldPath2"},
			{Id: "old3", Path: "oldPath3"},
		},
	}
	t.Run("MethodFilePath", func(t *testing.T) {
		assert.Equal(t, tcfg.filePath, tcfg.FilePath())
	})
	t.Run("MethodUpdateFrom", func(t *testing.T) {
		uAutoStart := false
		uDir := "uploaded"
		uTransfer := []TransferPath{
			{Id: "new", Path: "updated"},
		}
		update := TransferConfigUpdate{
			AutoStart: &uAutoStart,
			UploadDir: &uDir,
			TempDir:   &uDir,
			Transfer:  &uTransfer,
		}
		t.Run("UpdateWithChanges", func(t *testing.T) {
			got := tcfg.UpdateFrom(&update)
			assert.NotNil(t, got)
			assert.NotEqual(t, tcfg.UpdatedAt, got.UpdatedAt)
			assert.Equal(t, tcfg.CreatedAt, got.CreatedAt)
		})

		t.Run("UpdateWithNoChanges", func(t *testing.T) {
			got := tcfg.UpdateFrom(nil)
			assert.Nil(t, got)
		})

		t.Run("UpdateWithPartialChanges", func(t *testing.T) {
			got := tcfg.UpdateFrom(&TransferConfigUpdate{
				UploadDir: &uDir,
			})
			assert.NotNil(t, got)
			assert.NotEqual(t, tcfg.UploadDir, got.UploadDir)
			assert.Equal(t, tcfg.AutoStart, got.AutoStart)
			assert.Equal(t, tcfg.TempDir, got.TempDir)
			assert.Equal(t, tcfg.Transfer, got.Transfer)
		})
	})
}

func TestConfigTranferIntegr(t *testing.T) {
	notInitCfg := &Config{init: false}
	os.Setenv("FSON_DEBUG", "1")
	cfg := Instance()
	defer func() {
		if err := os.RemoveAll(cfg.debugDir); err != nil {
			t.Logf("Removing test dir failed err: %v path: %v", err.Error(), cfg.debugDir)
		}
	}()

	t.Run("NotInitConfig", func(t *testing.T) {
		assert.Nil(t, notInitCfg.tcfg)
		assert.NoFileExists(t, cfg.tcfgPath)
	})

	t.Run("NotExistingRead", func(t *testing.T) {
		// On first run read will raise error
		err := cfg.readTranfer()
		assert.NotNil(t, err)
		assert.Nil(t, cfg.tcfg)
	})

	t.Run("Setup", func(t *testing.T) {
		err := cfg.setupTransfer()
		assert.Nil(t, err)
		assert.NotNil(t, cfg.tcfg)
	})

	t.Run("EqualPaths", func(t *testing.T) {
		assert.Equal(t, cfg.tcfgPath, cfg.tcfg.FilePath())
	})

	t.Run("ExistingRead", func(t *testing.T) {
		cfg.tcfg = nil
		err := cfg.readTranfer()
		assert.Nil(t, err)
		assert.NotNil(t, cfg.tcfg)
	})

	t.Run("Saving", func(t *testing.T) {
		t.Run("GibberishPath", func(t *testing.T) {
			save := *cfg.tcfg
			save.filePath = "/someGibberishPath/file.cfg"
			err := cfg.saveTransfer(&save)
			assert.NotNil(t, err)
			assert.NotNil(t, cfg.tcfg)
		})

		oldField := cfg.tcfg.AutoStart
		newField := !oldField
		cfg.tcfg.AutoStart = newField

		t.Run("FinePath", func(t *testing.T) {
			cfg.tcfg.filePath = cfg.tcfgPath
			save := *cfg.tcfg
			err := cfg.saveTransfer(&save)
			assert.Nil(t, err)
			assert.NotNil(t, cfg.tcfg)
			assert.NotEqual(t, oldField, cfg.tcfg.AutoStart)
		})

		t.Run("ReadAfterSaving", func(t *testing.T) {
			err := cfg.readTranfer()
			assert.Nil(t, err)
			assert.NotNil(t, cfg.tcfg)
			assert.Equal(t, newField, cfg.tcfg.AutoStart)
		})

	})
	t.Run("Get", func(t *testing.T) {
		state := cfg.GetTransfer()
		assert.NotNil(t, state)
		assert.Equal(t, StatusSuccess, state.Status)
		assert.Equal(t, cfg.tcfg, state.Data)
	})
	t.Run("UpdateEmpty", func(t *testing.T) {
		state := cfg.UpdateTransfer(nil)
		assert.NotNil(t, state)
		assert.Equal(t, StatusFailUpdate, state.Status)
	})
	t.Run("Update", func(t *testing.T) {
		uAutoStart := false
		uDir := "uploaded"
		uTransfer := []TransferPath{
			{Id: "new", Path: "updated"},
		}
		update := TransferConfigUpdate{
			AutoStart: &uAutoStart,
			UploadDir: &uDir,
			TempDir:   nil,
			Transfer:  &uTransfer,
		}
		t.Run("InternalError", func(t *testing.T) {
			cfg.tcfg.filePath = ""
			state := cfg.UpdateTransfer(&update)
			assert.NotNil(t, state)
			assert.NotNil(t, state.Err)
			assert.Equal(t, true, state.ErrInternal)
			assert.Equal(t, StatusFailInternal, state.Status)

		})
		cfg.tcfg.filePath = cfg.tcfgPath
		state := cfg.UpdateTransfer(&update)
		assert.NotNil(t, state)

		assert.Equal(t, StatusUpdate, state.Status)
		assert.Equal(t, state.Data, cfg.tcfg)

	})

}
