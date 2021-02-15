// Package config provides configuration.
package config

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

var (
	cfg  Config
	once sync.Once
)

// Config is configuration for all project components.
type Config struct {
	ClientID     string `envconfig:"CLIENT_ID"`
	ClientSecret string `envconfig:"CLIENT_SECRET"`
	AuthCodeURL  string `envconfig:"AUTH_CODE_URL"`
}

// Get reads configuration once and returns it.
func Get() Config {
	once.Do(func() {
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatal("couldn't read configuration")
		}
	})

	return cfg
}
