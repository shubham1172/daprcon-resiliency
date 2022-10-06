package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Currency string `json:"currency"`
	Rating   string `json:"rating"`
}

const (
	defaultPort = 8001
)

var (
	data = map[string]Product{}
)

func handleProduct(w http.ResponseWriter, r *http.Request) {
	itemId := r.URL.Query().Get("id")
	if itemId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, ok := data[itemId]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func loadData() {
	// Well, in the real world we would read this from a database.
	file, err := os.Open("data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var products []Product
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&products); err != nil {
		log.Fatal(err)
	}

	for _, product := range products {
		data[product.ID] = product
	}
}

func main() {
	loadData()

	port, ok := os.LookupEnv("WHEREBUY_PRODUCT_PORT")
	if !ok {
		port = fmt.Sprintf("%d", defaultPort)
	}

	http.HandleFunc("/get", handleProduct)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
