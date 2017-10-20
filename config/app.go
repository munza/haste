package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type app struct {
	Environtment string `env:"APP_ENV"    envDefault:"development"`
	Name         string `env:"APP_NAME"   envDefault:"Haste"`
	URL          string `env:"APP_URL"    envDefault:"http://127.0.0.1"`
	Port         int    `env:"APP_PORT"   envDefault:"3000"`
	SecretKey    string `env:"SECRET_KEY" envDefault:"somesupersecretkey"`
}

func App() *app {
	cfg := &app{}

	err := env.Parse(cfg)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	return cfg
}
