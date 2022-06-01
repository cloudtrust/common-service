package events

import (
	"testing"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/cloudtrust/common-service/v2/events/fb"
	"github.com/cloudtrust/common-service/v2/events/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestNewEventOnUser(t *testing.T) {
	origin := "test"
	eventType := "test_event"
	agentRealmName := "TEST1"
	agentUserID := "testerID"
	agentUsername := "tester"
	targetRealmName := "TEST2"
	targetUserID := "targetID"
	targetUsername := "target"
	details := map[string]string{
		"my_custom_field": "my_custom_value",
	}

	t.Run("SUCCESS", func(t *testing.T) {
		event := NewEventOnUser(origin, eventType, agentRealmName, agentUserID, agentUsername, targetRealmName, targetUserID, targetUsername, details)

		assert.NotNil(t, event.time)
		assert.Equal(t, origin, event.origin)
		assert.Equal(t, eventType, event.eventType)
		assert.Equal(t, agentRealmName, event.details[CtEventAgentRealmName])
		assert.Equal(t, agentUserID, event.details[CtEventAgentUserID])
		assert.Equal(t, agentUsername, event.details[CtEventAgentUsername])
		assert.Equal(t, targetRealmName, event.details[CtEventTargetRealmName])
		assert.Equal(t, targetUserID, event.details[CtEventTargetUserID])
		assert.Equal(t, targetUsername, event.details[CtEventTargetUsername])
		assert.Equal(t, details["my_custom_field"], event.details["my_custom_field"])
	})

	t.Run("SUCCESS with nil details", func(t *testing.T) {
		event := NewEventOnUser(origin, eventType, agentRealmName, agentUserID, agentUsername, targetRealmName, targetUserID, targetUsername, nil)

		assert.NotNil(t, event.time)
		assert.Equal(t, origin, event.origin)
		assert.Equal(t, eventType, event.eventType)
		assert.Equal(t, agentRealmName, event.details[CtEventAgentRealmName])
		assert.Equal(t, agentUserID, event.details[CtEventAgentUserID])
		assert.Equal(t, agentUsername, event.details[CtEventAgentUsername])
		assert.Equal(t, targetRealmName, event.details[CtEventTargetRealmName])
		assert.Equal(t, targetUserID, event.details[CtEventTargetUserID])
		assert.Equal(t, targetUsername, event.details[CtEventTargetUsername])
	})
}

func TestNewEventOnUserFromContext(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := mock.NewLogger(mockCtrl)

	origin := "test"
	eventType := "test_event"
	agentRealmName := "TEST1"
	agentUserID := "testerID"
	agentUsername := "tester"
	targetRealmName := "TEST2"
	targetUserID := "targetID"
	targetUsername := "target"
	details := map[string]string{
		"my_custom_field": "my_custom_value",
	}

	t.Run("SUCCESS", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, cs.CtContextRealm, agentRealmName)
		ctx = context.WithValue(ctx, cs.CtContextUserID, agentUserID)
		ctx = context.WithValue(ctx, cs.CtContextUsername, agentUsername)

		event := NewEventOnUserFromContext(ctx, logger, origin, eventType, targetRealmName, targetUserID, targetUsername, details)

		assert.NotNil(t, event.time)
		assert.Equal(t, origin, event.origin)
		assert.Equal(t, eventType, event.eventType)
		assert.Equal(t, agentRealmName, event.details[CtEventAgentRealmName])
		assert.Equal(t, agentUserID, event.details[CtEventAgentUserID])
		assert.Equal(t, agentUsername, event.details[CtEventAgentUsername])
		assert.Equal(t, targetRealmName, event.details[CtEventTargetRealmName])
		assert.Equal(t, targetUserID, event.details[CtEventTargetUserID])
		assert.Equal(t, targetUsername, event.details[CtEventTargetUsername])
		assert.Equal(t, details["my_custom_field"], event.details["my_custom_field"])
	})

	t.Run("Agent not complete", func(t *testing.T) {
		logger.EXPECT().Warn(gomock.Any(), "msg", "failed to extract agent information from context", "key", cs.CtContextRealm)
		ctx := context.Background()
		ctx = context.WithValue(ctx, cs.CtContextUserID, agentUserID)
		ctx = context.WithValue(ctx, cs.CtContextUsername, agentUsername)

		event := NewEventOnUserFromContext(ctx, logger, origin, eventType, targetRealmName, targetUserID, targetUsername, details)

		assert.NotNil(t, event.time)
		assert.Equal(t, origin, event.origin)
		assert.Equal(t, eventType, event.eventType)
		assert.Equal(t, "", event.details[CtEventAgentRealmName])
		assert.Equal(t, agentUserID, event.details[CtEventAgentUserID])
		assert.Equal(t, agentUsername, event.details[CtEventAgentUsername])
		assert.Equal(t, targetRealmName, event.details[CtEventTargetRealmName])
		assert.Equal(t, targetUserID, event.details[CtEventTargetUserID])
		assert.Equal(t, targetUsername, event.details[CtEventTargetUsername])
		assert.Equal(t, details["my_custom_field"], event.details["my_custom_field"])
	})
	t.Run("SUCCESS with nil details", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, cs.CtContextRealm, agentRealmName)
		ctx = context.WithValue(ctx, cs.CtContextUserID, agentUserID)
		ctx = context.WithValue(ctx, cs.CtContextUsername, agentUsername)
		event := NewEventOnUserFromContext(ctx, logger, origin, eventType, targetRealmName, targetUserID, targetUsername, nil)

		assert.NotNil(t, event.time)
		assert.Equal(t, origin, event.origin)
		assert.Equal(t, eventType, event.eventType)
		assert.Equal(t, agentRealmName, event.details[CtEventAgentRealmName])
		assert.Equal(t, agentUserID, event.details[CtEventAgentUserID])
		assert.Equal(t, agentUsername, event.details[CtEventAgentUsername])
		assert.Equal(t, targetRealmName, event.details[CtEventTargetRealmName])
		assert.Equal(t, targetUserID, event.details[CtEventTargetUserID])
		assert.Equal(t, targetUsername, event.details[CtEventTargetUsername])
	})
}

