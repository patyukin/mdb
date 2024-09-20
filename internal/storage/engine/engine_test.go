package engine

import (
	"errors"
	"fmt"
	"testing"
)

func TestEngine_Set_Get(t *testing.T) {
	e := New()
	key := "testKey"
	value := "testValue"

	e.Set(key, value)

	got, err := e.Get(key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != value {
		t.Fatalf("expected value '%s', got '%s'", value, got)
	}
}

func TestEngine_Get_NotFound(t *testing.T) {
	e := New()
	key := "nonExistentKey"

	_, err := e.Get(key)
	if err == nil {
		t.Fatalf("expected error for non-existent key, got nil")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestEngine_Delete(t *testing.T) {
	e := New()
	key := "testKey"
	value := "testValue"

	e.Set(key, value)
	e.Delete(key)

	_, err := e.Get(key)
	if err == nil {
		t.Fatalf("expected error after deleting key, got nil")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound after deletion, got %v", err)
	}
}

func TestEngine_Del(t *testing.T) {
	e := New()
	key := "testKey"
	value := "testValue"

	e.Set(key, value)

	err := e.Del(key)
	if err != nil {
		t.Fatalf("expected no error when deleting existing key, got %v", err)
	}

	_, err = e.Get(key)
	if err == nil {
		t.Fatalf("expected error after deleting key, got nil")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound after deletion, got %v", err)
	}

	err = e.Del("nonExistentKey")
	if err == nil {
		t.Fatalf("expected error when deleting non-existent key, got nil")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound when deleting non-existent key, got %v", err)
	}
}

func TestEngine_GetByPattern(t *testing.T) {
	e := New()
	keys := []string{"apple", "apricot", "banana", "berry", "blueberry"}
	value := "fruit"

	for _, key := range keys {
		e.Set(key, value)
	}

	pattern := "b*"
	expected := map[string]string{
		"banana":    value,
		"berry":     value,
		"blueberry": value,
	}

	result, err := e.GetByPattern(pattern)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d results, got %d", len(expected), len(result))
	}

	for k, v := range expected {
		if result[k] != v {
			t.Fatalf("expected key '%s' to have value '%s', got '%s'", k, v, result[k])
		}
	}
}

func TestEngine_GetByPattern_NotFound(t *testing.T) {
	e := New()
	e.Set("apple", "fruit")
	e.Set("banana", "fruit")

	pattern := "c*"

	_, err := e.GetByPattern(pattern)
	if err == nil {
		t.Fatalf("expected error when no keys match pattern, got nil")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestEngine_DelByPattern(t *testing.T) {
	e := New()
	keys := []string{"apple", "apricot", "banana", "berry", "blueberry"}
	value := "fruit"

	for _, key := range keys {
		e.Set(key, value)
	}

	pattern := "a*"
	err := e.DelByPattern(pattern)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, key := range keys {
		if key[0] == 'a' {
			_, err := e.Get(key)
			if err == nil {
				t.Fatalf("expected key '%s' to be deleted, but it still exists", key)
			}
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("expected ErrNotFound for key '%s', got %v", key, err)
			}
		} else {
			_, err := e.Get(key)
			if err != nil {
				t.Fatalf("expected key '%s' to exist, got error %v", key, err)
			}
		}
	}
}

func TestEngine_DelByPattern_NotFound(t *testing.T) {
	e := New()
	e.Set("apple", "fruit")
	e.Set("banana", "fruit")

	pattern := "c*"

	err := e.DelByPattern(pattern)
	if err == nil {
		t.Fatalf("expected error when no keys match pattern, got nil")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestEngine_ConcurrentAccess(t *testing.T) {
	e := New()
	key := "concurrentKey"
	value := "initial"

	e.Set(key, value)

	done := make(chan bool)

	for i := 0; i < 100; i++ {
		go func(i int) {
			if i%2 == 0 {
				_, err := e.Get(key)
				if err != nil && !errors.Is(err, ErrNotFound) {
					t.Errorf("unexpected error during Get: %v", err)
				}
			} else {
				e.Set(key, fmt.Sprintf("value%d", i))
			}
			done <- true
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-done
	}

	finalValue, err := e.Get(key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if finalValue == "" {
		t.Fatalf("expected non-empty final value, got empty string")
	}
}
