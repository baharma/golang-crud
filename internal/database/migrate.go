package database

import (
	"crud-test/internal/entity"
	"fmt"
)

func Migrate() error {
	serverDB, err := connectServer()
	if err != nil {
		return err
	}

	dbName, err := databaseName()
	if err != nil {
		return err
	}

	if err := serverDB.Exec(
		fmt.Sprintf(
			"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
			dbName,
		),
	).Error; err != nil {
		return fmt.Errorf("create database %s: %w", dbName, err)
	}

	appDB, err := Connect()
	if err != nil {
		return err
	}

	if err := appDB.AutoMigrate(&entity.Category{}, &entity.Book{}); err != nil {
		return fmt.Errorf("auto migrate schema: %w", err)
	}

	return nil
}
