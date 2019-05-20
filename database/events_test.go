package database

import (
	"context"
	"testing"

	cs "github.com/cloudtrust/common-service"
	"github.com/cloudtrust/common-service/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
