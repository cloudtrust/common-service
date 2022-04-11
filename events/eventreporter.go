package events

import (
	"context"
	"encoding/json"

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
		event.details[CtEventAuditTime] = event.time.String()
		event.details[CtEventOrigin] = event.origin
		event.details[CtEventType] = event.eventType
		eventJSON, errMarshal := json.Marshal(event.details)
		if errMarshal == nil {
			e.logger.Error(ctx, "msg", "failed to persist event in kafka", "err", err.Error(), "event", string(eventJSON))
		} else {
			e.logger.Error(ctx, "msg", "failed to persist event in kafka", "err", err.Error())
		}
	}
}
