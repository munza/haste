package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type app struct {
	Name      string `env:"APP_NAME"   envDefault:"Haste"`
	URL       string `env:"APP_URL"    envDefault:"127.0.0.1"`
	SecretKey string `env:"SECRET_KEY" envDefault:"somesupersecretkey"`
}

func App() *app {
	cfg := &app{}

	err := env.Parse(cfg)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	return cfg
}
