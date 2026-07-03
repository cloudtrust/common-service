package healthcheck

import (
	"testing"
	"time"

	"github.com/cloudtrust/common-service/v2/healthcheck/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestKafkaConsumerHealthCheckLiveAndCached(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTime := mock.NewTimeProvider(mockCtrl)
	mockTime.EXPECT().Now().Return(testTime).AnyTimes()

	consumer := mock.NewKafkaConsumer(mockCtrl)
	consumer.EXPECT().IsLive().Return(true).Times(1)
	checker := newKafkaConsumerChecker("alias", consumer, 10*time.Second, mockTime)

	status := checker.CheckStatus()
	assert.Equal(t, "UP", *status.State)
	assert.Equal(t, "consuming", *status.Connection)
	assert.Nil(t, status.Message)

	status = checker.CheckStatus()
	assert.Equal(t, "UP", *status.State)
	assert.Equal(t, "consuming", *status.Connection)

	internal := checker.(*kafkaConsumerChecker)
	assert.Equal(t, 0, internal.failureCounter)
}

func TestKafkaConsumerHealthCheckDownAndCached(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTime := mock.NewTimeProvider(mockCtrl)
	mockTime.EXPECT().Now().Return(testTime).AnyTimes()

	consumer := mock.NewKafkaConsumer(mockCtrl)
	consumer.EXPECT().IsLive().Return(false).Times(1)
	checker := newKafkaConsumerChecker("alias", consumer, 10*time.Second, mockTime)

	status := checker.CheckStatus()
	assert.Equal(t, "DOWN", *status.State)
	assert.Nil(t, status.Connection)
	assert.NotNil(t, status.Message)
	assert.Equal(t, "Kafka consumer is not live", *status.Message)

	status = checker.CheckStatus()
	assert.Equal(t, "DOWN", *status.State)
	assert.Equal(t, "Kafka consumer is not live", *status.Message)

	internal := checker.(*kafkaConsumerChecker)
	assert.Equal(t, 1, internal.failureCounter)
}

func TestKafkaConsumerHealthCheckFailureCounterResetOnRecovery(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTime := mock.NewTimeProvider(mockCtrl)
	gomock.InOrder(
		mockTime.EXPECT().Now().Return(testTime),
		mockTime.EXPECT().Now().Return(testTime),
		mockTime.EXPECT().Now().Return(testTime.Add(10*time.Second).Add(time.Millisecond)),
		mockTime.EXPECT().Now().Return(testTime.Add(10*time.Second).Add(time.Millisecond)),
	)

	consumer := mock.NewKafkaConsumer(mockCtrl)
	gomock.InOrder(
		consumer.EXPECT().IsLive().Return(false),
		consumer.EXPECT().IsLive().Return(true),
	)
	checker := newKafkaConsumerChecker("alias", consumer, 10*time.Second, mockTime)

	status := checker.CheckStatus()
	assert.Equal(t, "DOWN", *status.State)
	assert.Equal(t, "Kafka consumer is not live", *status.Message)

	internal := checker.(*kafkaConsumerChecker)
	assert.Equal(t, 1, internal.failureCounter)

	status = checker.CheckStatus()
	assert.Equal(t, "UP", *status.State)
	assert.Equal(t, "consuming", *status.Connection)
	assert.Nil(t, status.Message)
	assert.Equal(t, 0, internal.failureCounter)
}
