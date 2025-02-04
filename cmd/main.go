/*
Copyright © 2025 Herman Gonçalves hermangoncalves@outlook.com
*/
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hermangoncalves/go-routeros-exporter/adapters/mikrotik"
	"github.com/hermangoncalves/go-routeros-exporter/adapters/prometheus"
	"github.com/hermangoncalves/go-routeros-exporter/config"
	"github.com/hermangoncalves/go-routeros-exporter/core/service"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func init() {
	config.Load()
	logger = config.NewLogger()
}

func main() {
	mikrotikClient, err := mikrotik.NewMikrotikClient(
		context.Background(),
		fmt.Sprintf("%s:%d", config.MikrotikDevice.Host, config.MikrotikDevice.Port),
		config.MikrotikDevice.Username,
		config.MikrotikDevice.Password,
	)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"host":     config.MikrotikDevice.Host,
			"port":     config.MikrotikDevice.Port,
			"username": config.MikrotikDevice.Username,
		}).Fatalf("Failed to connect to Mikrotik device: %v", err)
	}

	logger.Info("Successfully connected to Mikrotik device")

	metricsService := service.NewNetricsService(mikrotikClient)
	metricsHandler := prometheus.NewMetricsHandler(metricsService, logger)

	// Start a goroutine to update metrics periodically
	go func() {
		for {
			if err := metricsHandler.UpdateMetrics(); err != nil {
				logger.WithError(err).Error("Failed to update metrics")
			} else {
				logger.Info("Metrics updated successfully")
			}
			time.Sleep(15 * time.Second)
		}
	}()

	// Expose metrics endpoint
	port := ":8080"
	logger.Infof("Starting server on port %s", port)
	http.Handle("/metrics", metricsHandler)
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
