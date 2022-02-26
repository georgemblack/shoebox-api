package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/georgemblack/shoebox"
	"github.com/georgemblack/shoebox/pkg/config"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func newErrorResponse(message string) errorResponse {
	return errorResponse{
		Message:   message,
		Timestamp: time.Now().Format(time.RFC850),
	}
}

func PreflightHandler(config config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		if config.AddCORSHeaders {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		c.Next()
	}
}

func GetEntriesHandler(c *gin.Context) {
	entries, err := shoebox.GetEntries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("failed to get entries"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries.Entries})
}

func PostEntryHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, newErrorResponse("failed to read request body"))
		return
	}
	var parsedBody map[string]interface{}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, newErrorResponse("failed to parse request body"))
		return
	}
	entry, err := shoebox.ParseEntry(parsedBody)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, newErrorResponse("failed to parse request body"))
		return
	}
	err = shoebox.CreateEntry(entry)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, newErrorResponse("failed to create entry"))
		return
	}

	c.Header("Content-Type", "application/json")
	c.Status(http.StatusCreated)
}
