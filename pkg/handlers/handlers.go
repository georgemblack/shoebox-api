package handlers

import (
	"net/http"
	"time"

	"github.com/georgemblack/shoebox/pkg/config"
	"github.com/georgemblack/shoebox/pkg/firestore"
	"github.com/georgemblack/shoebox/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func PreflightHandler(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.AddCORSHeaders {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "*")
		}
		c.Next()
	}
}

func GetEntriesHandler(firestore firestore.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		entries, err := firestore.GetEntries()
		if err != nil {
			c.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{"entries": entries})
	}
}

func PostEntryHandler(firestore firestore.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var entry types.Entry

		err := c.BindJSON(&entry)
		if err != nil {
			c.JSON(http.StatusBadRequest, newErrorResponse(err.Error()))
			return
		}

		entry.ID = uuid.New().String()
		now := time.Now()
		entry.Created = now
		entry.Updated = now

		err = firestore.CreateEntry(entry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
			return
		}

		c.Status(http.StatusCreated)
	}
}

func PutEntryHandler(firestore firestore.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		entryID := c.Param("entry_id")

		var newEntry types.Entry
		err := c.BindJSON(&newEntry)
		if err != nil {
			c.JSON(http.StatusBadRequest, newErrorResponse(err.Error()))
			return
		}

		oldEntry, err := firestore.GetEntry(entryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
			return
		}

		entry := types.MergeEntries(oldEntry, newEntry)
		err = firestore.CreateEntry(entry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func DeleteEntryHandler(firestore firestore.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		entryID := c.Param("entry_id")
		err := firestore.DeleteEntry(entryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, newErrorResponse(err.Error()))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
