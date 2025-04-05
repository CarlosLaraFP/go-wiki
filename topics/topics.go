package topics

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Topic[T any] struct {
	Id             int
	PartitionCount int `json:"partitionCount"`
	Partitions     []chan T
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
		t.Partitions[i] = make(chan T) // should each channel be buffered?
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

// Sends message m of type T to Topic id
func (t *Topic[T]) Send(m T) error {
	// TODO: hash & router
	p := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(t.PartitionCount)
	t.Partitions[p] <- m
	return nil
}

// Sends n random messages to Topic id
func SendMessages(id, n int) {
	if _, exists := Topics[id]; !exists {
		t, _ := NewTopic[string](id, partitionCount)
		Topics[id] = t
	}
	t := Topics[id]

	for range n {
		wg.Add(1)
		t.Send(fmt.Sprintf("current time is %s", time.Now().Local().Format(time.RFC3339)))
	}
	wg.Wait() // wait until all n messages are processed by the Topic
}
