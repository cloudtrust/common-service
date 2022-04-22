package events

import (
	"context"
	"time"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/cloudtrust/common-service/v2/events/fb"
	"github.com/cloudtrust/common-service/v2/log"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/google/uuid"
)

const (
	CtEventType            = "ct_event_type"
	CtEventAgentUsername   = "agent_username"
	CtEventAgentRealmName  = "agent_realm_name"
	CtEventTargetUserID    = "target_user_id"
	CtEventGroupID         = "group_id"
	CtEventGroupName       = "group_name"
	CtEventRoleID          = "role_id"
	CtEventRoleName        = "role_name"
	CtEventOrigin          = "origin"
	CtEventAuditTime       = "audit_time"
	CtEventTargetRealmName = "target_realm_name"
	CtEventAgentUserID     = "agent_user_id"
	CtEventTargetUsername  = "target_username"
	CtEventKcEventType     = "kc_event_type"
	CtEventKcOperationType = "kc_operation_type"
	CtEventClientID        = "client_id"
	CtEventAdditionalInfo  = "additional_info"

	CtEventUnknownUsername = "--UNKNOWN--"
)

type Event struct {
	time      time.Time
	origin    string
	eventType string
	details   map[string]string
}

func newEvent(origin string, eventType string, agentRealmName string, agentUserID string, agentUsername string, targetRealmName string, targetUserID *string, targetUserName *string, details map[string]string) Event {
	if details == nil {
		details = map[string]string{}
	}

	details[CtEventAgentRealmName] = agentRealmName
	details[CtEventAgentUserID] = agentUserID
	details[CtEventAgentUsername] = agentUsername
	details[CtEventTargetRealmName] = targetRealmName
	if targetUserID != nil {
		details[CtEventTargetUserID] = *targetUserID
	}
	if targetUserName != nil {
		details[CtEventTargetUsername] = *targetUserName
	}

	return Event{
		time:      time.Now().UTC(),
		origin:    origin,
		eventType: eventType,
		details:   details,
	}
}

func NewEventOnUser(origin string, eventType string, agentRealmName string, agentUserID string, agentUsername string, targetRealmName string, targetUserID string, targetUserName string, details map[string]string) Event {
	return newEvent(origin, eventType, agentRealmName, agentUserID, agentUsername, targetRealmName, &targetUserID, &targetUserName, details)
}

func NewEventOnUserFromContext(ctx context.Context, logger log.Logger, origin string, eventType string, targetRealmName string, targetUserID string, targetUserName string, details map[string]string) Event {
	agentRealmName := extractAgentValueFromContext(ctx, logger, cs.CtContextRealm)
	agentUserID := extractAgentValueFromContext(ctx, logger, cs.CtContextUserID)
	agentUserName := extractAgentValueFromContext(ctx, logger, cs.CtContextUsername)

	return NewEventOnUser(origin, eventType, agentRealmName, agentUserID, agentUserName, targetRealmName, targetUserID, targetUserName, details)
}

func NewEvent(origin string, eventType string, agentRealmName string, agentUserID string, agentUsername string, targetRealmName string, details map[string]string) Event {
	return newEvent(origin, eventType, agentRealmName, agentUserID, agentUsername, targetRealmName, nil, nil, details)
}

func NewEventFromContext(ctx context.Context, logger log.Logger, origin string, eventType string, targetRealmName string, details map[string]string) Event {
	agentRealmName := extractAgentValueFromContext(ctx, logger, cs.CtContextRealm)
	agentUserID := extractAgentValueFromContext(ctx, logger, cs.CtContextUserID)
	agentUserName := extractAgentValueFromContext(ctx, logger, cs.CtContextUsername)

	return NewEvent(origin, eventType, agentRealmName, agentUserID, agentUserName, targetRealmName, details)
}

func extractAgentValueFromContext(ctx context.Context, logger log.Logger, key cs.CtContext) string {
	value := ctx.Value(key)
	if value == nil {
		logger.Warn(ctx, "msg", "failed to extract agent information from context", "key", key)
		return ""
	}
	return value.(string)
}

func (e *Event) serialize() []byte {
	builder := flatbuffers.NewBuilder(1024)

	origin := builder.CreateString(e.origin)
	eventType := builder.CreateString(e.eventType)

	var tuples []flatbuffers.UOffsetT
	for k, v := range e.details {
		tuples = append(tuples, createTuple(builder, k, v))
	}
	// Add a unique id (UUID) for cloudtrust events
	tuples = append(tuples, createTuple(builder, "uid", uuid.New().String()))

	fb.CloudtrustEventStartDetailsVector(builder, len(tuples))
	for _, offset := range tuples {
		builder.PrependUOffsetT(offset)
	}
	details := builder.EndVector(len(tuples))

	fb.CloudtrustEventStart(builder)

	fb.CloudtrustEventAddTime(builder, e.time.UTC().UnixMilli())
	fb.CloudtrustEventAddOrigin(builder, origin)
	fb.CloudtrustEventAddCtEventType(builder, eventType)
	fb.CloudtrustEventAddDetails(builder, details)

	eventOffset := fb.CloudtrustEventEnd(builder)
	builder.Finish(eventOffset)
	return builder.FinishedBytes()
}

func createTuple(builder *flatbuffers.Builder, k string, v string) flatbuffers.UOffsetT {
	keyOffset := builder.CreateString(k)
	valueOffset := builder.CreateString(v)
	fb.TupleStart(builder)
	fb.TupleAddKey(builder, keyOffset)
	fb.TupleAddValue(builder, valueOffset)
	return fb.TupleEnd(builder)
}