func TestNewEvent(t *testing.T) {
	origin := "test"
	eventType := "test_event"
	agentRealmName := "TEST1"
	agentUserID := "testerID"
	agentUsername := "tester"
	targetRealmName := "TEST2"
	details := map[string]string{
		"my_custom_field": "my_custom_value",
	}

	t.Run("SUCCESS", func(t *testing.T) {
		event := NewEvent(origin, eventType, agentRealmName, agentUserID, agentUsername, targetRealmName, details)

		assert.NotNil(t, event.time)
		assert.Equal(t, origin, event.origin)
		assert.Equal(t, eventType, event.eventType)
		assert.Equal(t, agentRealmName, event.details[CtEventAgentRealmName])
		assert.Equal(t, agentUserID, event.details[CtEventAgentUserID])
		assert.Equal(t, agentUsername, event.details[CtEventAgentUsername])
		assert.Equal(t, targetRealmName, event.details[CtEventTargetRealmName])
		assert.Equal(t, details["my_custom_field"], event.details["my_custom_field"])
	})

	t.Run("SUCCESS with nil details", func(t *testing.T) {
		event := NewEvent(origin, eventType, agentRealmName, agentUserID, agentUsername, targetRealmName, nil)

		assert.NotNil(t, event.time)
		assert.Equal(t, origin, event.origin)
		assert.Equal(t, eventType, event.eventType)
		assert.Equal(t, agentRealmName, event.details[CtEventAgentRealmName])
		assert.Equal(t, agentUserID, event.details[CtEventAgentUserID])
		assert.Equal(t, agentUsername, event.details[CtEventAgentUsername])
		assert.Equal(t, targetRealmName, event.details[CtEventTargetRealmName])
	})
}

