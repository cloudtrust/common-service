package metrics

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/configuration.go -package=mock -mock_names=Configuration=Configuration github.com/cloudtrust/common-service/v2 Configuration

import (
	"context"
	"testing"

	"github.com/cloudtrust/common-service/v2/database/mock"
	"github.com/cloudtrust/common-service/v2/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNoopInfluxClient(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	var prefix = "noop"

	mockConf.EXPECT().GetBool(prefix).Return(false).Times(1)

	var noop, err = NewMetrics(mockConf, prefix, nil)
	assert.Nil(t, err)
	defer noop.Close()

	// Coverage
	counter := noop.NewCounter("name")
	counter.Add(1.0)
	counter.With("value")

	gauge := noop.NewGauge("name")
	gauge.Add(1.0)
	gauge.Set(1.0)
	gauge.With("value")

	histo := noop.NewHistogram("name")
	histo.With("value")
	histo.Observe(1.0)

	noop.Ping(1)
	noop.WriteLoop(nil)
	noop.Stats(context.TODO(), "name", map[string]string{}, map[string]any{})
	noop.Close()
}

func TestInvalidConfigurationInfluxClient(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	var prefix = "name"

	mockConf.EXPECT().GetBool(prefix).Return(true).Times(1)
	mockConf.EXPECT().GetString(prefix + "-host-port").Return("influx.io#%").Times(1)
	mockConf.EXPECT().GetString(prefix + "-username").Return("username").Times(1)
	mockConf.EXPECT().GetString(prefix + "-password").Return("password").Times(1)

	var _, err = NewMetrics(mockConf, prefix, log.NewNopLogger())
	assert.NotNil(t, err)
}

func TestTrueInfluxClient(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	var prefix = "name"

	mockConf.EXPECT().GetBool(prefix).Return(true).Times(1)
	mockConf.EXPECT().GetString(prefix + "-host-port").Return("influx.io").Times(1)
	for _, suffix := range []string{"-username", "-password", "-database", "-retention-policy", "-write-consistency"} {
		mockConf.EXPECT().GetString(prefix + suffix).Return("value" + suffix).Times(1)
	}
	mockConf.EXPECT().GetString(prefix + "-precision").Return("s").Times(1)

	var influx, err = NewMetrics(mockConf, "name", log.NewNopLogger())
	assert.Nil(t, err)
	assert.NotNil(t, influx)
	defer influx.Close()

	influx.NewCounter("name")
	influx.NewGauge("name")
	influx.NewHistogram("name")
	influx.Ping(1)

	// Stats fails
	{
		var tags = map[string]string{}
		var fields = map[string]any{}
		err := influx.Stats(context.TODO(), "name", tags, fields)
		assert.NotNil(t, err)
	}
}
