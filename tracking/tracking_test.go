package tracking

//generate mockgen --build_flags=--mod=mod -destination=./mock/configuration.go -package=mock -mock_names=Configuration=Configuration github.com/cloudtrust/common-service Configuration

import (
	"errors"
	"testing"

	"github.com/cloudtrust/common-service/metrics/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNoopSentry(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	mockConf.EXPECT().GetBool("sentry").Return(false).Times(1)

	var sentry, _ = NewSentry(mockConf, "sentry")
	defer sentry.Close()

	// CaptureError
	assert.Zero(t, sentry.CaptureError(nil, nil))
	assert.Zero(t, sentry.CaptureError(errors.New("fail"), map[string]string{"key": "val"}))

	// URL
	assert.Zero(t, sentry.URL())
}

func TestTrueSentry(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	mockConf.EXPECT().GetBool("sentry").Return(true).Times(1)
	mockConf.EXPECT().GetString("sentry-dsn").Return("dsn").Times(1)

	var sentry, _ = NewSentry(mockConf, "sentry")
	defer sentry.Close()

	assert.NotNil(t, sentry)
}
