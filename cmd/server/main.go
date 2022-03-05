package main

import (
	"embed"
	"log"

	"github.com/georgemblack/shoebox"
	"github.com/georgemblack/shoebox/pkg/config"
	"github.com/georgemblack/shoebox/pkg/firestore"
	"github.com/georgemblack/shoebox/pkg/handlers"
	"github.com/gin-gonic/gin"
)

//go:embed config/*
var configFiles embed.FS
var datastore firestore.Datastore

func main() {
	// Load config
	config, err := config.LoadConfig(configFiles)
	if err != nil {
		log.Fatal(err)
	}

	// Init Shoebox service
	err = shoebox.Init()
	if err != nil {
		log.Fatalf("failed to initialize application; %v", err)
	}

	// Init datastore
	datastore, err = firestore.GetDatastoreClient(config)
	if err != nil {
		log.Fatalf("failed to initialize datastore; %v", err)
	}

	router := gin.Default()

	router.Use(handlers.PreflightHandler(config))
	router.GET("/api/entries", handlers.GetEntriesHandler)
	router.POST("/api/entries", handlers.PostEntryHandler)
	router.DELETE("/api/entries/:entry_id", handlers.DeleteEntryHandler(datastore))

	router.Run(":" + config.APIPort)
}
