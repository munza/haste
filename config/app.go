package config

type App struct {
	Name      string `default:"Haste"`
	URL       string `default:"127.0.0.1"`
	SecretKey string `default:"somesupersecretkey"`
}
