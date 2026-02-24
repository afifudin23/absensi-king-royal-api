package config

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

var (
	once     sync.Once
	initErr  error
	appEnv   *EnvConfig
	database *gorm.DB
)

func Init() error {
	once.Do(func() {
		appEnv, initErr = LoadEnv()
		if initErr != nil {
			initErr = fmt.Errorf("load env: %w", initErr)
			return
		}

		database, initErr = OpenMySQL(appEnv.DatabaseURL)
		if initErr != nil {
			initErr = fmt.Errorf("connect database: %w", initErr)
			return
		}
	})

	return initErr
}

func GetEnv() *EnvConfig {
	return appEnv
}

func GetDB() *gorm.DB {
	return database
}

func CloseDB() error {
	if database == nil {
		return nil
	}
	sqlDB, err := database.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
