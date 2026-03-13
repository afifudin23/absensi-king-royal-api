package config

import (
	"fmt"
	"sync"

	"github.com/afifudin23/absensi-king-royal-api/internal/database"
	"gorm.io/gorm"
)

var (
	once    sync.Once
	initErr error
	appEnv  *EnvConfig
	db      *gorm.DB
)

func Init() error {
	once.Do(func() {
		appEnv, initErr = LoadEnv()
		if initErr != nil {
			initErr = fmt.Errorf("load env: %w", initErr)
			return
		}

		db, initErr = database.OpenMySQL(appEnv.DatabaseURL)
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
	return db
}

func CloseDB() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
