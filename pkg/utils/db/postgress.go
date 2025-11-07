package db

import (
	"fmt"
	"os"

	"github.com/go-auth-microservice/pkg/utils/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgressDB struct {
	db         *gorm.DB
	connString string
	log        *zap.SugaredLogger
}

func (sdb *postgressDB) GetDB() *gorm.DB {
	return sdb.db
}
func (sdb *postgressDB) AutoMigrate(models ...interface{}) error {
	log := logger.InitializeAppLogger()
	err := sdb.db.AutoMigrate(models...)
	if err != nil {
		log.Panicf("Failed to migrate the database schema: %v", err)
	}
	return err
}
func (sdb *postgressDB) createConnectionString() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	timezone := os.Getenv("DB_TIMEZONE")

	if timezone == "" {
		timezone = "Asia/Kolkata"
	}
	if len(host) == 0 || len(port) == 0 || len(user) == 0 || len(password) == 0 || len(dbname) == 0 {
		sdb.log.Fatal("unable to read env variables for postgress connection")
	}
	// Build DSN
	sdb.connString = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		host, user, password, dbname, port, timezone,
	)
	sdb.log.Info("connection string for postgress has been created")
}
func (sdb *postgressDB) Initialize() {
	sdb.log = logger.InitializeAppLogger()
	sdb.createConnectionString()
	pgDB, err := gorm.Open(postgres.Open(sdb.connString), &gorm.Config{})
	sdb.db = pgDB
	if err != nil {
		sdb.log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	sdb.log.Info("connected to postgress database")
}
