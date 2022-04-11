package events

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/cloudtrust/common-service/v2/log"
)

type EventsReporterModule interface {
	ReportEvent(ctx context.Context, event Event)
}

type eventsReporterModule struct {
	producer sarama.SyncProducer
	topic    string
	logger   log.Logger
}

func NewEventReporterModule(producer sarama.SyncProducer, topic string, logger log.Logger) EventsReporterModule {
	return &eventsReporterModule{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}
}

func (e *eventsReporterModule) ReportEvent(ctx context.Context, event Event) {
	serializedEvent := event.serialize()
	msg := &sarama.ProducerMessage{Topic: e.topic, Value: sarama.StringEncoder(serializedEvent)}
	_, _, err := e.producer.SendMessage(msg)

	if err != nil {
		e.logger.Error(ctx, "msg", "failed to persist event in kafka", "err", err.Error(), "event", event)
	}
}
