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
	api := router.Group("/api")

	api.Use(handlers.PreflightHandler(config))
	api.OPTIONS("/*path", handlers.OptionsHanlder)
	api.GET("/entries", handlers.GetEntriesHandler(datastore))
	api.POST("/entries", handlers.PostEntryHandler(datastore))
	api.GET("/entries/:entry_id", handlers.GetEntryHandler(datastore))
	api.PUT("/entries/:entry_id", handlers.PutEntryHandler(datastore))
	api.DELETE("/entries/:entry_id", handlers.DeleteEntryHandler(datastore))

	router.Run(":" + config.APIPort)
}
