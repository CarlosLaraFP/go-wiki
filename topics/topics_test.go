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
	topic.Send("Hello Kafka!")
	wg.Wait()
}

func TestSendMessages(t *testing.T) {
	clear(Topics)
	_, exists := Topics[4]
	assert.Equal(t, false, exists)

	fmt.Println("Sending random messages...")
	SendMessages(4, 3)

	_, exists = Topics[4]
	assert.Equal(t, true, exists)
}
