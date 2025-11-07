package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second, // optional but recommended
	}
	log.Info("Starting API Server on Port ", port)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("unable to start server on port %d ", port, err)
	}
}
