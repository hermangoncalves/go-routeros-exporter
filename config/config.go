package config

import (
	"log"
	"os"
	"strconv"

	"github.com/hermangoncalves/go-routeros-exporter/core/domain"
	"github.com/joho/godotenv"
)

var MikrotikDevice *domain.MikrotikDevice

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port, _ := strconv.Atoi(os.Getenv("MIKROTIK_PORT"))
	MikrotikDevice = &domain.MikrotikDevice{
		Name:     os.Getenv("MIKROTIK_DEVICE_NAME"),
		Host:     os.Getenv("MIKROTIK_HOST"),
		Port:     port,
		Username: os.Getenv("MIKROTIK_USER"),
		Password: os.Getenv("MIKROTIK_PASSWORD"),
	}
}
