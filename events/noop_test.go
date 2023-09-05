package events

import (
	"testing"

	"github.com/IBM/sarama"
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

func TestTxnStatus(t *testing.T) {
	noopProducer := NoopKafkaProducer{}
	txnStatus := noopProducer.TxnStatus()
	assert.Equal(t, sarama.ProducerTxnStatusFlag(0), txnStatus)
}

func TestIsTransactional(t *testing.T) {
	noopProducer := NoopKafkaProducer{}
	b := noopProducer.IsTransactional()
	assert.True(t, b)
}

func TestBeginTxn(t *testing.T) {
	noopProducer := NoopKafkaProducer{}
	err := noopProducer.BeginTxn()
	assert.Nil(t, err)
}

func TestCommitTxn(t *testing.T) {
	noopProducer := NoopKafkaProducer{}
	err := noopProducer.CommitTxn()
	assert.Nil(t, err)
}

func TestAbortTxn(t *testing.T) {
	noopProducer := NoopKafkaProducer{}
	err := noopProducer.AbortTxn()
	assert.Nil(t, err)
}

func TestAddOffsetsToTxn(t *testing.T) {
	noopProducer := NoopKafkaProducer{}
	err := noopProducer.AddOffsetsToTxn(map[string][]*sarama.PartitionOffsetMetadata{}, "")
	assert.Nil(t, err)
}

func TestAddMessageToTxn(t *testing.T) {
	noopProducer := NoopKafkaProducer{}
	err := noopProducer.AddMessageToTxn(nil, "", nil)
	assert.Nil(t, err)
}
