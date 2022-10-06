package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	defaultPort   = 8000
	contosoApiUrl = "http://localhost:8080"
	client        = &http.Client{Timeout: 5 * time.Second}
)

func handleAvailability(w http.ResponseWriter, r *http.Request) {
	res, err := client.Get(contosoApiUrl + "/check?id=" + r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.StatusCode != http.StatusOK {
		w.WriteHeader(res.StatusCode)
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
