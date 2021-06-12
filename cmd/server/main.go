package main

import (
	"log"
	"net/http"
	"os"
)

// Build is the return payload
type Build struct {
	BuildID string `json:"buildID"`
}

func main() {
	port := getEnv("PORT", "9002")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bogus!"))
	})

	log.Println("Listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
