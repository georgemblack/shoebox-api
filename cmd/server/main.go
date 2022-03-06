package main

import (
	"embed"
	"log"

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
		log.Fatalf("failed to load config; %v", err)
	}

	// Init datastore
	datastore, err = firestore.GetDatastoreClient(config)
	if err != nil {
		log.Fatalf("failed to initialize datastore; %v", err)
	}

	router := gin.Default()

	router.Use(handlers.PreflightHandler(config))
	router.GET("/api/entries", handlers.GetEntriesHandler(datastore))
	router.POST("/api/entries", handlers.PostEntryHandler(datastore))
	router.PUT("/api/entries/:entry_id", handlers.PutEntryHandler(datastore))
	router.DELETE("/api/entries/:entry_id", handlers.DeleteEntryHandler(datastore))

	router.Run(":" + config.APIPort)
}
