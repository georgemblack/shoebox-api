package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/georgemblack/shoebox"
	"github.com/gin-gonic/gin"
)

func getOptions(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func getEntries(c *gin.Context) {
	entries, err := shoebox.GetEntries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get entries"})
		return
	}

	
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.JSON(http.StatusOK, gin.H{"entries": entries.Entries})
}

func postEntry(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to read request body",
		})
		return
	}
	var parsedBody map[string]any
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse request body",
		})
		return
	}
	entry, err := shoebox.ParseEntry(parsedBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse request body",
		})
		return
	}
	err = shoebox.CreateEntry(entry)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create entry",
		})
		return
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, gin.H{
		"message": "entry created",
	})
}

func main() {
	err := shoebox.Init()
	if err != nil {
		log.Fatalf("failed to initialize application; %v", err)
	}
	port := getEnv("PORT", "8080")
	router := gin.Default()
	router.GET("/entries", getEntries)
	router.POST("/entries", postEntry)
	router.OPTIONS("/entries" , getOptions)
	router.Run(":" + port)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
