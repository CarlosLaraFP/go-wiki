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

type Storage interface {
	Put(key, value string)
	Get(key string) (string, bool)
	Delete(key string) bool
}

type MemoryStorage struct {
	cache map[string]string
}

func (m *MemoryStorage) Put(key, value string) {
	m.cache[key] = value
}

func (m *MemoryStorage) Get(key string) (string, bool) {
	value, exists := m.cache[key]
	return value, exists
}

func (m *MemoryStorage) Delete(key string) bool {
	if _, exists := m.cache[key]; exists {
		delete(m.cache, key)
		return true
	}
	return false
}

type LoggingStorage struct {
	storage Storage
}

func (l *LoggingStorage) Put(key, value string) {
	l.storage.Put(key, value)
	fmt.Printf("[PUT] key=%q value=%q\n", key, value)
}

func (l *LoggingStorage) Get(key string) (string, bool) {
	value, exists := l.storage.Get(key)
	if exists {
		fmt.Printf("[GET] value: %s\n", value)
	} else {
		fmt.Printf("Key '%s' does not exist\n", key)
	}
	return value, exists
}

func (l *LoggingStorage) Delete(key string) {
	if l.storage.Delete(key) {
		fmt.Printf("Key '%s' successfully deleted", key)
	} else {
		fmt.Printf("Key '%s' does not exist", key)
	}
}

// When the parameter is an interface, it automatically expects a pointer
func New(storage Storage) *LoggingStorage {
	l := LoggingStorage{storage}
	return &l
}
