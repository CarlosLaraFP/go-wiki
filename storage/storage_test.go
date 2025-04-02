package storage

import (
	"testing"
)

func TestStorage(t *testing.T) {
	ms := &MemoryStorage[string]{cache: make(map[string]string)}
	ms.Put("hello", "world")
	_, exists := ms.Get("hello")
	if !exists {
		t.Errorf("'hello' key is expected to exist")
	}

	ls := New(ms)
	ls.Put("h", "v")
	_, exists = ls.Get("h")
	if !exists {
		t.Errorf("'h' key is expected to exist")
	}

	if _, exists := ls.Get("k"); exists {
		t.Errorf("'k' key should not exist")
	}
}

/*
FailNow is okay for unresolvable errors, but in most cases you want the test to continue and show all failures.
*/
