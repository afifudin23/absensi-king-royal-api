package config

import (
	"fmt"
	"net/url"
	"strings"

	driverMysql "github.com/go-sql-driver/mysql"

	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenMySQL(databaseURL string) (*gorm.DB, error) {
	dsn, err := normalizeMySQLDSN(databaseURL)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return db, nil
}

func normalizeMySQLDSN(databaseURL string) (string, error) {
	databaseURL = strings.TrimSpace(databaseURL)
	if databaseURL == "" {
		return "", fmt.Errorf("DATABASE_URL is required")
	}

	// Raw DSN already in mysql driver style.
	if !strings.HasPrefix(databaseURL, "mysql://") {
		return databaseURL, nil
	}

	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		return "", fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	if parsedURL.User == nil {
		return "", fmt.Errorf("parse DATABASE_URL: missing user info")
	}

	username := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()
	host := parsedURL.Host
	dbName := strings.TrimPrefix(parsedURL.Path, "/")

	if host == "" || dbName == "" {
		return "", fmt.Errorf("parse DATABASE_URL: invalid host/database")
	}

	mysqlConfig := driverMysql.NewConfig()
	mysqlConfig.User = username
	mysqlConfig.Passwd = password
	mysqlConfig.Addr = host
	mysqlConfig.Net = "tcp"
	mysqlConfig.DBName = dbName
	mysqlConfig.AllowNativePasswords = true
	mysqlConfig.ParseTime = true

	return mysqlConfig.FormatDSN(), nil
}
