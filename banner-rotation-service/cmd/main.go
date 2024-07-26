package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yuriiwanchev/banner-rotation-service/internal/api"
)

func main() {
	port := ":8080"

	server := &http.Server{
		Addr:         ":8080",
		Handler:      http.DefaultServeMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	http.HandleFunc("/add-banner", api.AddBannerHandler)
	http.HandleFunc("/remove-banner", api.RemoveBannerHandler)
	http.HandleFunc("/record-click", api.RecordClickHandler)
	http.HandleFunc("/select-banner", api.SelectBannerHandler)

	fmt.Printf("Starting server on %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
