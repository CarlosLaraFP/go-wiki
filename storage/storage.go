/*
Create a Storage interface with Put(key string, value string) and Get(key string) (string, bool).
Implement two concrete types:

MemoryStorage (a map-based implementation)
LoggingStorage (wraps another Storage and logs all operations)
Use the Storage interface to store and retrieve a value.

Follow-up if time: Add unit tests using Go’s testing package.
*/

package storage

import "fmt"

type Storage[T any] interface {
	Put(key string, value T)
	Get(key string) (T, bool)
	Delete(key string) bool
}

type MemoryStorage[T any] struct {
	cache map[string]T
}

func (m *MemoryStorage[T]) Put(key string, value T) {
	m.cache[key] = value
}

func (m *MemoryStorage[T]) Get(key string) (T, bool) {
	value, exists := m.cache[key]
	return value, exists
}

func (m *MemoryStorage[T]) Delete(key string) bool {
	if _, exists := m.cache[key]; exists {
		delete(m.cache, key)
		return true
	}
	return false
}

type LoggingStorage[T any] struct {
	storage Storage[T]
}

func (l *LoggingStorage[T]) Put(key string, value T) {
	l.storage.Put(key, value)
	fmt.Printf("[PUT] key=%s value=%v\n", key, value)
}

func (l *LoggingStorage[T]) Get(key string) (T, bool) {
	value, exists := l.storage.Get(key)
	if exists {
		fmt.Printf("[GET] value: %v\n", value)
	} else {
		fmt.Printf("Key '%s' does not exist\n", key)
	}
	return value, exists
}

func (l *LoggingStorage[T]) Delete(key string) {
	if l.storage.Delete(key) {
		fmt.Printf("Key '%s' successfully deleted", key)
	} else {
		fmt.Printf("Key '%s' does not exist", key)
	}
}

// When the parameter is an interface, it automatically expects a pointer
func New[T any](storage Storage[T]) *LoggingStorage[T] {
	l := LoggingStorage[T]{storage}
	return &l
}
