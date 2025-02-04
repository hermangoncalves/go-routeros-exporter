/*
Copyright © 2025 Herman Gonçalves hermangoncalves@outlook.com
*/
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hermangoncalves/go-routeros-exporter/adapters/mikrotik"
	"github.com/hermangoncalves/go-routeros-exporter/adapters/prometheus"
	"github.com/hermangoncalves/go-routeros-exporter/config"
	"github.com/hermangoncalves/go-routeros-exporter/core/service"
)

func init() {
	config.Load()
}

func main() {

	mikrotikClient, err := mikrotik.NewMikrotikClient(
		context.Background(),
		fmt.Sprintf("%s:%d", config.MikrotikDevice.Host, config.MikrotikDevice.Port),
		config.MikrotikDevice.Username,
		config.MikrotikDevice.Password,
	)

	if err != nil {
		panic(err)
	}

	metricsService := service.NewNetricsService(mikrotikClient)
	metricsHandler := prometheus.NewMetricsHandler(metricsService)

	// Start a goroutine to update metrics periodically
	go func() {
		for {
			if err := metricsHandler.UpdateMetrics(); err != nil {
				log.Printf("Failed to update metrics: %v", err)
			}
			time.Sleep(15 * time.Second)
		}
	}()

	// Expose metrics endpoint
	port := ":8080"
	http.Handle("/metrics", metricsHandler)
	log.Println("Starting server on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
