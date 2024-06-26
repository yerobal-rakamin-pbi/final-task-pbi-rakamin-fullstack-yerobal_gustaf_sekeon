package main

import (
	"os"

	"rakamin-final-task/config"
	"rakamin-final-task/database"
	"rakamin-final-task/helpers/configbuilder"
	"rakamin-final-task/helpers/configreader"
	"rakamin-final-task/helpers/files"
	"rakamin-final-task/helpers/log"
)

const (
	configFile = "./config/config.json"
)

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

	// Init DB Connection
	dbConfig := database.Config{
		Host:     config.SQL.Host,
		Port:     config.SQL.Port,
		Username: config.SQL.Username,
		Password: config.SQL.Password,
		Database: config.SQL.Database,
	}

	db := database.Init(logger, dbConfig)
}
