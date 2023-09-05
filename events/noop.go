package events

import (
	"context"

	"github.com/IBM/sarama"
)

// NoopKafkaConsumerGroup is an consumer group that does nothing.
type NoopKafkaConsumerGroup struct{}

// Consume does noop
func (n *NoopKafkaConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	return nil
}

// Errors creates a channel for errors
func (n *NoopKafkaConsumerGroup) Errors() <-chan error {
	return make(<-chan error)
}

// Close does noop
func (n *NoopKafkaConsumerGroup) Close() error {
	return nil
}

// Pause does noop
func (n *NoopKafkaConsumerGroup) Pause(partitions map[string][]int32) {}

// Resume does noop
func (n *NoopKafkaConsumerGroup) Resume(partitions map[string][]int32) {}

// PauseAll does noop
func (n *NoopKafkaConsumerGroup) PauseAll() {}

// ResumeAll does noop
func (n *NoopKafkaConsumerGroup) ResumeAll() {}

// NoopKafkaProducer struct
type NoopKafkaProducer struct{}

// SendMessage does noop
func (n *NoopKafkaProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return 0, 0, nil
}

// SendMessages does noop
func (n *NoopKafkaProducer) SendMessages(msgs []*sarama.ProducerMessage) error { return nil }

// Close does noop
func (n *NoopKafkaProducer) Close() error { return nil }

// noop
func (n *NoopKafkaProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	return 0
}

// noop
func (n *NoopKafkaProducer) IsTransactional() bool { return true }

// noop
func (n *NoopKafkaProducer) BeginTxn() error { return nil }

// noop
func (n *NoopKafkaProducer) CommitTxn() error { return nil }

// noop
func (n *NoopKafkaProducer) AbortTxn() error { return nil }

// noop
func (n *NoopKafkaProducer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	return nil
}

// noop
func (n *NoopKafkaProducer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	return nil
}
