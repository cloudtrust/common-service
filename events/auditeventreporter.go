package events

import (
	"context"
	"encoding/base64"
)

// Producer defines an interface to send ProducerMessage
type Producer interface {
	SendWithMap(content []byte, partitionKey string, details map[string]any) error
}

// AuditEventsReporterModule interface
type AuditEventsReporterModule interface {
	ReportEvent(ctx context.Context, event Event)
}

type auditEventsReporterModule struct {
	producer    Producer
	serviceName string
}

// NewAuditEventReporterModule creates an instance of AuditEventsReporterModule
func NewAuditEventReporterModule(producer Producer, serviceName string) AuditEventsReporterModule {
	return &auditEventsReporterModule{
		producer:    producer,
		serviceName: serviceName,
	}
}

func (e *auditEventsReporterModule) ReportEvent(ctx context.Context, event Event) {
	serializedEvent := event.serialize()
	base64Event := base64.StdEncoding.EncodeToString(serializedEvent)

	key, ok := event.details[CtEventAgentUserID]
	if !ok {
		key = "DEFAULT-KEY"
	}
	realm, ok := event.details[CtEventAgentRealmName]
	if !ok {
		realm = "UNDEFINED"
	}

	// Does not take care of any error... the publisher function is assumed to handle it internally.
	_ = e.producer.SendWithMap([]byte(base64Event), key, map[string]any{
		"svc":        e.serviceName,
		"realm":      realm,
		"time":       event.time.String(),
		"content":    base64Event,
		"origin":     event.origin,
		"event_type": event.eventType,
	})
}
