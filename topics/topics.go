// Package topics simulates Kafka-style partitioning in Go.
// Key concepts:
// - Topics split into partitions (channels) for parallel processing.
// - Keys hash to consistent partitions (like Kafka's sticky partitioning).
// - Buffered channels prevent producer blocking (backpressure).
package topics

import (
	"fmt"
	"sync"
	"time"
)

type Topic[T any] struct {
	Id             int      `json:"id"`
	PartitionCount int      `json:"partitionCount"`
	Partitions     []chan T `json:"partitions"`
}

var Topics = make(map[int]*Topic[string])
var partitionCount = 3
var wg = &sync.WaitGroup{}

// NewTopic creates a Topic with id and n partitions
func NewTopic[T any](id, n int) (*Topic[T], error) {
	t := &Topic[T]{
		Id:             id,
		PartitionCount: n,
		Partitions:     make([]chan T, n),
	}
	for i := range t.PartitionCount {
		// Kafka partitions are buffered (in-memory queues)
		// Prevents backpressure from blocking the producer
		// If channels are unbuffered (make(chan T)), the producer blocks on t.Partitions[p] <- m until a consumer reads.
		t.Partitions[i] = make(chan T, 100)

		go func(id int, c chan T) {
			for m := range c {
				time.Sleep(1 * time.Second)                                            // simulates work
				fmt.Printf("topic id %d, partition id %d, message: %v\n", t.Id, id, m) // send back into a receiver channel?
				wg.Done()
			}
		}(i, t.Partitions[i])
	}
	return t, nil
}

func (t *Topic[T]) Delete() {
	for _, p := range t.Partitions {
		close(p)
	}
}

// Sends message m of type T to Topic id
// By applying hash(key) % t.PartitionCount, we constrain the result to:
// Minimum: 0 (when hash mod = 0)
// Maximum: t.PartitionCount - 1 (when hash mod = PartitionCount-1)
func (t *Topic[T]) Send(key string, m T) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	// Kafka also uses message keys to assign partitions
	p := hash(key) % t.PartitionCount // Consistent hashing
	t.Partitions[p] <- m
	return nil
}

// hash converts any string to a (potentially large) integer
func hash(s string) int {
	h := 0
	for _, c := range s {
		h = 31*h + int(c)
	}
	return h
}

// Sends messages to Topic id
func SendMessages(id int, keys []string) {
	if _, exists := Topics[id]; !exists {
		t, _ := NewTopic[string](id, partitionCount)
		Topics[id] = t
	}
	t := Topics[id]

	for _, k := range keys {
		wg.Add(1)
		t.Send(k, fmt.Sprintf("current time is %s", time.Now().Local().Format(time.RFC3339)))
	}
	wg.Wait() // wait until all messages are processed by the Topic
}
