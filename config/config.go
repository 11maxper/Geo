package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Dadata struct {
	Key       string `env:"DADATA_KEY"`
	SecretKey string `env:"DADATA_SECRET_KEY"`
}

type JWT struct {
	JWTKey string `env:"JWT_KEY"`
}

type HTTP struct {
	Port string `env:"PORT" env-default:":8080"`
}

type DB struct {
	Password string `env:"DB_PASSWORD"`
	User     string `env:"DB_USER"`
	Name     string `env:"DB_NAME"`
	Port     string `env:"DB_PORT"`
	Host     string `env:"DB_HOST"`
}

type Cache struct {
	Addr     string `env:"REDIS_ADDR" env-default:"redis:6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
}

type Config struct {
	HTTP   HTTP
	Dadata Dadata
	JWTKey JWT
	DB     DB
	Cache  Cache
}

func Init(env ...string) (*Config, error) {
	err := godotenv.Load(env...)
	if err != nil {
		return nil, err
	}

	var cfg = &Config{
		HTTP: HTTP{
			Port: os.Getenv("PORT"),
		},
		Dadata: Dadata{
			Key:       os.Getenv("DADATA_KEY"),
			SecretKey: os.Getenv("DADATA_SECRET_KEY"),
		},
		JWTKey: JWT{
			JWTKey: os.Getenv("JWT_KEY"),
		},
		DB: DB{
			Password: os.Getenv("DB_PASSWORD"),
			User:     os.Getenv("DB_USER"),
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
			Host:     os.Getenv("DB_HOST"),
		},
		Cache: Cache{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
	}

	if cfg.HTTP.Port == "" {
		return nil, errors.New("HTTP config is empty")
	}

	if cfg.Dadata.Key == "" || cfg.Dadata.SecretKey == "" {
		return nil, errors.New("API config is empty")
	}

	if cfg.JWTKey.JWTKey == "" {
		return nil, errors.New("JWT config is empty")
	}
	if cfg.DB.Password == "" {
		return nil, errors.New("DB password is empty")
	}
	if cfg.DB.User == "" {
		return nil, errors.New("DB user is empty")
	}
	if cfg.DB.Name == "" {
		return nil, errors.New("DB name is empty")
	}
	if cfg.DB.Port == "" {
		return nil, errors.New("DB port is empty")
	}
	if cfg.DB.Host == "" {
		return nil, errors.New("DB host is empty")
	}

	if cfg.Cache.Addr == "" {
		return nil, errors.New("REDIS addr is empty")
	}
	if cfg.Cache.Password == "" {
		return nil, errors.New("REDIS password is empty")
	}

	return cfg, nil
}
