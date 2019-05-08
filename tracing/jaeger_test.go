package tracing

//go:generate mockgen -destination=./mock/configuration.go -package=mock -mock_names=Configuration=Configuration github.com/cloudtrust/common-service Configuration

import (
	"context"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/metrics/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateNoopJaegerClient(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	var prefix = "prefix"
	var expected = "expected"

	mockConf.EXPECT().GetBool(prefix).Return(false).Times(1)

	var jaeger, err = CreateJaegerClient(mockConf, prefix, "")
	assert.Nil(t, err)

	jaeger.TryStartSpanWithTag(context.TODO(), "name", "tagName", "tagValue")
	jaeger.MakeHTTPTracingMW("component", "operation")
	jaeger.Close()

	var e = jaeger.MakeEndpointTracingMW("name")(
		func(_ context.Context, _ interface{}) (interface{}, error) {
			return expected, nil
		},
	)
	result, _ := e(context.TODO(), nil)
	assert.Equal(t, expected, result)
}

func TestCreateJaegerClientSucceeds(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	var prefix = "prefix"

	mockConf.EXPECT().GetBool(prefix).Return(true).Times(1)
	mockConf.EXPECT().GetString(prefix + "-sampler-type").Return("").Times(1)
	mockConf.EXPECT().GetFloat64(prefix + "-sampler-param").Return(0.0).Times(1)
	mockConf.EXPECT().GetString(prefix + "-sampler-host-port").Return("").Times(1)
	mockConf.EXPECT().GetBool(prefix + "-reporter-logspan").Return(false).Times(1)
	mockConf.EXPECT().GetDuration(prefix + "-write-interval").Return(time.Duration(0)).Times(1)

	var jaeger, err = CreateJaegerClient(mockConf, "prefix", "cloudtrusttester")
	defer jaeger.Close()

	assert.Nil(t, err)
}
