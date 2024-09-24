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
	if err := e.Delete(key); err != nil {
		t.Fatalf("expected error after deleting key, got nil")
	}

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

	err := e.Delete(key)
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

	err = e.Delete("nonExistentKey")
	if err == nil {
		t.Fatalf("expected error when deleting non-existent key, got nil")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound when deleting non-existent key, got %v", err)
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
