package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type database struct {
	Driver   string `env:"DB_DRIVER"   envDefault:"sqlite"`
	Host     string `env:"DB_HOST"     envDefault:"127.0.0.1"`
	Port     int    `env:"DB_PORT"     envDefault:"3306"`
	Name     string `env:"DB_NAME"     envDefault:"haste_db"`
	Username string `env:"DB_USERNAME" envDefault:"root"`
	Password string `env:"DB_PASSWORD" envDefault:"root"`
}

func Database() *database {
	cfg := &database{}

	err := env.Parse(cfg)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	return cfg
}
