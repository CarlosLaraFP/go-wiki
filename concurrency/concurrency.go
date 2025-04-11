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

type DeadLetterQueue[T any] struct {
	sync.Mutex
	Failed []T
}

func (dlq *DeadLetterQueue[T]) Add(message T) {
	dlq.Lock()
	defer dlq.Unlock()
	dlq.Failed = append(dlq.Failed, message)
}

type Message[T any] struct {
	message T
	ctx     context.Context
	wg      *sync.WaitGroup
	dlq     *DeadLetterQueue[T]
}

type Worker[T any] struct {
	Id    int
	Queue chan *Message[T]
}

type WorkerPool[T any] struct {
	Workers []Worker[T]
}

func (wp *WorkerPool[T]) Cleanup() {
	for _, w := range wp.Workers {
		close(w.Queue)
	}
}

// NewWorkerPool accepts n worker count and l buffer size and returns a WorkerPool
func NewWorkerPool[T any](n, l int) (*WorkerPool[T], error) {
	var w []Worker[T]
	for i := range n {
		w = append(w, Worker[T]{
			Id:    i,
			Queue: make(chan *Message[T], l),
		})
		go messageProcessor(w[i].Queue)
	}
	wp := WorkerPool[T]{Workers: w}
	return &wp, nil
}

// process respects the context deadline
func process[T any](ctx context.Context, m T, d time.Duration, retry int) error {
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
		process(ctx, m, d*2, retry+1)
	}

	fmt.Printf("Fetched Resource: %v\n", m)
	return nil
}

func messageProcessor[T any](ch chan *Message[T]) {
	for m := range ch {
		if err := process(m.ctx, m.message, time.Millisecond*200, 1); err != nil {
			fmt.Println(err)
			m.dlq.Add(m.message)
		}
		m.wg.Done()
	}
}

type Request[T any] struct {
	Context    context.Context
	WorkerPool *WorkerPool[T]
	DLQueue    *DeadLetterQueue[T]
}

// ProcessResourceIds blocks until processing is complete.
// The whole process must respect a context.Context timeout (e.g., 5 seconds).
// If the context is canceled (timeout hit), workers must immediately stop.
func ProcessResources[T any](r Request[T], s []T) {
	wg := &sync.WaitGroup{}
	c, cancel := context.WithTimeout(r.Context, 5*time.Second)
	defer cancel()

	for _, m := range s {
		wg.Add(1)
		i := rand.Intn(len(r.WorkerPool.Workers))
		r.WorkerPool.Workers[i].Queue <- &Message[T]{m, c, wg, r.DLQueue}
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
