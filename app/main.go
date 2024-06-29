package main

import (
	"os"

	"rakamin-final-task/config"
	repo "rakamin-final-task/controllers/repository"
	uc "rakamin-final-task/controllers/usecase"
	"rakamin-final-task/database"
	"rakamin-final-task/helpers/configbuilder"
	"rakamin-final-task/helpers/configreader"
	"rakamin-final-task/helpers/files"
	"rakamin-final-task/helpers/jwt"
	"rakamin-final-task/helpers/log"
	"rakamin-final-task/router"
)

const (
	configFile = "./config/config.json"
)

// @title Rakamin Backend
// @description API for Rakamin Backend
// @version 1.0

// @contact.name 	Yerobal Gustaf Sekeon
// @contact.email 	yerobalg@gmail.com

// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @value Bearer {token}

func main() {
	if !files.IsExist(configFile) {
		// Build config
		configbuilder.Init(configbuilder.GCPConfig{
			ProjectID:   os.Getenv("GCP_PROJECT_ID"),
			SecretID:    os.Getenv("GCP_SECRET_ID"),
			PrivateKey:  os.Getenv("GCP_PRIVATE_KEY"),
			ClientEmail: os.Getenv("GCP_CLIENT_EMAIL"),
			SecretName:  os.Getenv("GCP_SECRET_NAME"),
		}, configFile).BuildConfig()
	}

	// Init Config Reader
	configReader := configreader.Init(configFile)

	var config config.Application
	configReader.ReadConfig(&config)

	// Init Logger
	logger := log.Init()

	// Init JWT
	jwtLib := jwt.Init(config.Server.JWT.ExpSec, config.Server.JWT.Secret)

	// Init DB Connection
	dbConfig := database.Config{
		Host:     config.SQL.Host,
		Port:     config.SQL.Port,
		Username: config.SQL.Username,
		Password: config.SQL.Password,
		Database: config.SQL.Database,
	}

	db := database.Init(logger, dbConfig)

	// Init repository
	repository := repo.Init(db)

	// Init Usecase
	ucParam := uc.InitParam{
		Repo:       repository,
		ServerConf: config.Server,
		JwtLib:     jwtLib,
	}
	usecase := uc.Init(ucParam)

	// Init Router
	routerParam := router.InitParam{
		Config:  config,
		Log:     logger,
		DB:      db,
		Usecase: usecase,
	}
	router := router.Init(routerParam)

	router.Run()
}
