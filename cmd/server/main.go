package main

import (
	"fmt"
	"net/http"

	"github.com/AksanovK/url-monitor/internal/api"
)

func main() {
	router := api.NewRouter()

	port := ":8080"
	fmt.Printf("Server starting on %s\n", port)

	if err := http.ListenAndServe(port, router); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
