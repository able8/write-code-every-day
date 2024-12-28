package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	// Initialize Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Watch for config changes
	viper.WatchConfig()

	// Set up callback function for config changes
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	// Example configuration values
	port := viper.GetInt("port")
	host := viper.GetString("host")

	fmt.Printf("Server running on port %d at host %s\n", port, host)

	// Main application logic
	for {
		// Your main application logic here
	}
}
