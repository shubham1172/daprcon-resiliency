package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort = 9000
)

var (
	client             = &http.Client{Timeout: 5 * time.Second}
	availabilityApiUrl = "http://localhost:8000"
	productApiUrl      = "http://localhost:8001"
)

func getAvailability(id string) (bool, error) {
	res, err := client.Get(availabilityApiUrl + "/check?id=" + id)
	if err != nil {
		log.Println(err)
		return false, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		return false, fmt.Errorf("availability API returned status code %d", res.StatusCode)
	}

	if res.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

func getProductInfo(id string) (string, error) {
	res, err := client.Get(productApiUrl + "/get?id=" + id)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		return "", fmt.Errorf("product API returned status code %d", res.StatusCode)
	}

	if res.StatusCode == http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		productInfo := buf.String()
		log.Printf("Product info: %s", productInfo)
		return productInfo, nil
	}

	return "", nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to WhereBuy! Use /query?id= to query for an item.")
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	itemId := r.URL.Query().Get("id")
	if itemId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	available, err := getAvailability(itemId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	productInfo, err := getProductInfo(itemId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !available {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Item %s is not available.", itemId)
		return
	}

	fmt.Fprintf(w, "Item %s is available. Details: %s", itemId, productInfo)
}

func main() {
	apiUrl, ok := os.LookupEnv("WHEREBUY_AVAILABILITY_API_URL")
	if ok {
		availabilityApiUrl = apiUrl
	}

	apiUrl, ok = os.LookupEnv("WHEREBUY_PRODUCT_API_URL")
	if ok {
		productApiUrl = apiUrl
	}

	port, ok := os.LookupEnv("WHEREBUY_FRONTEND_PORT")
	if !ok {
		port = fmt.Sprintf("%d", defaultPort)
	}

	http.HandleFunc("/query", handleInfo)
	http.HandleFunc("/", handleRoot)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
