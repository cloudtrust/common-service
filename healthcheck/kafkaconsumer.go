package healthcheck

import (
	"time"
)

type KafkaConsumer interface {
	IsLive() bool
}

type kafkaConsumerChecker struct {
	alias          string
	consumer       KafkaConsumer
	response       HealthStatus
	failureCounter int
}

func newKafkaConsumerChecker(alias string, consumer KafkaConsumer, cacheDuration time.Duration, timeProvider TimeProvider) BasicChecker {
	healthStatusType := "kafkaConsumer"
	response := HealthStatus{Name: &alias, Type: &healthStatusType, CacheDuration: cacheDuration, TimeProvider: timeProvider}
	response.connection("init")
	response.stateUp()
	return &kafkaConsumerChecker{
		alias:          alias,
		consumer:       consumer,
		response:       response,
		failureCounter: 0,
	}
}

func (k *kafkaConsumerChecker) CheckStatus() HealthStatus {
	if !k.response.hasExpired() {
		return k.response
	}

	if k.consumer.IsLive() {
		k.response.connection("consuming")
		k.response.stateUp()
		k.failureCounter = 0
	} else {
		k.response.stateDown("Kafka consumer is not live")
		k.failureCounter++
	}

	k.response.touch()
	return k.response
}
