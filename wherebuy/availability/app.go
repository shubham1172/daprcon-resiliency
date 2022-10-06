package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort = 8000
)

var (
	contosoApiUrl = "http://localhost:8080"
	client        = &http.Client{Timeout: 5 * time.Second}
)

func handleAvailability(w http.ResponseWriter, r *http.Request) {
	itemID := r.FormValue("id")
	if itemID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scenario := r.FormValue("scenario")

	// Check availability with Contoso API
	res, err := client.Get(contosoApiUrl + "?id=" + itemID + "&scenario=" + scenario)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		w.WriteHeader(res.StatusCode)
		return
	}

	if res.StatusCode == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Item %s is not available", itemID)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Item %s is available", itemID)
}

func main() {
	apiUrl, ok := os.LookupEnv("CONTOSO_API_URL")
	if ok {
		contosoApiUrl = apiUrl
	}

	port, ok := os.LookupEnv("WHEREBUY_AVAILABILITY_PORT")
	if !ok {
		port = fmt.Sprintf("%d", defaultPort)
	}

	http.HandleFunc("/check", handleAvailability)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
