package events

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/v2/events/mock"
	"go.uber.org/mock/gomock"
)

func TestReportEvent(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	producer := mock.NewProducer(mockCtrl)
	logger := mock.NewLogger(mockCtrl)
	eventReporter := NewAuditEventReporterModule(producer, logger)
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
		producer.EXPECT().SendPartitionedMessageBytes(gomock.Any(), gomock.Any()).Return(nil)

		eventReporter.ReportEvent(ctx, event)
	})

	t.Run("FAILURE", func(t *testing.T) {
		producer.EXPECT().SendPartitionedMessageBytes(gomock.Any(), gomock.Any()).Return(fmt.Errorf("Kafka failure"))
		logger.EXPECT().Error(gomock.Any(), "msg", "failed to persist event in kafka", "err", "Kafka failure", "eventJSON", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		eventReporter.ReportEvent(ctx, event)
	})
}
