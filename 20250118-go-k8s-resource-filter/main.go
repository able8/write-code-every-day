package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

const DefaultConfigFileName = "config.yaml"
const DefaultConfigWatchInterval = 2 * time.Second

func main() {
	c := NewConfig(
		WithFileName(DefaultConfigFileName),
		WithWatchInterval(DefaultConfigWatchInterval),
	)

	err := c.ReadConfig()
	if err != nil {
		log.Println(err)
		return
	}

	go c.StartConfigWatcher()

	resources := []ResourceInfo{
		{ClusterName: "cluster1", Name: "pod1", Namespace: "ns1"},
		{ClusterName: "cluster2", Name: "pod2", Namespace: "ns2"},
		{ClusterName: "cluster3", Name: "pod3", Namespace: "ns3"},
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		for _, resource := range resources {
			if c.IsMatchedResouce(resource) {
				log.Printf("found %#v\n", resource)
			}
		}
	}
}

type ResourceInfo struct {
	ClusterName string
	Name        string
	Namespace   string
}

type Filter struct {
	Name             string `yaml:"name"`
	ClusterRegex     string `yaml:"clusterRegex"`
	NamespaceRegex   string `yaml:"namespaceRegex"`
	ResouceNameRegex string `yaml:"resourceNameRegex"`

	compiledClusterRegex      *regexp.Regexp
	compiledNamespaceRegex    *regexp.Regexp
	compiledResourceNameRegex *regexp.Regexp
}

func (f *Filter) NewFilter() (*Filter, error) {
	var err error

	if f.ClusterRegex == "" && f.NamespaceRegex == "" && f.ResouceNameRegex == "" {
		return nil, fmt.Errorf("at least one of ClusterRegex, NamespaceRegex, or PodNameRegex must be provided")
	}

	if f.ClusterRegex != "" {
		f.compiledClusterRegex, err = regexp.Compile(f.ClusterRegex)
		if err != nil {
			return nil, err
		}
	}
	if f.NamespaceRegex != "" {
		f.compiledNamespaceRegex, err = regexp.Compile(f.NamespaceRegex)
		if err != nil {
			return nil, err
		}
	}
	if f.ResouceNameRegex != "" {
		f.compiledResourceNameRegex, err = regexp.Compile(f.ResouceNameRegex)
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

type Config struct {
	mux              sync.Mutex
	filters          []Filter
	lastModifiedTime time.Time
	fileName         string
	watchInterval    time.Duration
}

type ConfigOption func(*Config)

func WithFileName(fileName string) ConfigOption {
	return func(c *Config) {
		c.fileName = fileName
	}
}

func WithWatchInterval(interval time.Duration) ConfigOption {
	return func(c *Config) {
		c.watchInterval = interval
	}
}

func NewConfig(opts ...ConfigOption) *Config {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Config) ReadConfig() error {
	data, err := os.ReadFile(c.fileName)
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	var config struct {
		Filters []Filter `yaml:"filters"`
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("Error unmarshaling config: %v", err)
	}

	var newFilters []Filter
	for _, rawFilter := range config.Filters {
		filter, err := rawFilter.NewFilter()
		if err != nil {
			return err
		}
		newFilters = append(newFilters, *filter)
	}

	fileInfo, err := os.Stat(c.fileName)
	if err != nil {
		return fmt.Errorf("Error stating config file: %v", err)
	}

	c.mux.Lock()
	c.filters = newFilters
	c.lastModifiedTime = fileInfo.ModTime()
	c.mux.Unlock()
	return nil
}

func (c *Config) CheckAndReloadConfig() {
	fileInfo, err := os.Stat(c.fileName)
	if err != nil {
		log.Println("Error stating config file:", err)
		return
	}
	currentModifiedTime := fileInfo.ModTime()
	if currentModifiedTime.After(c.lastModifiedTime) {
		err := c.ReadConfig()
		if err != nil {
			log.Println("Error re-loading config:", err)
		}

		log.Println("Config re-loaded successfully.")
	}
}

func (c *Config) StartConfigWatcher() {
	ticker := time.NewTicker(c.watchInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.CheckAndReloadConfig()
		}
	}
}

func (c *Config) Filters() []Filter {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.filters
}

func (c *Config) IsMatchedResouce(resource ResourceInfo) bool {
	for _, filter := range c.Filters() {
		matchCluster := filter.ClusterRegex == "" || (filter.compiledClusterRegex != nil && filter.compiledClusterRegex.MatchString(resource.ClusterName))
		matchNamespace := filter.NamespaceRegex == "" || (filter.compiledNamespaceRegex != nil && filter.compiledNamespaceRegex.MatchString(resource.Namespace))
		matchName := filter.ResouceNameRegex == "" || (filter.compiledResourceNameRegex != nil && filter.compiledResourceNameRegex.MatchString(resource.Name))

		if matchCluster && matchNamespace && matchName {
			return true
		}
	}
	return false
}
