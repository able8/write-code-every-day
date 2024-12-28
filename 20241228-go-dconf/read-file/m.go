package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Config represents our configuration
type Config struct {
	IntervalSecond int
}

// ConfigReader handles reading and watching the config file
type ConfigReader struct {
	config  Config
	watcher *fsnotify.Watcher
	stopCh  chan struct{}
	mu      sync.Mutex
}

// NewConfigReader creates a new ConfigReader
func NewConfigReader(filename string) (*ConfigReader, error) {
	reader := &ConfigReader{
		stopCh: make(chan struct{}),
		config: Config{
			IntervalSecond: 5,
		},
	}

	go reader.readPeriodically()
	return reader, nil
}

// readPeriodically reads the config file every interval
func (cr *ConfigReader) readPeriodically() {
	ticker := time.NewTicker(time.Duration(cr.config.IntervalSecond) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := cr.updateConfig()
			if err != nil {
				fmt.Printf("Error updating config: %v\n", err)
			}
		case <-cr.stopCh:
			ticker.Stop()
			return
		}
	}
}

// updateConfig reads the latest config from the file
func (cr *ConfigReader) updateConfig() error {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}

	var conf Config
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return err
	}

	cr.mu.Lock()
	defer cr.mu.Unlock()

	if conf.IntervalSecond != cr.config.IntervalSecond {
		fmt.Printf("Config interval changed: %v\n", conf.IntervalSecond)
	}

	cr.config = conf
	return nil
}

// GetConfig returns the current config
func (cr *ConfigReader) GetConfig() Config {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	return cr.config
}

// Stop stops the config reader goroutine
func (cr *ConfigReader) Stop() {
	close(cr.stopCh)
}

func main() {
	configReader, err := NewConfigReader("config.json")
	if err != nil {
		fmt.Printf("Failed to create config reader: %v\n", err)
		return
	}
	defer configReader.Stop()

	for {
		config := configReader.GetConfig()
		fmt.Printf("Current config interval: %v\n", config.IntervalSecond)
		time.Sleep(time.Second * 5)
	}
}
