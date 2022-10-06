package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort    = 8080
	delayInSeconds = 8
)

func isAvailable(id string) bool {
	// for this example, only IDs 1 to 10 are available
	iid, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	if iid < 1 || iid > 10 {
		return false
	}

	return true
}

func handleAvailability(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scenario := r.URL.Query().Get("scenario")
	if scenario == "slow" {
		time.Sleep(time.Duration(rand.Intn(delayInSeconds)) * time.Second)
	} else if scenario == "error" {
		time.Sleep(time.Duration(rand.Intn(delayInSeconds)) * time.Second)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isAvailable(itemID) {
		fmt.Fprintf(w, "Item %s is available", itemID)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Item %s is not available", itemID)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	port, ok := os.LookupEnv("CONTOSO_PORT")
	if !ok {
		port = fmt.Sprintf("%d", defaultPort)
	}

	http.HandleFunc("/check", handleAvailability)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
