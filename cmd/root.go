/*
Copyright © 2025 Herman hermangoncalves@outlook.com
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hermangoncalves/go-routeros-exporter/internal/adapters"
	"github.com/hermangoncalves/go-routeros-exporter/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-routeros-exporter",
	Short: "A Prometheus exporter for MikroTik RouterOS metrics.",
	Long: `go-routeros-exporter is an open-source tool written in Go that connects to MikroTik RouterOS devices via the API. 

It collects network metrics such as bandwidth usage, device status, and connected clients, exposing them in a format 
compatible with Prometheus for monitoring and alerting.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		address := fmt.Sprintf("%s:%d", cfg.MikrotikDevice.Address, cfg.MikrotikDevice.Port)

		authAdapter := adapters.NewMikrotikAuthAdapter(5 * time.Second)

		client, err := authAdapter.Authenticate(
			context.Background(),
			address,
			cfg.MikrotikDevice.Username,
			cfg.MikrotikDevice.Password,
		)

		if err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			return
		}
		defer client.Close()

		reply, err := client.RunCommand("/system/identity/print")
		if err != nil {
			log.Fatalf("Failed to execute command: %v", err)
		}

		for _, re := range reply.Re {
			fmt.Println("MikroTik Identity:", re.Map["name"])
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-routeros-exporter.yaml)")

	rootCmd.Flags().StringP("address", "r", "", "Mikrotik RouterOs address")
	rootCmd.Flags().StringP("username", "u", "", "Mikrotik Mikrotik API username")
	rootCmd.Flags().StringP("password", "p", "", "Mikrotik API password")
	rootCmd.Flags().IntP("port", "P", 8728, "Mikrotik RouterOs address")

	viper.BindPFlag("address", rootCmd.Flags().Lookup("address"))
	viper.BindPFlag("username", rootCmd.Flags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.Flags().Lookup("password"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".go-routeros-exporter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
