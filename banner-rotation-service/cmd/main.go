package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yuriiwanchev/banner-rotation-service/internal/api"
)

func main() {
	http.HandleFunc("/add-banner", api.AddBannerHandler)
	http.HandleFunc("/remove-banner", api.RemoveBannerHandler)
	http.HandleFunc("/record-click", api.RecordClickHandler)
	http.HandleFunc("/select-banner", api.SelectBannerHandler)

	port := ":8080"
	fmt.Printf("Starting server on %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
