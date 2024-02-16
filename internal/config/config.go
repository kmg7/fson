package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/kmg7/fson/internal/logger"
)

const APP_NAME string = "fson"

// Config struct is designed to manipule how this applications works.
type Config struct {
	init      bool
	debugMode bool
	debugDir  string
	configDir string
	libDir    string
	logsDir   string
	fa        *FileAdapter
	log       logger.AppLogger
}

var si *Config
var instanciateOnce sync.Once

// A single initialized instance of Config
func Instance() *Config {
	instanciateOnce.Do(func() {
		//Creating config instance
		i := &Config{}
		if debugEnv := os.Getenv("FSON_DEBUG"); debugEnv != "" {
			i.debugMode = true
		}
		if err := i.initialize(); err != nil {
			log.Fatal(err.Error())
		}
		si = i
	})
	return si
}

// Initializes given config if its not init already.
// Sets up configs directories (config, lib and log).
// For debug purposes be sure to set debug mode before calling.
func (c *Config) initialize() error {
	if c.init {
		return errors.New("config already init")
	}
	if c.debugMode {
		if err := c.setDebugDir(); err != nil {
			return err
		}
	}

	if err := c.setAppConfigDir(); err != nil {
		return err
	}
	if err := c.setAppLibDir(); err != nil {
		return err
	}
	if err := c.setAppLogsDir(); err != nil {
		return err
	}
	if err := c.setupLogger(); err != nil {
		return err

	}

	c.fa = &FileAdapter{
		Parse:   json.Marshal,
		Unparse: json.Unmarshal,
	}

	c.init = true
	return nil
}
