package config

import "os"

var config Config

type Config struct {
	Database
}

type Database struct {
	User     string
	Password string
	Host     string
	Name     string
	Port     string
}

func Instance() Config {
	return config
}

func Init(c ...Config) Config {
	if len(c) == 1 {
		config = c[0]
		return config
	}
	config = Config{
		Database: Database{
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Host:     os.Getenv("DATABASE_HOST"),
			Name:     os.Getenv("DATABASE_NAME"),
			Port:     os.Getenv("DATABASE_PORT"),
		},
	}
	return config
}
