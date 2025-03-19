package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"case-service/api"
	"case-service/util"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"

	log "github.com/sirupsen/logrus"
)

func main() {
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		log.Info("Log level set to DEBUG level")
		log.SetLevel(log.DebugLevel)
	}

	if strings.ToLower(os.Getenv("PROFILE")) != "dev" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	log.Info("===========================================")
	log.Infof("Starting %s apps on port %s", os.Getenv("APP_NAME"), os.Getenv("PORT"))

	// init db
	db := util.GetDBConnection()
	sqlDB, err := db.DB()
	if err != nil {
		sqlDB.Close()
	}

	// config API router
	log.Info("Configuring API router...")
	router := api.ConfigRouter()

	// Start server
	serverPort := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Infof("%s listen on port %s", os.Getenv("APP_NAME"), serverPort)
	router.Server.Addr = serverPort

	go func() {
		if err := router.Start(serverPort); err != nil {
			router.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := router.Shutdown(ctx); err != nil {
		router.Logger.Fatal(err)
	}
}
