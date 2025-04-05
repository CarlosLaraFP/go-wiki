package topics

import "fmt"

type Topic[T any] struct {
	Id             int
	PartitionCount int `json:"partitionCount"`
	Partitions     []chan T
}

// NewTopic creates a Topic with n partitions
func NewTopic[T any](n int) (*Topic[T], error) {
	t := &Topic[T]{}
	t.PartitionCount = n
	t.Partitions = make([]chan T, t.PartitionCount)
	for i := range t.PartitionCount {
		t.Partitions[i] = make(chan T) // should each channel be buffered?
		go func(id int, c chan T) {
			for m := range c {
				fmt.Printf("topic id %d, partition id %d, message: %v\n", t.Id, id, m)
			}
		}(i, t.Partitions[i])
	}
	return t, nil
}
