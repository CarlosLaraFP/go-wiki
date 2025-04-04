package main

import (
	c "go-wiki/concurrency"
	"time"
)

func main() {
	urls := []string{
		"http://google.com",
		"http://facebook.com",
		"http://golang.org",
		"http://amazon.com",
		"fake",
	}

	c.LaunchWorkerPool([]string{"1.txt", "2.txt", "3.txt"}, 3)

	time.Sleep(30 * time.Second)

	ch := make(chan string, len(urls))

	for _, url := range urls {
		go c.CheckURL(url, ch)
	}
	// for each iteration, wait to receive a value from the channel
	for url := range ch {
		go c.CheckURL(url, ch)
	}

	//web.Serve()
}

// go mod init go-wiki
// go run main.go
// go mod tidy
// go test ./...
