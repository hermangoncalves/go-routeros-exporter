package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Devices         []DeviceConfig `yaml:"devices"`          // List of MikroTik devices
	CollectInterval string         `yaml:"collect_interval"` // How often to collect metrics (e.g., "30s")
	ListenAddress   string         `yaml:"listen_address"`   // Address to expose Prometheus metrics (e.g., ":9283")
}

type DeviceConfig struct {
	Name     string `yaml:"name"`     // Friendly name for the device
	Host     string `yaml:"host"`     // IP or hostname of the MikroTik device
	Port     int    `yaml:"port"`     // API port (default: 8728)
	Username string `yaml:"username"` // API username
	Password string `yaml:"password"` // API password
}

func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
