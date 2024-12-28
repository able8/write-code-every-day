package manager

import (
	"context"
	"log"
	"sync"
	"time"
)

type ConfigManager[T any] struct {
	repo         Repository
	scanInterval time.Duration
	mx           sync.Mutex
	dynCfg       T
}

func New[T any](repo Repository, initCfg T, scanInterval time.Duration) *ConfigManager[T] {
	return &ConfigManager[T]{
		dynCfg:       initCfg,
		repo:         repo,
		scanInterval: scanInterval,
	}
}

func (m *ConfigManager[T]) LoadConfig(ctx context.Context) error {
	var cfg T
	if err := m.repo.GetConfig(ctx, &cfg, m.dynCfg); err != nil {
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

func (m *ConfigManager[T]) Run(ctx context.Context, wg *sync.WaitGroup) error {
	if err := m.LoadConfig(ctx); err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(m.scanInterval):
				if err := m.LoadConfig(ctx); err != nil {
					log.Printf("error while loading the config: %v", err)
				}
			}
		}
	}()
	return nil
}
