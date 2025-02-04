package config

import "github.com/spf13/viper"

type Config struct {
	MikrotikDevice MikrotikDevice
}

type MikrotikDevice struct {
	Name     string
	Address  string
	Username string
	Password string
	Port     int
}

func LoadConfig() *Config {
	return &Config{
		MikrotikDevice: MikrotikDevice{
			Name: viper.GetString("device_name"),
			Address: viper.GetString("address"),
			Username: viper.GetString("username"),
			Password: viper.GetString("password"),
			Port: viper.GetInt("port"),
		},
	}
}
