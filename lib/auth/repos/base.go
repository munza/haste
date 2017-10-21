package authrepo

import (
	"haste/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/joho/godotenv/autoload"
)

type BaseRepo struct{}

func DB() *gorm.DB {
	conn, err := gorm.Open(config.Database().Driver, config.DBConnectionSource())
	if err != nil {
		panic(err)
	}
	return conn
}
