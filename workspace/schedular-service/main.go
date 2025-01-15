package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"schedular-service/api"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func main() {
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		log.Info("Log level set to DEBUG level")
		log.SetLevel(log.DebugLevel)
	}

	log.Info("===========================================")
	log.Infof("Starting %s apps on port %s", os.Getenv("APP_NAME"), os.Getenv("PORT"))

	// get access token for first time
	success := api.GetCase()
	if !success {
		log.Error("Unable to get GetCase!! Program exit")
		os.Exit(1)
	}

	c := cron.New()

	log.Infof("Starting CRON job for get GetCase : %s", os.Getenv("ACCESS_CRON"))
	c.AddFunc(os.Getenv("ACCESS_CRON"), func() {
		api.GetCase()
	})

	c.Start()

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
