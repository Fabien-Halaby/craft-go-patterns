package main

import (
	"fmt"
	"log"

	"github.com/Fabien-Halaby/craft-go-patterns/pkg/singleton/config"
	"github.com/Fabien-Halaby/craft-go-patterns/pkg/singleton/config/source"
)

type AppConfig struct {
	DatabaseURL string `env:"DATABASE_URL" default:"postgres://localhost/db"`
	Port        int    `env:"PORT" default:"8080"`
	Debug       bool   `env:"DEBUG" default:"false"`
}

func main() {
	cfg, err := config.GetInstance(source.NewEnv[AppConfig]("APP_"))
	if err != nil {
		log.Fatal(err)
	}

	c := cfg.Get()
	fmt.Printf("DB: %s\n", c.DatabaseURL)
	fmt.Printf("Port: %d\n", c.Port)
	fmt.Printf("Debug: %v\n", c.Debug)

	cfg.OnChange(func(old, new AppConfig) {
		log.Printf("Config changed! Port: %d -> %d", old.Port, new.Port)
	})
}