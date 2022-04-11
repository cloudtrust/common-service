package events

import (
	"github.com/Shopify/sarama"
)

type NoopKafkaProducer struct{}

// noop
func (n *NoopKafkaProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return 0, 0, nil
}

// noop
func (n *NoopKafkaProducer) SendMessages(msgs []*sarama.ProducerMessage) error { return nil }

// noop
func (n *NoopKafkaProducer) Close() error { return nil }
