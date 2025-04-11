package main

import (
	"context"
	"fmt"
	c "go-wiki/concurrency"

	//i "go-wiki/interfaces"
	//"net/http"
	"os"
)

func main() {
	wp, err := c.NewWorkerPool[string](5, 10)
	if err != nil {
		fmt.Printf("error setting up worker pool: %v\n", err)
		os.Exit(1)
	}
	defer wp.Cleanup()

	dlq := &c.DeadLetterQueue[string]{}

	capacity := 20
	ids := make([]string, 0, capacity)
	for i := range capacity {
		ids = append(ids, fmt.Sprintf("event-%d", i))
	}
	c.ProcessResourceIds(context.Background(), wp, dlq, ids)

	if len(dlq.Failed) > 0 {
		fmt.Printf("%d messages failed to process", len(dlq.Failed))
	}

	/*
		s := i.Square{Length: 4}
		t := i.Triangle{Base: 5, Height: 4}

		i.ShowShape(&s)
		i.ShowShape(&t)

		urls := []string{
			"http://google.com",
			"http://facebook.com",
			"http://golang.org",
			"http://amazon.com",
			"fake",
		}

		_, err := http.Get("http://google.com")
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		//io.Copy(os.Stdout, rsp.Body)

		c.LaunchWorkerPool([]string{"1.txt", "2.txt", "3.txt"}, 3)

		time.Sleep(30 * time.Second)

		ch := make(chan string, len(urls))

		for _, url := range urls {
			go c.CheckURL(url, ch)
		}
		// for each iteration, wait to receive a value from the channel
		//for url := range ch {
		//	go c.CheckURL(url, ch)
		//}

		//web.Serve()
	*/
}

// go mod init go-wiki
// go run main.go
// go mod tidy
// go test ./...
