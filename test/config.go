package test

import "github.com/Jimeux/go-generic-dao/config"

func InitConfig() {
	config.Init(config.Config{
		Database: config.Database{
			User:     "root",
			Password: "dev",
			Host:     "localhost",
			Name:     "genericdaotest",
			Port:     "33306",
		},
	})
}
