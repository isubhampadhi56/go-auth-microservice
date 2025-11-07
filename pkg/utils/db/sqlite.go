package db

import (
	"os"

	"github.com/go-auth-microservice/pkg/utils/logger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqliteDB struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

func (sdb *sqliteDB) GetDB() *gorm.DB {
	return sdb.db
}
func (sdb *sqliteDB) AutoMigrate(models ...interface{}) error {
	log := logger.InitializeAppLogger()
	err := sdb.db.AutoMigrate(models...)
	if err != nil {
		log.Panicf("Failed to migrate the database schema: %v", err)
	}
	return err
}
func (sdb *sqliteDB) Initialize() {
	sdb.log = logger.InitializeAppLogger()
	if _, err := os.Stat("users.db"); err == nil {
		sdb.log.Info("Database already exists. Skipping initialization.")
		sdb.db, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
		if err != nil {
			sdb.log.Fatalf("Failed to connect to the database: %v", err)
		}
		sdb.log.Info("connected to sqlite database")
	} else if os.IsNotExist(err) {
		sdb.log.Info("Database does not exist. Initializing database.")
		sdb.db, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
		if err != nil {
			sdb.log.Fatalf("Failed to connect to the database: %v", err)
		}
		sdb.log.Info("connected to sqlite database")
	} else {
		sdb.log.Fatalf("Error checking the database file: %v", err)
	}
}
