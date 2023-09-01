package config

import (
	"strings"

	"github.com/caarlos0/env"
)

type IConfig interface {
	GetAppConfig() Config
}

type Config struct {
	DbHost         string
	DbPort         string
	DbName         string
	DbUser         string
	DbPassword     string
	JwtSecret      []byte
	TrustedOrigins []string
}

type DevelopmentConfig struct {
	DbHost     string `env:"DB_HOST" envDefault:"mysql-service"`
	DbPort     string `env:"DB_PORT" envDefault:"3306"`
	DbName     string `env:"DB_NAME" envDefault:"jwt_demo"`
	DbUser     string `env:"DB_USER" envDefault:"root"`
	DbPassword string `env:"DB_PASSWORD" envDefault:"admin"`
	JwtSecret  string `env:"JWT_SECRET" envDefault:"invisiblekey!"`
	Cors       string `env:"CORS" envDefault:"*"`
}

type ProductionConfig struct {
	DbHost     string `env:"DB_HOST" envDefault:"mysql-service"`
	DbPort     string `env:"DB_PORT" envDefault:"3306"`
	DbName     string `env:"DB_NAME" envDefault:"jwt_demo"`
	DbUser     string `env:"DB_USER" envDefault:"root"`
	DbPassword string `env:"DB_PASSWORD" envDefault:"password"`
	JwtSecret  string `env:"JWT_SECRET" envDefault:"hsaldfhjlaslvgosdhf!"`
	Cors       string `env:"CORS" envDefault:"localhost:8080"`
}

// DEVELOPMENT CONFIG
func NewDevelopmentConfig() IConfig {
	return &DevelopmentConfig{}
}

func (dev DevelopmentConfig) GetAppConfig() Config {
	origins := strings.Split(dev.Cors, ",")

	return Config{
		DbHost:         dev.DbHost,
		DbPort:         dev.DbPort,
		DbName:         dev.DbName,
		DbUser:         dev.DbUser,
		DbPassword:     dev.DbPassword,
		JwtSecret:      []byte(dev.JwtSecret),
		TrustedOrigins: origins,
	}
}

// PROD CONFIG
func NewProductionConfig() IConfig {
	return &ProductionConfig{}
}

func (dev ProductionConfig) GetAppConfig() Config {
	origins := strings.Split(dev.Cors, ",")
	return Config{
		DbHost:         dev.DbHost,
		DbPort:         dev.DbPort,
		DbName:         dev.DbName,
		DbUser:         dev.DbUser,
		DbPassword:     dev.DbPassword,
		JwtSecret:      []byte(dev.JwtSecret),
		TrustedOrigins: origins,
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
