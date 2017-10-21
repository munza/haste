package config

import (
	"strconv"

	"github.com/caarlos0/env"
)

type database struct {
	Driver   string `env:"DB_DRIVER"   envDefault:"sqlite"`
	Host     string `env:"DB_HOST"     envDefault:"127.0.0.1"`
	Port     int    `env:"DB_PORT"     envDefault:"3306"`
	Name     string `env:"DB_NAME"     envDefault:"haste_db"`
	Username string `env:"DB_USERNAME" envDefault:"root"`
	Password string `env:"DB_PASSWORD" envDefault:"root"`

	MigrationPath string `envDefault:"db/migrations"`
}

func Database() *database {
	cfg := &database{}
	env.Parse(cfg)
	return cfg
}

func DBConnectionSource() string {
	db := Database()
	return db.Username + ":" + db.Password + "@tcp(" + db.Host + ":" + strconv.Itoa(db.Port) + ")/" + db.Name
}
