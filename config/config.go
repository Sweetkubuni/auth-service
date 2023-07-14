package config

import "github.com/caarlos0/env"

type IConfig interface {
	GetAppConfig() Config
}

type Config struct {
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string
	JwtSecret  []byte
}

type DevelopmentConfig struct {
	DbHost     string `env:"DB_HOST" envDefault:"localhost"`
	DbPort     string `env:"DB_PORT" envDefault:"3600"`
	DbName     string `env:"DB_NAME" envDefault:"jwt_demo"`
	DbUser     string `env:"DB_USER" envDefault:"root"`
	DbPassword string `env:"DB_PASSWORD" envDefault:"password"`
	JwtSecret  string `env:"JWT_SECRET" envDefault:"invisiblekey!"`
}

type ProductionConfig struct {
	DbHost     string `env:"DB_HOST" envDefault:"localhost"`
	DbPort     string `env:"DB_PORT" envDefault:"3600"`
	DbName     string `env:"DB_NAME" envDefault:"jwt_demo"`
	DbUser     string `env:"DB_USER" envDefault:"root"`
	DbPassword string `env:"DB_PASSWORD" envDefault:"password"`
	JwtSecret  string `env:"JWT_SECRET" envDefault:"hsaldfhjlaslvgosdhf!"`
}

// DEVELOPMENT CONFIG
func NewDevelopmentConfig() IConfig {
	return &DevelopmentConfig{}
}

func (dev DevelopmentConfig) GetAppConfig() Config {
	return Config{
		DbHost:     dev.DbHost,
		DbPort:     dev.DbPort,
		DbName:     dev.DbName,
		DbUser:     dev.DbUser,
		DbPassword: dev.DbPassword,
		JwtSecret:  []byte(dev.JwtSecret),
	}
}

// PROD CONFIG
func NewProductionConfig() IConfig {
	return &ProductionConfig{}
}

func (dev ProductionConfig) GetAppConfig() Config {
	return Config{
		DbHost:     dev.DbHost,
		DbPort:     dev.DbPort,
		DbName:     dev.DbName,
		DbUser:     dev.DbUser,
		DbPassword: dev.DbPassword,
		JwtSecret:  []byte(dev.JwtSecret),
	}
}

func GetConfig(e string) (IConfig, error) {
	var cfg IConfig = nil

	if e == "development" {
		cfg = NewDevelopmentConfig()
	}

	if e == "production" {
		cfg = NewProductionConfig()
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
