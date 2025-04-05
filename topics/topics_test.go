package topics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTopic(t *testing.T) {
	_, exists := Topics[5]
	assert.Equal(t, false, exists)
	topic, _ := NewTopic[string](5, 3)
	assert.Equal(t, 5, topic.Id)
	assert.Equal(t, 3, topic.PartitionCount)
	assert.Equal(t, 3, len(topic.Partitions))
	_, exists = Topics[5]
	assert.Equal(t, false, exists)
}

func TestSend(t *testing.T) {
	topic, _ := NewTopic[string](5, 3)
	wg.Add(1)
	topic.Send("12345", "Hello Kafka!")
	wg.Wait()
}

func TestSendMessages(t *testing.T) {
	clear(Topics)
	id := 4
	m := []string{"12345", "67890", "555550"}

	_, exists := Topics[id]
	assert.Equal(t, false, exists)

	fmt.Println("Sending user messages...")
	SendMessages(id, m)

	_, exists = Topics[id]
	assert.Equal(t, true, exists)
}

func TestPartitionConsistency(t *testing.T) {
	topic, _ := NewTopic[string](1, 3)
	key := "user123"
	p1 := hash(key) % topic.PartitionCount
	p2 := hash(key) % topic.PartitionCount
	assert.Equal(t, p1, p2, "same key should map to same partition")
}
