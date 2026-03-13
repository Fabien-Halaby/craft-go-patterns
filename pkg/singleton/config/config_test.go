package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Fabien-Halaby/craft-go-patterns/pkg/singleton/config/source"
)

type TestConfig struct {
	Name  string `env:"NAME" default:"default"`
	Count int    `env:"COUNT" default:"0"`
	Live  bool   `env:"LIVE" default:"false"`
}

func (c TestConfig) Validate() error {
	if c.Count < 0 {
		return fmt.Errorf("count must be positive")
	}
	return nil
}

func TestEnvSource(t *testing.T) {
	ResetForTest[TestConfig]()
	
	os.Setenv("TEST_NAME", "prod")
	os.Setenv("TEST_COUNT", "42")
	os.Setenv("TEST_LIVE", "true")
	defer os.Unsetenv("TEST_NAME")
	defer os.Unsetenv("TEST_COUNT")
	defer os.Unsetenv("TEST_LIVE")

	cfg, err := GetInstance(source.NewEnv[TestConfig]("TEST_"))
	require.NoError(t, err)

	c := cfg.Get()
	assert.Equal(t, "prod", c.Name)
	assert.Equal(t, 42, c.Count)
	assert.True(t, c.Live)
}

func TestValidationFailure(t *testing.T) {
	ResetForTest[TestConfig]()
	
	os.Setenv("TEST_COUNT", "-5")
	defer os.Unsetenv("TEST_COUNT")

	_, err := GetInstance(source.NewEnv[TestConfig]("TEST_"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestSingleton(t *testing.T) {
	ResetForTest[TestConfig]()
	
	os.Setenv("TEST_NAME", "first")
	defer os.Unsetenv("TEST_NAME")

	cfg1, _ := GetInstance(source.NewEnv[TestConfig]("TEST_"))
	cfg2, _ := GetInstance(source.NewEnv[TestConfig]("TEST_"))

	assert.Same(t, cfg1, cfg2)
	assert.Equal(t, "first", cfg2.Get().Name)
}

func TestOnChange(t *testing.T) {
	ResetForTest[TestConfig]()
	
	os.Setenv("TEST_NAME", "initial")
	defer os.Unsetenv("TEST_NAME")

	cfg, _ := GetInstance(source.NewEnv[TestConfig]("TEST_"))
	
	var called bool
	var oldName, newName string
	
	unsub := cfg.OnChange(func(old, new TestConfig) {
		called = true
		oldName = old.Name
		newName = new.Name
	})
	defer unsub()

	os.Setenv("TEST_NAME", "reloaded")
	cfg.Reload()

	assert.True(t, called)
	assert.Equal(t, "initial", oldName)
	assert.Equal(t, "reloaded", newName)
}