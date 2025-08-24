package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")

	if name == "" {
		name = "World"
	}

	response := map[string]string{
		"message": "Hello " + name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	fmt.Printf("Starting server on :8000\n")

	http.HandleFunc("/hello", helloHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