func TestNewEventFromContext(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := mock.NewLogger(mockCtrl)

	origin := "test"
	eventType := "test_event"
	agentRealmName := "TEST1"
	agentUserID := "testerID"
	agentUsername := "tester"
	targetRealmName := "TEST2"
	details := map[string]string{
		"my_custom_field": "my_custom_value",
	}

	t.Run("SUCCESS", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, cs.CtContextRealm, agentRealmName)
		ctx = context.WithValue(ctx, cs.CtContextUserID, agentUserID)
		ctx = context.WithValue(ctx, cs.CtContextUsername, agentUsername)
		event := NewEventFromContext(ctx, logger, origin, eventType, targetRealmName, details)

		assert.NotNil(t, event.time)
		assert.Equal(t, origin, event.origin)
		assert.Equal(t, eventType, event.eventType)
		assert.Equal(t, agentRealmName, event.details[CtEventAgentRealmName])
		assert.Equal(t, agentUserID, event.details[CtEventAgentUserID])
		assert.Equal(t, agentUsername, event.details[CtEventAgentUsername])
		assert.Equal(t, targetRealmName, event.details[CtEventTargetRealmName])
		assert.Equal(t, details["my_custom_field"], event.details["my_custom_field"])
	})

	t.Run("Agent not complete", func(t *testing.T) {
		logger.EXPECT().Warn(gomock.Any(), "msg", "failed to extract agent information from context", "key", cs.CtContextUsername)
		ctx := context.Background()
		ctx = context.WithValue(ctx, cs.CtContextRealm, agentRealmName)
		ctx = context.WithValue(ctx, cs.CtContextUserID, agentUserID)
		event := NewEventFromContext(ctx, logger, origin, eventType, targetRealmName, details)

		assert.NotNil(t, event.time)
		assert.Equal(t, origin, event.origin)
		assert.Equal(t, eventType, event.eventType)
		assert.Equal(t, agentRealmName, event.details[CtEventAgentRealmName])
		assert.Equal(t, agentUserID, event.details[CtEventAgentUserID])
		assert.Equal(t, "", event.details[CtEventAgentUsername])
		assert.Equal(t, targetRealmName, event.details[CtEventTargetRealmName])
		assert.Equal(t, details["my_custom_field"], event.details["my_custom_field"])
	})
	t.Run("SUCCESS with nil details", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, cs.CtContextRealm, agentRealmName)
		ctx = context.WithValue(ctx, cs.CtContextUserID, agentUserID)
		ctx = context.WithValue(ctx, cs.CtContextUsername, agentUsername)

		event := NewEventFromContext(ctx, logger, origin, eventType, targetRealmName, nil)

		assert.NotNil(t, event.time)
		assert.Equal(t, origin, event.origin)
		assert.Equal(t, eventType, event.eventType)
		assert.Equal(t, agentRealmName, event.details[CtEventAgentRealmName])
		assert.Equal(t, agentUserID, event.details[CtEventAgentUserID])
		assert.Equal(t, agentUsername, event.details[CtEventAgentUsername])
		assert.Equal(t, targetRealmName, event.details[CtEventTargetRealmName])
	})
}

func TestSerialize(t *testing.T) {
	ctEventColumns := map[string]bool{
		CtEventType: true, CtEventAgentUsername: true, CtEventAgentRealmName: true, CtEventTargetUserID: true, CtEventOrigin: true, CtEventAuditTime: true, CtEventTargetRealmName: true,
		CtEventAgentUserID: true, CtEventTargetUsername: true, CtEventKcEventType: true, CtEventKcOperationType: true, CtEventClientID: true, CtEventAdditionalInfo: true}

	origin := "test"
	eventType := "test_event"
	agentRealmName := "TEST1"
	agentUserID := "testerID"
	agentUsername := "tester"
	targetRealmName := "TEST2"
	details := map[string]string{
		"my_custom_field": "my_custom_value",
	}

	e := NewEvent(origin, eventType, agentRealmName, agentUserID, agentUsername, targetRealmName, details)
	serializedEvent := e.serialize()

	event := fb.GetRootAsCloudtrustEvent(serializedEvent, 0)
	deserializedDetails := map[string]string{}
	for i := 0; i < event.DetailsLength(); i++ {
		var tuple = new(fb.Tuple)
		event.Details(tuple, i)
		key := string(tuple.Key())
		if _, ok := ctEventColumns[key]; ok {
			deserializedDetails[key] = string(tuple.Value())
		} else {
			deserializedDetails[key] = string(tuple.Value())
		}
	}

	assert.NotEqual(t, int64(0), event.Time())
	assert.Equal(t, origin, string(event.Origin()))
	assert.Equal(t, eventType, string(event.CtEventType()))
	assert.Equal(t, agentRealmName, deserializedDetails[CtEventAgentRealmName])
	assert.Equal(t, agentUserID, deserializedDetails[CtEventAgentUserID])
	assert.Equal(t, agentUsername, deserializedDetails[CtEventAgentUsername])
	assert.Equal(t, targetRealmName, deserializedDetails[CtEventTargetRealmName])
	assert.Equal(t, details["my_custom_field"], deserializedDetails["my_custom_field"])
}
