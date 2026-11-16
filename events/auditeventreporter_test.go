package events

import (
	"context"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/v2/events/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReportEvent(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	serviceName := "test_service"
	realmName := "TEST"
	eventType := "test_event"
	userID := "testerID"
	producer := mock.NewProducer(mockCtrl)
	eventReporter := NewAuditEventReporterModule(producer, serviceName)
	ctx := context.Background()
	event := Event{
		time:      time.Now().UTC(),
		origin:    "test",
		eventType: eventType,
		details: map[string]string{
			CtEventAgentUserID:    userID,
			CtEventAgentUsername:  "tester",
			CtEventAgentRealmName: realmName,
		},
	}

	producer.EXPECT().SendWithMap(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(bytes []byte, key string, details map[string]any) error {
		assert.NotNil(t, bytes)
		assert.Equal(t, userID, key)
		assert.Equal(t, serviceName, details["svc"])
		assert.Equal(t, realmName, details["realm"])
		assert.Equal(t, eventType, details["event_type"])
		return nil
	})
	eventReporter.ReportEvent(ctx, event)
}
