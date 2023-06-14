package database

import (
	"cloud-app-hive/domain"
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ErrDatabaseConnection = errors.New("failed to connect to database")

var ErrDatabaseMigration = errors.New("failed to migrate database")

func ConnectToDatabase() (*gorm.DB, error) {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database,
	)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, ErrDatabaseConnection
	}

	return db, nil
}

func MigrateDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&domain.Application{}, &domain.Namespace{}, &domain.NamespaceMembership{})
	if err != nil {
		return ErrDatabaseMigration
	}

	return nil
}
