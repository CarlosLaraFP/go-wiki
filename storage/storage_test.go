package storage

import (
	"testing"
)

func TestStorage(t *testing.T) {
	ms := &MemoryStorage{cache: make(map[string]string)}
	ms.Put("hello", "world")
	_, exists := ms.Get("hello")
	if !exists {
		t.FailNow()
	}

	ls := New(ms)
	ls.Put("h", "v")
	_, exists = ls.Get("h")
	if !exists {
		t.FailNow()
	}

	if _, exists := ls.Get("k"); exists {
		t.FailNow()
	}
}
