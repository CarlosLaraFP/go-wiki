package concurrency

import (
	"fmt"
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
