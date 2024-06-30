package database

import (
	"context"
	"fmt"
	"time"

	"rakamin-final-task/config"
	"rakamin-final-task/helpers/log"
	"rakamin-final-task/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	ORM    *gorm.DB
	Config config.SQL
}

func Init(dbLogger log.LogInterface, config config.SQL) *DB {
	orm, err := initPostgres(dbLogger, config)
	if err != nil {
		dbLogger.Fatal(context.Background(), err.Error())
		panic(err)
	}

	return &DB{ORM: orm, Config: config}
}

func initPostgres(dbLogger log.LogInterface, config config.SQL) (*gorm.DB, error) {
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
	sqlDB.SetMaxIdleConns(int(config.PoolConfig.MaxIdle))
	sqlDB.SetMaxOpenConns(int(config.PoolConfig.MaxOpen))
	sqlDB.SetConnMaxIdleTime(time.Duration(config.PoolConfig.ConnIdleSec) * time.Second)
	sqlDB.SetConnMaxLifetime(time.Duration(config.PoolConfig.ConnMaxLifetimeSec) * time.Second)

	

	return db, nil
}

func (db *DB) Migrate() {
	db.ORM.AutoMigrate(&models.Users{})
	db.ORM.AutoMigrate(&models.UserToken{})
	db.ORM.AutoMigrate(&models.Photo{})
}
