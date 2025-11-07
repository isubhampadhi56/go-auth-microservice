package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	router "github.com/go-auth-microservice/pkg/routes"
	"github.com/go-auth-microservice/pkg/utils/db"
	"github.com/go-auth-microservice/pkg/utils/logger"
	"github.com/joho/godotenv"
)

func main() {
	if strings.ToUpper(os.Getenv("APP_ENV")) != "PRODUCTION" {
		err := godotenv.Load()
		if err != nil {
			log.Panicf("⚠️ No .env file found, using system environment variables")
		}
	}
	router := router.MainRouter()
	port := os.Getenv("API_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log := logger.InitializeAppLogger()
	_ = db.GetDBConn()
	log.Info("Starting API Server on Port ", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("unable to start server on port %d ", port, err)
	}
}
