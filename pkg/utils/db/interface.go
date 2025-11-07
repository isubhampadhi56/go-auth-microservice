package db

import (
	"log"
	"os"

	"gorm.io/gorm"
)

type DB interface {
	GetDB() *gorm.DB
	AutoMigrate(...interface{}) error
	Initialize()
}

var dBConn DB

func GetDBConn() DB {
	if dBConn != nil {
		return dBConn
	}

	dbType := os.Getenv("DB_TYPE")
	switch dbType {
	case "sqlite":
		dBConn = &sqliteDB{}
		dBConn.Initialize()
	case "postgress":
		dBConn = &postgressDB{}
		dBConn.Initialize()
	default:
		log.Fatalf("invalid database type %s", dbType)
	}

	return dBConn
}
