package concurrency

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	count = 2
	size  = 2
)

func TestNewWorkerPool(t *testing.T) {
	wp, err := NewWorkerPool[string](count, size)
	assert.NoError(t, err)
	assert.Equal(t, count, len(wp.Workers))
	for i, w := range wp.Workers {
		assert.Equal(t, i, w.Id)
		assert.Equal(t, size, cap(w.Queue))
	}
}

func TestCleanup(t *testing.T) {
	wp, err := NewWorkerPool[string](count, size)
	assert.NoError(t, err)
	wp.Cleanup()
	for _, w := range wp.Workers {
		select {
		case _, ok := <-w.Queue:
			if ok {
				assert.Fail(t, "expected channel to be closed for receive")
			}
		default:
		}
	}
}
