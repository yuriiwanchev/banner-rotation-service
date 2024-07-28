package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yuriiwanchev/banner-rotation-service/internal/api"
)

func main() {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	api.InitKafkaProducer([]string{kafkaBrokers}, kafkaTopic)

	port := ":8080"

	server := &http.Server{
		Addr:         port,
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
