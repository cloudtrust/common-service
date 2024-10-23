package events

import (
	"testing"

	"github.com/cloudtrust/common-service/v2/events/mock"
	"go.uber.org/mock/gomock"
)

func TestWrite(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := mock.NewLogger(mockCtrl)

	writer := cloudtrustLoggerWrapper{
		logger: logger,
	}
	test := "test"
	logger.EXPECT().Info(gomock.Any(), "msg", test, "tag", "sarama")

	writer.Write([]byte(test))
}
