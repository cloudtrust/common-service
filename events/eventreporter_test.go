package events

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/v2/events/mock"
	"github.com/golang/mock/gomock"
)

func TestReportEvent(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	producer := mock.NewSyncProducer(mockCtrl)
	topic := "testTopic"
	logger := mock.NewLogger(mockCtrl)
	eventReporter := NewEventReporterModule(producer, topic, logger)
	ctx := context.Background()
	event := Event{
		time:      time.Now().UTC(),
		origin:    "test",
		eventType: "test_event",
		details: map[string]string{
			CtEventAgentUserID:    "testerID",
			CtEventAgentUsername:  "tester",
			CtEventAgentRealmName: "TEST",
		},
	}

	t.Run("SUCCESS", func(t *testing.T) {
		producer.EXPECT().SendMessage(gomock.Any()).Return(int32(0), int64(0), nil)

		eventReporter.ReportEvent(ctx, event)
	})

	t.Run("FAILURE", func(t *testing.T) {
		producer.EXPECT().SendMessage(gomock.Any()).Return(int32(0), int64(0), fmt.Errorf("Kafka failure"))
		logger.EXPECT().Error(gomock.Any(), "msg", "failed to persist event in kafka", "err", "Kafka failure", "event", gomock.Any())
		eventReporter.ReportEvent(ctx, event)
	})
}
