package database

//go:generate mockgen -destination=./mock/module.go -package=mock -mock_names=CloudtrustDB=CloudtrustDB github.com/cloudtrust/common-service/database CloudtrustDB

import (
	"context"
	"testing"

	"github.com/cloudtrust/common-service/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEventsDBModule(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockDB = mock.NewCloudtrustDB(mockCtrl)

	var eventsDBModule = NewEventsDBModule(mockDB)

	var ctx = context.WithValue(context.Background(), "username", "my name")
	ctx = context.WithValue(ctx, "realm", "myrealm")

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
