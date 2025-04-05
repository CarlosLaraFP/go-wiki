package topics

import (
	"fmt"
	"time"
)

type Topic[T any] struct {
	Id             int
	PartitionCount int `json:"partitionCount"`
	Partitions     []chan T
}

var topics = make(map[int]*Topic[string])
var partitionCount = 3

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
				time.Sleep(100 * time.Millisecond) // simulates work
				fmt.Printf("topic id %d, partition id %d, message: %v\n", t.Id, id, m)
			}
		}(i, t.Partitions[i])
	}
	return t, nil
}

// Sends message m of type T to Topic id
func (t *Topic[T]) Send(m T) error {
	// Router
	return nil
}

// Sends random messages to Topic id every s seconds
func SendMessages(id, s int) {
	for {
		if _, exists := topics[id]; !exists {
			t, _ := NewTopic[string](id, partitionCount)
			topics[id] = t
		}
		t := topics[id]
		t.Send(fmt.Sprintf("current time is %s", time.Now().Local().Format(time.RFC3339)))
		time.Sleep(time.Duration(s) * time.Second)
	}
}
