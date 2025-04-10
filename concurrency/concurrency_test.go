package concurrency

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func TestProcess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	err := process(ctx, 100, time.Millisecond*250, 0)
	assert.NoError(t, err)
	fmt.Println(err)
	err = process(ctx, 10, time.Millisecond*250, 3)
	assert.NoError(t, err)
	cancel()
	ctx, cancel = context.WithTimeout(context.TODO(), 5*time.Millisecond)
	defer cancel()
	time.Sleep(200 * time.Millisecond)
	err = process(ctx, 100, time.Millisecond*250, 0)
	assert.Error(t, err)
}
