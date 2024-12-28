package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// AppConfig 应用配置结构体
type AppConfig struct {
	Fee   float64
	Limit int
}

// ConfigManager 配置管理器
type ConfigManager[T any] struct {
	scanInterval time.Duration
	mx           sync.Mutex
	dynCfg       T
	configFile   string
}

func New[T any](configFile string, initCfg T, scanInterval time.Duration) *ConfigManager[T] {
	return &ConfigManager[T]{
		scanInterval: scanInterval,
		dynCfg:       initCfg,
		configFile:   configFile,
	}
}

func (m *ConfigManager[T]) LoadConfig() error {
	file, err := os.Open(m.configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var cfg T
	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		return err
	}

	m.mx.Lock()
	defer m.mx.Unlock()
	m.dynCfg = cfg
	return nil
}

func (m *ConfigManager[T]) GetConfig() T {
	m.mx.Lock()
	defer m.mx.Unlock()
	return m.dynCfg
}

func (m *ConfigManager[T]) Run(ctx context.Context) error {
	if err := m.LoadConfig(); err != nil {
		return err
	}
	go func() {
		timer := time.NewTimer(m.scanInterval)
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				if err := m.LoadConfig(); err != nil {
					log.Printf("Error while loading the config: %v", err)
				}

				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(m.scanInterval)
			}
		}
	}()
	return nil
}

func main() {
	// 初始化 AppConfig
	initCfg := AppConfig{Fee: 0.1, Limit: 10}
	// 初始化 ConfigManager，使用 config.json 作为配置文件
	manager := New("config.json", initCfg, 5*time.Second)

	if err := manager.Run(context.TODO()); err != nil {
		log.Fatal(err)
	}

	config := manager.GetConfig()
	fmt.Printf("当前配置: Fee: %.2f, Limit: %d\n", config.Fee, config.Limit)

	// 模拟更新配置文件，这里可以手动修改 config.json 文件
	// 例如将 Fee 修改为 0.05，Limit 修改为 20
	// 然后等待一段时间，让后台任务更新配置
	time.Sleep(10 * time.Second)

	// 再次获取最新配置
	config = manager.GetConfig()
	fmt.Printf("更新后的配置: Fee: %.2f, Limit: %d\n", config.Fee, config.Limit)
}
