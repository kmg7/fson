package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigInsance(t *testing.T) {
	os.Setenv("FSON_DEBUG", "1")

	cfg := Instance()

	for i := 0; i < 10; i++ {
		go func() {
			o := Instance()
			assert.Equal(t, cfg, o)
		}()
	}
}

func TestConfigInitializer(t *testing.T) {
	os.Setenv("FSON_DEBUG", "1")

	t.Run("InitAlready", func(t *testing.T) {
		conf := Config{
			init: true,
		}
		err := conf.initialize()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "already init")
	})

	cfg := Config{}
	err := cfg.initialize()
	assert.Nil(t, err)
	assert.Equal(t, true, cfg.init)
	assert.NotNil(t, cfg.fa)
}
