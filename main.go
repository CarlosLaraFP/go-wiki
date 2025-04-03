package main

import (
	"fmt"
	"go-wiki/concurrency"
)

func main() {
	urls := []string{
		"http://google.com",
		"http://facebook.com",
		"http://golang.org",
		"http://amazon.com",
		"fake",
	}

	ch := make(chan string, len(urls))

	for _, url := range urls {
		go concurrency.CheckURL(url, ch)
	}

	for range urls {
		fmt.Println(<-ch)
	}

	//web.Serve()
}

// go mod init go-wiki
// go run main.go
// go mod tidy
// go test ./...
