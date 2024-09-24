package engine

import (
	"fmt"
	"sync"
)

var (
	ErrNotFound = fmt.Errorf("key not found")
)

type Engine struct {
	mu   sync.RWMutex
	data map[string]string
}

func New() *Engine {
	return &Engine{
		data: make(map[string]string),
	}
}

func (e *Engine) Set(key string, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.data[key] = value
}

func (e *Engine) Get(key string) (string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	value, exists := e.data[key]
	if !exists {
		return "", fmt.Errorf("'%s' - %w", key, ErrNotFound)
	}

	return value, nil
}

func (e *Engine) Delete(key string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.data[key]
	if !ok {
		return fmt.Errorf("'%s' - %w", key, ErrNotFound)
	}

	delete(e.data, key)

	return nil
}
