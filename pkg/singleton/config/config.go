package config

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/Fabien-Halaby/craft-go-patterns/pkg/singleton"
)

type Validator interface {
	Validate() error
}

/**
	Config:
		- holds configuration of type T with thread-safe access.
**/
type Config[T any] struct {
	mu       sync.RWMutex
	data     T
	source   Source[T]
	onChange []func(old, new T)
	loaded   bool
}

/**
	Source:
		- loads configuration from external source
**/
type Source[T any] interface {
	Load() (T, error)
	Watch(func(T)) error
}

var instances = &sync.Map{}

/**
	GetInstance:
		- returns singleton Config for type T.
	Usage:
		- config.GetInstance[MyConfig](source)
**/
func GetInstance[T any](source Source[T]) (*Config[T], error) {
	key := reflect.TypeOf((*T)(nil)).Elem()

	actual, loaded := instances.LoadOrStore(key, singleton.NewInstance(func() *Config[T] {
		return &Config[T]{
			source:   source,
			onChange: make([]func(old, new T), 0),
		}
	}))

	cfg := actual.(*singleton.Instance[*Config[T]]).Get()

	if !loaded || !cfg.loaded {
		if err := cfg.load(); err != nil {
			return nil, fmt.Errorf("config load failed: %w", err)
		}
	}

	return cfg, nil
}

/**
	MustGetInstance:
		- panics on error (pour init simple)
**/
func MustGetInstance[T any](source Source[T]) *Config[T] {
	cfg, err := GetInstance(source)
	if err != nil {
		panic(err)
	}
	return cfg
}

func (c *Config[T]) load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := c.source.Load()
	if err != nil {
		return err
	}

	if v, ok := any(data).(Validator); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	old := c.data
	c.data = data
	c.loaded = true

	c.mu.Unlock()
	for _, fn := range c.onChange {
		fn(old, data)
	}
	c.mu.Lock()

	return nil
}

/**
	Get:
		- returns current config (copy)
**/
func (c *Config[T]) Get() T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data
}

/**
	Reload:
		- force reload from source
**/
func (c *Config[T]) Reload() error {
	return c.load()
}

/**
	OnChange:
		-registers callback for config changes
**/
func (c *Config[T]) OnChange(fn func(old, new T)) func() {
	c.mu.Lock()
	id := len(c.onChange)
	c.onChange = append(c.onChange, fn)
	c.mu.Unlock()

	return func() {
		c.mu.Lock()
		if id < len(c.onChange) {
			c.onChange[id] = nil // soft delete
		}
		c.mu.Unlock()
	}
}

/**
	ResetForTest:
		- clears singleton (testing only)
**/
func ResetForTest[T any]() {
	key := reflect.TypeOf((*T)(nil)).Elem()
	instances.Delete(key)
}