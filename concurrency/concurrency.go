package concurrency

import (
	"fmt"
	"net/http"
)

func CheckURL(url string, ch chan string) {
	var status string
	defer func() {
		ch <- status
	}()
	_, err := http.Get(url)
	if err != nil {
		status = fmt.Sprintf("%s is unresponsive...", url)
		return
	}
	status = fmt.Sprintf("%s is live", url)
}
