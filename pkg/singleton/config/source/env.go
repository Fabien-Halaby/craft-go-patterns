package source

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

/**
	EnvSource: 
		- loads config from environment variables with prefix
**/
type EnvSource[T any] struct {
	Prefix string
}

func NewEnv[T any](prefix string) *EnvSource[T] {
	return &EnvSource[T]{Prefix: prefix}
}

func (s *EnvSource[T]) Load() (T, error) {
	var cfg T
	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(&cfg).Elem()

	if t.Kind() != reflect.Struct {
		return cfg, fmt.Errorf("config must be a struct, got %v", t.Kind())
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envKey := s.Prefix + field.Name
		
		if tag := field.Tag.Get("env"); tag != "" {
			envKey = s.Prefix + tag
		}

		envVal := os.Getenv(strings.ToUpper(envKey))
		if envVal == "" {
			if tag := field.Tag.Get("default"); tag != "" {
				envVal = tag
			} else {
				continue
			}
		}

		if err := s.setField(v.Field(i), envVal); err != nil {
			return cfg, fmt.Errorf("field %s: %w", field.Name, err)
		}
	}

	return cfg, nil
}

func (s *EnvSource[T]) Watch(fn func(T)) error {
	// Env vars don't support watch natively
	// Could implement with fsnotify on .env file
	return nil
}

func (s *EnvSource[T]) setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	default:
		return fmt.Errorf("unsupported type: %v", field.Kind())
	}
	return nil
}