package engine

import (
	"fmt"
	"github.com/patyukin/mdb/pkg/utils"
	"sync"
)

var (
	ErrNotFound = fmt.Errorf("key not found")
)

type Engine struct {
	data map[string]string
	mu   sync.RWMutex
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

func (e *Engine) Delete(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.data, key)
}

func (e *Engine) GetByPattern(pattern string) (map[string]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make(map[string]string)
	for key, value := range e.data {
		match, err := utils.MatchPattern(key, pattern)
		if err != nil {
			return nil, err
		}
		if match {
			result[key] = value
		}
	}

	if len(result) == 0 {
		return nil, ErrNotFound
	}

	return result, nil
}

func (e *Engine) Del(key string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.data[key]
	if !ok {
		return ErrNotFound
	}

	delete(e.data, key)

	return nil
}

func (e *Engine) DelByPattern(pattern string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	var keysToDelete []string
	for key := range e.data {
		match, err := utils.MatchPattern(key, pattern)
		if err != nil {
			return err
		}

		if match {
			keysToDelete = append(keysToDelete, key)
		}
	}

	if len(keysToDelete) == 0 {
		return ErrNotFound
	}

	for _, key := range keysToDelete {
		delete(e.data, key)
	}

	return nil
}
