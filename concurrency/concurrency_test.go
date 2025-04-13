package concurrency

import (
	"context"
	"fmt"
	"sync"
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
	ch := make(chan string, 10)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	err := process(ctx, 100, ch, time.Millisecond*250, 0)
	assert.NoError(t, err)
	fmt.Println(err)
	err = process(ctx, 10, ch, time.Millisecond*250, 3)
	assert.NoError(t, err)
	cancel()
	ctx, cancel = context.WithTimeout(context.TODO(), 5*time.Millisecond)
	defer cancel()
	time.Sleep(200 * time.Millisecond)
	err = process(ctx, 100, ch, time.Millisecond*250, 0)
	assert.Error(t, err)
}

func TestMessageProcessor(t *testing.T) {
	log := make(chan string, 10)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}
	dlq := &DeadLetterQueue[float64]{}
	m := &Job[float64]{3.14159, ctx, dlq, log, wg}
	ch := make(chan *Job[float64])
	go messageProcessor(ch)
	time.Sleep(200 * time.Millisecond)
	wg.Add(1)
	ch <- m
	wg.Wait()
}

func TestProcessResourceIds(t *testing.T) {
	wp, err := NewWorkerPool[string](5, 10)
	assert.NoError(t, err)
	defer wp.Cleanup()

	dlq := &DeadLetterQueue[string]{}

	capacity := 20
	ids := make([]string, 0, capacity)
	for i := range capacity {
		ids = append(ids, fmt.Sprintf("event-%d", i))
	}
	request := Request[string]{
		Context:    context.Background(),
		WorkerPool: wp,
		DLQueue:    dlq,
		Log:        make(chan string, capacity), // no capacity == deadlock
	}
	ProcessResources(request, ids)

	for r := range request.Log {
		fmt.Println(r)
	}

	assert.LessOrEqual(t, len(dlq.Failed), 5)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	copy := make([]string, 0, capacity)

	for i := range capacity {
		copy = append(copy, fmt.Sprintf("event-%d", i))
	}

	new := Request[string]{
		Context:    ctx,
		WorkerPool: wp,
		DLQueue:    dlq,
		Log:        make(chan string, capacity), // no capacity == deadlock
	}

	ProcessResources(new, copy)
	assert.NotEmpty(t, dlq.Failed)
}

func TestScaleUp(t *testing.T) {
	wp, _ := NewWorkerPool[string](5, 10)
	defer wp.Cleanup()

	wp.ScaleUp(5, 20)
	assert.Equal(t, 10, len(wp.Workers))
}

func TestScaleDown(t *testing.T) {
	wp, _ := NewWorkerPool[string](5, 10)
	defer wp.Cleanup()

	closed := wp.ScaleDown(4)
	assert.Equal(t, 4, closed)
}
