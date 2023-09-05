package events

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/cloudtrust/common-service/v2/log"
)

// AuditEventsReporterModule interface
type AuditEventsReporterModule interface {
	ReportEvent(ctx context.Context, event Event)
}

type auditEventsReporterModule struct {
	producer sarama.SyncProducer
	topic    string
	logger   log.Logger
}

// NewAuditEventReporterModule creates an instance of AuditEventsReporterModule
func NewAuditEventReporterModule(producer sarama.SyncProducer, topic string, logger log.Logger) AuditEventsReporterModule {
	return &auditEventsReporterModule{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}
}

func (e *auditEventsReporterModule) ReportEvent(ctx context.Context, event Event) {
	serializedEvent := event.serialize()
	base64Event := base64.StdEncoding.EncodeToString(serializedEvent)

	key, ok := event.details[CtEventAgentUserID]
	if !ok {
		key = "DEFAULT-KEY"
	}

	msg := &sarama.ProducerMessage{Topic: e.topic, Key: sarama.StringEncoder(key), Value: sarama.StringEncoder(base64Event)}
	_, _, err := e.producer.SendMessage(msg)

	if err != nil {
		event.details[CtEventAuditTime] = event.time.String()
		event.details[CtEventOrigin] = event.origin
		event.details[CtEventType] = event.eventType
		eventJSON, errMarshal := json.Marshal(event.details)
		if errMarshal == nil {
			e.logger.Error(ctx, "msg", "failed to persist event in kafka", "err", err.Error(), "eventJSON", string(eventJSON), "key", key, "event64", base64Event)
		} else {
			e.logger.Error(ctx, "msg", "failed to persist event in kafka", "err", err.Error(), "key", key, "event64", base64Event)
		}
	}
}
