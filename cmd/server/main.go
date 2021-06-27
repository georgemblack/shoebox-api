package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/georgemblack/shoebox"
)

func entryHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Println(err)
			http.Error(response, "Bad request", http.StatusBadRequest)
			return
		}
		var parsedBody map[string]interface{}
		err = json.Unmarshal(body, &parsedBody)
		if err != nil {
			log.Println(err)
			http.Error(response, "Bad request", http.StatusBadRequest)
			return
		}
		entry, err := shoebox.ParseEntry(parsedBody)
		if err != nil {
			log.Println(err)
			http.Error(response, "Bad request", http.StatusBadRequest)
			return
		}
		err = shoebox.CreateEntry(entry)
		if err != nil {
			log.Println(err)
			http.Error(response, "Internal error", http.StatusInternalServerError)
			return
		}
		responseBody, err := json.Marshal(entry)
		if err != nil {
			log.Println(err)
			http.Error(response, "Internal error", http.StatusInternalServerError)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		_, err = response.Write(responseBody)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	err := shoebox.Init()
	if err != nil {
		log.Fatalf("failed to initialize application; %v", err)
	}
	port := getEnv("PORT", "8080")
	http.HandleFunc("/entries", entryHandler)
	log.Println("Listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
