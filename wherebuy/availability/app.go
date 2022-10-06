package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// for this example, only IDs 1 to 10 are available
	itemIdInt, err := strconv.Atoi(itemID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Item %s is not available", r.URL.Query().Get("id"))

		return
	}

	if itemIdInt < 1 || itemIdInt > 10 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Item %s is not available", r.URL.Query().Get("id"))
		return
	}

	// Check availability with Contoso API
	res, err := client.Get(contosoApiUrl + "/check?id=" + itemID)
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
		fmt.Fprintf(w, "Item %s is not available", r.URL.Query().Get("id"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Item %s is available", r.URL.Query().Get("id"))
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
