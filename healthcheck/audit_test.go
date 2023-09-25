package healthcheck

import (
	"testing"
	"time"

	"github.com/cloudtrust/common-service/v2/healthcheck/mock"
	log "github.com/cloudtrust/common-service/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuditEventsReporterChecker(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockAuditEventReporter = mock.NewAuditEventsReporterModule(mockCtrl)

	t.Run("Success ", func(t *testing.T) {
		var auditEventReporterChecker = newAuditEventsReporterChecker("alias", mockAuditEventReporter, 1*time.Second, 1*time.Second, log.NewNopLogger())
		// first call
		mockAuditEventReporter.EXPECT().ReportEvent(gomock.Any(), gomock.Any()).Times(1)

		var res = auditEventReporterChecker.CheckStatus()
		assert.NotNil(t, res.Connection)
		assert.Equal(t, "init", *res.Connection)

		time.Sleep(2 * time.Second)
		mockAuditEventReporter.EXPECT().ReportEvent(gomock.Any(), gomock.Any()).Times(1)

		// usual success call
		res = auditEventReporterChecker.CheckStatus()
		assert.NotNil(t, res.Connection)
		assert.Equal(t, "established", *res.Connection)

		time.Sleep(500 * time.Millisecond)

		// Get succes from the cache
		res = auditEventReporterChecker.CheckStatus()
		assert.NotNil(t, res.Connection)
		assert.Equal(t, "established", *res.Connection)
	})

	t.Run("Failure ", func(t *testing.T) {
		var auditEventReporterChecker = newAuditEventsReporterChecker("alias", mockAuditEventReporter, 1*time.Second, 10*time.Second, log.NewNopLogger())
		mockAuditEventReporter.EXPECT().ReportEvent(gomock.Any(), gomock.Any()).Do(func(arg0 interface{}, arg1 interface{}) {
			time.Sleep(2 * time.Second)
		})

		// first call
		var res = auditEventReporterChecker.CheckStatus()
		assert.NotNil(t, res.Connection)
		assert.Equal(t, "init", *res.Connection)

		time.Sleep(2 * time.Second)

		// Call with down status
		res = auditEventReporterChecker.CheckStatus()
		assert.NotNil(t, res.Message)
		assert.Equal(t, "Events reporter timeout", *res.Message)
	})
}
