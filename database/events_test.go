package database

import (
	"context"
	"strings"
	"testing"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/cloudtrust/common-service/v2/database/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateAdditionalInfo(t *testing.T) {
	var addInfo = CreateAdditionalInfo("a", "b", "c", "d", "z")
	assert.True(t, strings.Contains(addInfo, `"a":"b"`))
	assert.True(t, strings.Contains(addInfo, `"c":"d"`))
	assert.False(t, strings.Contains(addInfo, `"z"`))
}

func TestEventsDBModule(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockDB = mock.NewCloudtrustDB(mockCtrl)

	var eventsDBModule = NewEventsDBModule(mockDB)

	var ctx = context.WithValue(context.Background(), cs.CtContextUsername, "my name")
	ctx = context.WithValue(ctx, cs.CtContextRealm, "myrealm")

	// Missing ct_event_type
	{
		var err = eventsDBModule.ReportEvent(ctx, "", "back-office", "type", "val")
		assert.Nil(t, err)
	}

	// ct_event_type is present
	{
		mockDB.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
		var err = eventsDBModule.ReportEvent(ctx, "event_type", "back-office", "type", "val")
		assert.Nil(t, err)
	}
}
