package config

type Database struct {
	Host     string `default:"127.0.0.1"`
	Port     int    `default:"3306"`
	Name     string `default:"haste_db"`
	Username string `default:"root"`
	Password string `default:"root"`
}
