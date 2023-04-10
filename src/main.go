package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	usernames, err := GetUsernames()
	if err != nil {
		http.Error(w, "Error fetching usernames from Redis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(usernames)
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = Initialize(config.LogFile)
	if err != nil {
		log.Fatalf("Error initializing the Minecraft package: %v", err)
	}

	go WatchLogFile()

	// Log usernames from Redis
	logUsernamesFromRedis()

	http.HandleFunc("/api/users", getUsersHandler)

	fmt.Printf("Server listening on port %d...\n", config.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.ServerPort), nil))

}
