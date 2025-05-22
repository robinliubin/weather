package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/robinliubin/weather/weather"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/alerts", weather.HandleAlerts)
	http.HandleFunc("/forecast", weather.HandleForecast)
	http.HandleFunc("/forecast/city", weather.HandleForecastByCity)

	fmt.Printf("Starting weather server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}