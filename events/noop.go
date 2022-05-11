package events

import (
	"context"

	"github.com/Shopify/sarama"
)

// NoopKafkaConsumerGroup is an consumer group that does nothing.
type NoopKafkaConsumerGroup struct{}

// noop
func (n *NoopKafkaConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	return nil
}

// noop
func (n *NoopKafkaConsumerGroup) Errors() <-chan error {
	return make(<-chan error)
}

// noop
func (n *NoopKafkaConsumerGroup) Close() error {
	return nil
}

// noop
func (n *NoopKafkaConsumerGroup) Pause(partitions map[string][]int32) {}

// noop
func (n *NoopKafkaConsumerGroup) Resume(partitions map[string][]int32) {}

// noop
func (n *NoopKafkaConsumerGroup) PauseAll() {}

// noop
func (n *NoopKafkaConsumerGroup) ResumeAll() {}

type NoopKafkaProducer struct{}

// noop
func (n *NoopKafkaProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return 0, 0, nil
}

// noop
func (n *NoopKafkaProducer) SendMessages(msgs []*sarama.ProducerMessage) error { return nil }

// noop
func (n *NoopKafkaProducer) Close() error { return nil }
