package config

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Dialect  string
	Host     string
	Username string
	Password string
	Name     string
	Port     string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  "postgres",
			Host:     "localhost",
			Username: "postgres",
			Password: "123456",
			Name:     "demo_db",
			Port:     "5432",
		},
	}
}
