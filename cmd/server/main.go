package main

import (
	"log"
	"os"

	"github.com/georgemblack/shoebox"
	"github.com/georgemblack/shoebox/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	err := shoebox.Init()
	if err != nil {
		log.Fatalf("failed to initialize application; %v", err)
	}

	port := getEnv("PORT", "8080")
	router := gin.Default()

	router.GET("/api/entries", handlers.GetEntriesHandler)
	router.POST("/api/entries", handlers.PostEntryHandler)
	router.Run(":" + port)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
