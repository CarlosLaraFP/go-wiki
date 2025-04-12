package concurrency

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func CheckURL(url string, ch chan string) {
	defer func() {
		time.Sleep(5 * time.Second)
		ch <- url
	}()
	if _, err := http.Get(url); err != nil {
		fmt.Printf("%s is unresponsive...\n", url)
		return
	}
	fmt.Printf("%s is live\n", url)
}

// Program accepts a slice of filenames.
// Uses a worker pool pattern to process each file:
//
//	Each worker is a goroutine that:
//	Receives a filename from a channel
//	Simulates work with time.Sleep(100 * time.Millisecond)
//	Prints: processed: <filename>
func LaunchWorkerPool(fileNames []string, n int) {
	// Launch N workers (e.g., N := 3)
	// Close the file channel properly
	// Use a sync.WaitGroup to wait for all workers to finish
	// Avoid goroutine leaks or deadlocks
	wg := &sync.WaitGroup{}
	c := make(chan string)

	for i := range n {
		wg.Add(1)
		go func(id int) {
			defer wg.Done() // avoids deadlock even if the goroutine fails
			/*
				If the channel is open but empty → the read blocks (waits for producer to send).
				If the channel is closed:
				If there are buffered items → keep consuming buffered items.
				Once the buffer is drained → exit loop cleanly (no panic, no error).
			*/
			for fn := range c {
				time.Sleep(100 * time.Millisecond)
				fmt.Printf("worker pool %d processed %s\n", id, fn)
			}
		}(i)
	}

	for _, n := range fileNames {
		c <- n
	}
	close(c)
	wg.Wait()
}

////////////////////////////////////////////////////////////////////////

type DeadLetterQueue[T comparable] struct {
	sync.Mutex
	Failed []T
}

func (dlq *DeadLetterQueue[T]) Add(message T) {
	dlq.Lock()
	defer dlq.Unlock()
	dlq.Failed = append(dlq.Failed, message)
}

// fmt.Stringer
type Message[T comparable] struct {
	message T
	ctx     context.Context
	dlq     *DeadLetterQueue[T]
	log     chan T
	wg      *sync.WaitGroup
}

type Worker[T comparable] struct {
	Id    int
	Queue chan *Message[T]
}

type WorkerPool[T comparable] struct {
	Workers []Worker[T]
}

func (wp *WorkerPool[T]) Cleanup() {
	for _, w := range wp.Workers {
		close(w.Queue)
	}
}

// NewWorkerPool accepts n worker count and c buffer capacity and returns a WorkerPool
func NewWorkerPool[T comparable](n, c int) (*WorkerPool[T], error) {
	var w []Worker[T]
	for i := range n {
		w = append(w, Worker[T]{
			Id:    i,
			Queue: make(chan *Message[T], c),
		})
		go messageProcessor(w[i].Queue)
	}
	wp := WorkerPool[T]{Workers: w}
	return &wp, nil
}

func messageProcessor[T comparable](ch chan *Message[T]) {
	for m := range ch {
		if err := process(m.ctx, m.message, m.log, time.Millisecond*200, 1); err != nil {
			fmt.Println(err)
			m.dlq.Add(m.message)
		}
		m.wg.Done()
	}
}

// process respects the context deadline
func process[T comparable](ctx context.Context, m T, l chan T, d time.Duration, retry int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if retry > 3 {
		return fmt.Errorf("3rd retry failed for resource: %v", m)
	}

	time.Sleep(d)

	if n := rand.Float64(); n <= 0.20 {
		fmt.Printf("Simulating random failure for resource: %v. Retrying...\n", m)
		return process(ctx, m, l, d*2, retry+1)
	}

	l <- m
	return nil
}

type Request[T comparable] struct {
	Context    context.Context
	WorkerPool *WorkerPool[T]
	DLQueue    *DeadLetterQueue[T]
	Log        chan T
}

// ProcessResources processes a slice of resources concurrently using a WorkerPool
// The whole process respects a context.Context timeout (e.g., 5 seconds).
// If the context is canceled (timeout hit), workers stop immediately.
// The total time receiving from the Log also counts.
func ProcessResources[T comparable](r Request[T], s []T) {
	wg := &sync.WaitGroup{}
	c, cancel := context.WithTimeout(r.Context, 5*time.Second)
	defer cancel()
	defer close(r.Log)

	for _, m := range s {
		wg.Add(1)
		// Random dispatcher fans out jobs evenly to available workers (Fan-Out)
		i := rand.Intn(len(r.WorkerPool.Workers))
		r.WorkerPool.Workers[i].Queue <- &Message[T]{m, c, r.DLQueue, r.Log, wg}
	}
	wg.Wait()
	/*
		If the context deadline is hit during wg.Wait():
		The context (c) is automatically canceled by Go’s runtime.
		Workers (goroutines) notice ctx.Done() inside process().
		Workers immediately stop processing further retries.
		Workers still call wg.Done() once they return from process().
		wg.Wait() will continue waiting for all workers to call Done().
		After all workers call Done(), wg.Wait() will unblock normally.
	*/
}

// ScaleUp dynamically adds n workers to the pool, each with buffer capacity c
func (wp *WorkerPool[T]) ScaleUp(n, c int) {
	wc := len(wp.Workers)
	for i := range n {
		wp.Workers = append(wp.Workers, Worker[T]{
			Id:    i + wc,
			Queue: make(chan *Message[T], c),
		})
	}
}

// ScaleDown dynamically removes n workers from the pool by closing their channels
// Asynchronous behavior: Only closed or empty channels can be removed (waits for non-empty channels to drain)
func (wp *WorkerPool[T]) ScaleDown(n int) error {
	if n > len(wp.Workers) {
		return fmt.Errorf("%d is greater than the number of workers: %d. Use WorkerPool.Cleanup() instead", n, len(wp.Workers))
	}
	closed := 0
	go func() {
		for closed < n {
			for _, w := range wp.Workers {
				if len(w.Queue) > 0 {
					time.Sleep(1 * time.Second)
					continue
				}
				close(w.Queue)
				closed++
			}
		}
	}()
	return nil
}
