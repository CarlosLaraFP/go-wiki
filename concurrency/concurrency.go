package concurrency

import (
	"fmt"
	"net/http"
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
