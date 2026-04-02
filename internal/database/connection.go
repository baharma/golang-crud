package database

import (
	"crud-test/internal/config"
	"fmt"
	"regexp"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var databaseNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func Connect() (*gorm.DB, error) {
	godotenv.Load()

	dsn, err := applicationDSN()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	return db, nil
}

func connectServer() (*gorm.DB, error) {
	godotenv.Load()

	db, err := gorm.Open(mysql.Open(serverDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect mysql server: %w", err)
	}

	return db, nil
}

func applicationDSN() (string, error) {
	dbName, err := databaseName()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetEnv("DB_USER", "root"),
		config.GetEnv("DB_PASS", "baharma1899"),
		config.GetEnv("DB_HOST", "127.0.0.1"),
		config.GetEnv("DB_PORT", "3306"),
		dbName,
	), nil
}

func serverDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetEnv("DB_USER", "root"),
		config.GetEnv("DB_PASS", "baharma1899"),
		config.GetEnv("DB_HOST", "127.0.0.1"),
		config.GetEnv("DB_PORT", "3306"),
	)
}

func databaseName() (string, error) {
	name := config.GetEnv("DB_NAME", "crud_go")
	if !databaseNamePattern.MatchString(name) {
		return "", fmt.Errorf("invalid DB_NAME %q", name)
	}

	return name, nil
}
