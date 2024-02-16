package config

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/kmg7/fson/internal/logger"
)

// Sets up application config directory.
func (c *Config) setAppConfigDir() error {
	var configDir string
	if !c.debugMode {
		userConfig, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("user config directory not reachable err: %w", err)
		}
		configDir = userConfig

	} else {
		configDir = c.JoinDebugDir("config")
	}
	p := path.Join(configDir, APP_NAME)
	if err := os.MkdirAll(p, 0740); err != nil {
		return fmt.Errorf("cannot create app config dir err: %w", err)
	}
	c.configDir = p
	return nil
}

// Sets up application directory in users home.
func (c *Config) setAppLibDir() error {
	var homeDir string
	if !c.debugMode {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot read app lib err: %w", err)
		}
		homeDir = userHome
	} else {
		homeDir = c.JoinDebugDir("home")
	}
	p := path.Join(homeDir, "fson")
	if err := os.MkdirAll(p, 0740); err != nil {
		return fmt.Errorf("cannot create application lib dir err: %w", err)
	}
	c.libDir = p
	return nil
}

// Sets up logs directory in applications lib dir which located in users home.
func (c *Config) setAppLogsDir() error {
	ld := c.JoinLibDir("logs")
	cld := path.Join(ld, "config")
	if err := os.MkdirAll(cld, 0740); err != nil {
		return fmt.Errorf("cannot create logs dir  err: %w", err)
	}
	c.logsDir = ld
	return nil
}

// Returns temp directory on current working directory.
// Useful when observing changes during applications development.
// Also protects developers configurations
func (c *Config) setDebugDir() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get current working directory err: %w", err)
	}
	c.debugDir = path.Join(wd, "tmp")

	if err := os.MkdirAll(c.debugDir, 0740); err != nil {
		return fmt.Errorf("cannot create temp debug dir err: %w", err)
	}
	return nil
}

// Sets up a logger for config only. Logs will created under config folder.
func (c *Config) setupLogger() error {
	now := time.Now()
	name := fmt.Sprintf("%d-%d-%d-%d.log", now.Day(), now.Month(), now.Year(), now.Hour())
	p := c.JoinLogDir(path.Join("config", name))
	file, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0740)
	if err != nil {
		return fmt.Errorf("cannot create file for logging err: %w", err)
	}
	c.log = logger.New(logger.Options{Files: []io.Writer{file}, Stdout: os.Stdout})
	return nil
}

//TODO consider join bulk dirs like args

func (c *Config) JoinConfigDir(name string) string {
	return path.Join(c.configDir, name)
}

func (c *Config) JoinLibDir(name string) string {
	return path.Join(c.libDir, name)
}

func (c *Config) JoinLogDir(name string) string {
	return path.Join(c.logsDir, name)
}

func (c *Config) JoinDebugDir(name string) string {
	return path.Join(c.debugDir, name)
}
