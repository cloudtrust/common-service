package events

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/cloudtrust/common-service/v2/log"
)

type Producer interface {
	SendPartitionedMessageBytes(partitionKey string, content []byte) error
}

// AuditEventsReporterModule interface
type AuditEventsReporterModule interface {
	ReportEvent(ctx context.Context, event Event)
}

type auditEventsReporterModule struct {
	producer Producer
	logger   log.Logger
}

// NewAuditEventReporterModule creates an instance of AuditEventsReporterModule
func NewAuditEventReporterModule(producer Producer, logger log.Logger) AuditEventsReporterModule {
	return &auditEventsReporterModule{
		producer: producer,
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

	err := e.producer.SendPartitionedMessageBytes(key, []byte(base64Event))

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
