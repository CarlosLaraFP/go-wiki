package main

import (
	"context"
	"fmt"
	c "go-wiki/concurrency"
	"go-wiki/elasticsearch"
	"log"
	"strconv"

	//i "go-wiki/interfaces"
	//"net/http"
	"os"
)

func main() {
	wp, err := c.NewWorkerPool[string](3, 10)
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

	request := c.Request[string]{
		Context:    context.Background(),
		WorkerPool: wp,
		DLQueue:    dlq,
		Log:        make(chan string, capacity),
	}
	c.ProcessResources(request, ids)

	for r := range request.Log {
		fmt.Println(r)
	}

	if len(dlq.Failed) > 0 {
		fmt.Printf("%d messages failed to process", len(dlq.Failed))
	}

	// ElasticSearch
	years := []int{2020, 2021, 2022, 2023, 2024, 2025}
	//books := make([]*elasticsearch.Book, len(years))
	// TODO: dependency injection
	client, err := elasticsearch.CreateClient()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//defer client.ClosePointInTime()

	for i, y := range years {
		b := &elasticsearch.Book{
			Title:  strconv.Itoa(i),
			Author: strconv.Itoa(i + 1),
			Price:  10.99,
			Year:   y,
			Rating: 5.0,
			Genre:  "Action",
		}
		if err := b.Index(client); err != nil {
			log.Printf("Failed to index book %d: %v", i, err)
		}
	}

	if err := elasticsearch.SearchBooks(client, "Book 0"); err != nil {
		log.Printf("Search failed: %v", err)
	}

	if err := elasticsearch.AggregateRatingsByGenre(client); err != nil {
		log.Printf("Aggregation failed: %v", err)
	}
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

// go mod init go-wiki
// go run main.go
// go mod tidy
// go test ./... -v
