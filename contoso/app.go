package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort = 8080
)

func handleAvailability(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Randomly wait between 0 to 20 seconds
	time.Sleep(time.Duration(rand.Intn(20)) * time.Second)

	// Check availability
	fmt.Fprintf(w, "Item %s is available", itemID)
}

func main() {
	port, ok := os.LookupEnv("CONTOSO_PORT")
	if !ok {
		port = fmt.Sprintf("%d", defaultPort)
	}

	http.HandleFunc("/check", handleAvailability)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
