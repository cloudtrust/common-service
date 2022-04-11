package events

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	noopProducer := NoopKafkaProducer{}
	msg := sarama.ProducerMessage{Topic: "test"}
	partition, offset, err := noopProducer.SendMessage(&msg)

	assert.Nil(t, err)
	assert.Zero(t, partition)
	assert.Zero(t, offset)
}

func TestSendMessages(t *testing.T) {
	noopProducer := NoopKafkaProducer{}
	msg := sarama.ProducerMessage{Topic: "test"}
	err := noopProducer.SendMessages([]*sarama.ProducerMessage{&msg})

	assert.Nil(t, err)
}

func TestClose(t *testing.T) {
	noopProducer := NoopKafkaProducer{}
	err := noopProducer.Close()
	assert.Nil(t, err)
}
