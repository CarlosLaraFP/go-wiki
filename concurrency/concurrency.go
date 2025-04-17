package concurrency

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
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

type Job[T comparable] struct {
	message T
	ctx     context.Context
	dlq     *DeadLetterQueue[T]
	log     chan string
	wg      *sync.WaitGroup
	ctr     *atomic.Int32
}

type Worker[T comparable] struct {
	Id    int
	Queue chan *Job[T]
}

type WorkerPool[T comparable] struct {
	sync.RWMutex
	Jobs    *atomic.Int32
	Workers []Worker[T]
	Alive   bool
}

func (wp *WorkerPool[T]) Dispatch(m *Job[T]) {
	wp.Lock()
	defer wp.Unlock()
	// Random dispatcher fans out jobs evenly to available workers (Fan-Out)
	i := rand.Intn(len(wp.Workers))
	wp.Workers[i].Queue <- m
	wp.Jobs.Add(1)
}

// StartAutoscalerController simulates an autoscaling controller
func (wp *WorkerPool[T]) StartAutoscalerController(ctx context.Context, cap int) {
	select {
	case <-ctx.Done():
		return
	default:
		// TODO: Make sleep a percentage of remaining context time
		time.Sleep(2 * time.Second)
		wp.ScaleUp(2, cap)
		time.Sleep(2 * time.Second)
		_ = wp.ScaleDown(2)
	}
}

// StartLogger prints the number of jobs in progress every d
func (wp *WorkerPool[T]) StartLogger(d time.Duration) {
	for wp.Alive {
		time.Sleep(d)
		fmt.Printf("%d jobs in progress\n", wp.Jobs.Load())
	}
}

func (wp *WorkerPool[T]) Cleanup() {
	wp.Alive = false
	wp.Jobs.Store(0)

	for _, w := range wp.Workers {
		select {
		case _, ok := <-w.Queue:
			if ok {
				close(w.Queue)
			}
		default:
		}
	}
}

// NewWorkerPool accepts n worker count and c buffer capacity and returns a WorkerPool
func NewWorkerPool[T comparable](n, c int) (*WorkerPool[T], error) {
	var w []Worker[T]
	for i := range n {
		w = append(w, Worker[T]{
			Id:    i,
			Queue: make(chan *Job[T], c),
		})
		go messageProcessor(w[i].Queue)
	}
	wp := WorkerPool[T]{Jobs: &atomic.Int32{}, Workers: w, Alive: true}
	wp.Jobs.Store(0)

	return &wp, nil
}

func messageProcessor[T comparable](ch chan *Job[T]) {
	for m := range ch {
		if err := process(m.ctx, m.message, m.log, time.Millisecond*200, 1); err != nil {
			fmt.Println(err)
			m.dlq.Add(m.message)
		}
		m.ctr.Add(-1)
		m.wg.Done()
	}
}

// process respects the context deadline
func process[T comparable](ctx context.Context, m T, l chan string, d time.Duration, retry int) error {
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

	l <- fmt.Sprintf("Processed message: %v", m)

	return nil
}

type Request[T comparable] struct {
	Messages       []T
	Context        context.Context
	WorkerPool     *WorkerPool[T]
	MaxParallelism int
	DLQueue        *DeadLetterQueue[T]
	Log            chan string
}

// ProcessRequest processes a slice of messages concurrently using a WorkerPool
// The entire process respects a context.Context timeout (e.g., 5 seconds).
// If the context is canceled (timeout hit), workers stop immediately.
// The total time receiving from the Log also counts.
func ProcessRequest[T comparable](r Request[T]) {
	wg := &sync.WaitGroup{}
	c, cancel := context.WithTimeout(r.Context, 5*time.Second)
	defer cancel()
	defer close(r.Log)

	go r.WorkerPool.StartAutoscalerController(c, cap(r.Log))
	go r.WorkerPool.StartLogger(500 * time.Millisecond)

	for _, m := range r.Messages {
		for r.WorkerPool.Jobs.Load() >= int32(r.MaxParallelism) {
			time.Sleep(50 * time.Millisecond)
		}
		wg.Add(1)
		r.WorkerPool.Dispatch(&Job[T]{m, c, r.DLQueue, r.Log, wg, r.WorkerPool.Jobs})
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
	if wp.Alive {
		wp.Lock()
		defer wp.Unlock()
		l := len(wp.Workers)

		for i := range n {
			wp.Workers = append(wp.Workers, Worker[T]{
				Id:    i + l,
				Queue: make(chan *Job[T], c),
			})
			go messageProcessor(wp.Workers[i+l].Queue)
		}
		fmt.Printf("Added %d workers to the pool\n", n)
	}
}

// ScaleDown dynamically removes n workers from the pool by closing their channels.
// Asynchronous behavior: Only closed or empty channels can be removed (waits for non-empty channels to drain).
func (wp *WorkerPool[T]) ScaleDown(n int) int {
	wp.Lock()
	defer wp.Unlock()
	if n > len(wp.Workers) {
		fmt.Printf("%d is greater than the number of workers: %d. Use WorkerPool.Cleanup() instead", n, len(wp.Workers))
		return 0
	}
	closed := 0
	for wp.Alive && closed < n {
		for _, w := range wp.Workers {
			if closed == n {
				break
			}
			select {
			case _, ok := <-w.Queue: // non-blocking check
				if !ok {
					// Channel is already closed
					closed++
				}
				// Channel still active; do nothing
			default:
				// Closing the channel and letting it drain on its own if it still contains messages
				close(w.Queue)
				closed++
			}
		}
		fmt.Printf("Removed %d workers from the pool by closing their channels\n", closed)
	}
	return closed
}
