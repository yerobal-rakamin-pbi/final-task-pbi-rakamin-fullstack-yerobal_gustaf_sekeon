package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rakamin-final-task/helpers/log"
)

type DB struct {
	ORM    *gorm.DB
	Config Config
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func Init(dbLogger log.LogInterface, config Config) (*DB) {
	orm, err := initPostgres(dbLogger, config)
	if err != nil {
		dbLogger.Fatal(nil, err.Error())
	}

	return &DB{ORM: orm, Config: config}
}

func initPostgres(dbLogger log.LogInterface, config Config) (*gorm.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
	)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: log.New(dbLogger),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Pool configuration
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	return db, nil
}
