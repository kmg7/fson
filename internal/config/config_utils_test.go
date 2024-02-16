package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupsDebug(t *testing.T) {
	c := &Config{
		debugMode: true,
	}
	wd, err := os.Getwd()
	assert.Nil(t, err)

	err = c.setDebugDir()

	assert.Nil(t, err)
	assert.DirExists(t, c.debugDir)

	assert.Equal(t, c.debugDir, path.Join(wd, "tmp"))

	// defer func() {
	// 	if err := os.RemoveAll(c.debugDir); err != nil {
	// 		t.Logf("Failed to delete temporary files. Path: %v", c.debugDir)
	// 	}
	// }()
	t.Run("AppConfig", func(t *testing.T) {
		err = c.setAppConfigDir()
		assert.Nil(t, err)
		assert.DirExists(t, c.configDir)
		assert.Contains(t, c.configDir, c.debugDir)
	})
	t.Run("AppLib", func(t *testing.T) {
		err = c.setAppLibDir()
		assert.Nil(t, err)
		assert.DirExists(t, c.libDir)
		assert.Contains(t, c.libDir, c.debugDir)
	})
	t.Run("AppLogs", func(t *testing.T) {
		err = c.setAppLogsDir()
		assert.Nil(t, err)
		assert.DirExists(t, c.logsDir)
		assert.Contains(t, c.logsDir, c.debugDir)
	})
	t.Run("AppLogger", func(t *testing.T) {
		err = c.setupLogger()
		assert.Nil(t, err)
		assert.DirExists(t, c.JoinLogDir("config"))
	})
}

func TestSetups(t *testing.T) {
	c := &Config{}
	userConfigDir, err := os.UserConfigDir()
	assert.Nil(t, err)

	userHomeDir, err := os.UserHomeDir()
	assert.Nil(t, err)

	t.Run("AppConfig", func(t *testing.T) {
		err = c.setAppConfigDir()
		assert.Nil(t, err)
		assert.DirExists(t, c.configDir)
		assert.Contains(t, c.configDir, userConfigDir)
	})
	t.Run("AppLib", func(t *testing.T) {
		err = c.setAppLibDir()
		assert.Nil(t, err)
		assert.DirExists(t, c.libDir)
		assert.Contains(t, c.libDir, userHomeDir)
	})
	t.Run("AppLogs", func(t *testing.T) {
		err = c.setAppLogsDir()
		assert.Nil(t, err)
		assert.DirExists(t, c.logsDir)
		assert.Contains(t, c.logsDir, userHomeDir)
	})
	t.Run("AppLogger", func(t *testing.T) {
		err = c.setupLogger()
		assert.Nil(t, err)
		assert.DirExists(t, c.JoinLogDir("config"))
	})
}
