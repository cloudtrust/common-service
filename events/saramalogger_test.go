package events

import (
	"testing"

	"github.com/cloudtrust/common-service/v2/events/mock"
	"github.com/golang/mock/gomock"
)

func TestWrite(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := mock.NewLogger(mockCtrl)

	writer := cloudtrustLoggerWrapper{
		logger: logger,
	}
	test := "test"
	logger.EXPECT().Info(gomock.Any(), "saramaMsg", test)

	writer.Write([]byte(test))
}
