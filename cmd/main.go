package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/api-assignment/pkg/config"
	router "github.com/api-assignment/pkg/routes"
	"github.com/api-assignment/pkg/utils/db"
	"github.com/api-assignment/pkg/utils/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("⚠️ No .env file found, using system environment variables")
	}
	router := router.MainRouter()
	port := strconv.Itoa(config.GetConfig().GetAppPort())
	log := logger.InitializeAppLogger()
	_ = db.GetDBConn()
	log.Info("Starting API Server on Port ", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("unable to start server on port %d ", port, err)
	}
}
