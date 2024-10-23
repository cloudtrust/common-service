package security

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/keycloak_client.go -package=mock -mock_names=KeycloakClient=KeycloakClient github.com/cloudtrust/common-service/v2/security KeycloakClient
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/authentication_db_reader.go -package=mock -mock_names=AuthorizationDBReader=AuthorizationDBReader github.com/cloudtrust/common-service/v2/security AuthorizationDBReader
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/detailederr.go -package=mock -mock_names=DetailedError=DetailedError github.com/cloudtrust/common-service/v2/errors DetailedError

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/cloudtrust/common-service/v2/configuration"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/cloudtrust/common-service/v2/log"
	"github.com/cloudtrust/common-service/v2/security/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCheckAuthorizationOnTargetRealm(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockAuthorizationDBReader = mock.NewAuthorizationDBReader(mockCtrl)

	var accessToken = "TOKEN=="
	var groups = []string{"toe", "svc"}
	var master = "master"
	var toe = "toe"
	var getRealm = "GetRealm"
	var createUser = "CreateUser"
	var any = "*"
	var anyNonMasterRealm = "/"

	// Authorized for all realm (test wildcard)
	{
		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:       &master,
				GroupName:     &toe,
				Action:        &getRealm,
				TargetRealmID: &any,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, "master")

		err = authorizationManager.CheckAuthorizationOnTargetRealm(ctx, "GetRealm", "master")

		assert.Nil(t, err)
	}

	// Authorized for non admin realm
	{
		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:       &master,
				GroupName:     &toe,
				Action:        &getRealm,
				TargetRealmID: &anyNonMasterRealm,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, "master")

		err = authorizationManager.CheckAuthorizationOnTargetRealm(ctx, "GetRealm", "toto")
		assert.Nil(t, err)

		err = authorizationManager.CheckAuthorizationOnTargetRealm(ctx, "GetRealm", "master")
		assert.NotNil(t, err)
		assert.Equal(t, "ForbiddenError: Operation not permitted", err.Error())

	}

	// Authorized for specific realm
	{
		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:       &master,
				GroupName:     &toe,
				Action:        &getRealm,
				TargetRealmID: &master,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, "master")

		err = authorizationManager.CheckAuthorizationOnTargetRealm(ctx, "GetRealm", "master")
		assert.Nil(t, err)

		err = authorizationManager.CheckAuthorizationOnTargetRealm(ctx, "GetRealm", "other")
		assert.Equal(t, "ForbiddenError: Operation not permitted", err.Error())
	}

	// Deny by default
	{
		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:       &master,
				GroupName:     &toe,
				Action:        &createUser,
				TargetRealmID: &master,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, "master")

		err = authorizationManager.CheckAuthorizationOnTargetRealm(ctx, "GetRealm", "master")
		assert.Equal(t, ForbiddenError{}, err)
	}
}

func TestCheckAuthorizationOnTargetGroupID(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockAuthorizationDBReader = mock.NewAuthorizationDBReader(mockCtrl)
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)

	var accessToken = "TOKEN=="
	var groups = []string{"toe", "svc"}
	var realm = "master"
	var master = "master"
	var toe = "toe"
	var deleteUser = "DeleteUser"
	var any = "*"

	// Authorized for all groups (test wildcard)
	{
		var targetRealm = "master"
		var targetGroupID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &deleteUser,
				TargetRealmID:   &master,
				TargetGroupName: &any,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupName(ctx, accessToken, targetRealm, targetGroupID).Return(groupName, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetGroupID(ctx, "DeleteUser", "master", targetGroupID)
		assert.Nil(t, err)
	}

	// Error management
	{
		var targetRealm = "master"
		var targetGroupID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &deleteUser,
				TargetRealmID:   &master,
				TargetGroupName: &any,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		mockKeycloakClient.EXPECT().GetGroupName(ctx, accessToken, targetRealm, targetGroupID).Return("", errors.New("ERROR")).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetGroupID(ctx, "DeleteUser", "master", targetGroupID)
		assert.Equal(t, ForbiddenError{}, err)
	}

	// No group found with this ID
	{
		var targetRealm = "master"
		var targetGroupID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &deleteUser,
				TargetRealmID:   &master,
				TargetGroupName: &any,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		mockKeycloakClient.EXPECT().GetGroupName(ctx, accessToken, targetRealm, targetGroupID).Return("", nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetGroupID(ctx, "DeleteUser", "master", targetGroupID)
		assert.Equal(t, ForbiddenError{}, err)
	}

}

func TestCheckAuthorizationOnTargetUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockAuthorizationDBReader = mock.NewAuthorizationDBReader(mockCtrl)

	var accessToken = "TOKEN=="
	var groups = []string{"toe", "svc"}
	var realm = "master"
	var master = "master"
	var toe = "toe"
	var deleteUser = "DeleteUser"
	var any = "*"
	var anyNonMasterRealm = "/"

	t.Run("Authorized for all groups (test wildcard)", func(t *testing.T) {
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &deleteUser,
				TargetRealmID:   &master,
				TargetGroupName: &any,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var userID = "123-456-789"
		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(ctx, accessToken, targetRealm, targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", "master", userID)
		assert.Nil(t, err)
	})

	t.Run("Test no groups assigned to targetUser", func(t *testing.T) {
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &deleteUser,
				TargetRealmID:   &master,
				TargetGroupName: &any,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var userID = "123-456-789"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(ctx, accessToken, targetRealm, targetUserID).Return([]string{}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", "master", userID)
		assert.Equal(t, ForbiddenError{}, err)
	})

	t.Run("Test allowed only for non master realm", func(t *testing.T) {
		var targetRealm = "toto"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &deleteUser,
				TargetRealmID:   &anyNonMasterRealm,
				TargetGroupName: &any,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(ctx, accessToken, targetRealm, targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", targetRealm, targetUserID)
		assert.Nil(t, err)

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(ctx, accessToken, "master", targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", "master", targetUserID)
		assert.Equal(t, ForbiddenError{}, err)
	})

	t.Run("Authorized for all realms (test wildcard) and all groups", func(t *testing.T) {
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &deleteUser,
				TargetRealmID:   &any,
				TargetGroupName: &any,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(ctx, accessToken, targetRealm, targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", "master", targetUserID)
		assert.Nil(t, err)
	})

	t.Run("Test cannot GetUser infos", func(t *testing.T) {
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &deleteUser,
				TargetRealmID:   &any,
				TargetGroupName: &any,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(ctx, accessToken, targetRealm, targetUserID).Return([]string{}, errors.New("Error")).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", "master", targetUserID)
		assert.Equal(t, ForbiddenError{}, err)
	})

	t.Run("Test for a specific target group", func(t *testing.T) {
		var targetRealm = "toto"
		var targetUserID = "123-456-789"
		var targetGroupName = "customer"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &deleteUser,
				TargetRealmID:   &targetRealm,
				TargetGroupName: &targetGroupName,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(ctx, accessToken, targetRealm, targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", targetRealm, targetUserID)
		assert.Nil(t, err)
	})

	t.Run("Deny", func(t *testing.T) {
		var targetRealm = "toto"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:   &master,
				GroupName: &toe,
				Action:    &deleteUser,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(ctx, accessToken, targetRealm, targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", targetRealm, targetUserID)
		assert.Equal(t, ForbiddenError{}, err)
	})

	t.Run("SelfUser-Deny", func(t *testing.T) {
		var targetRealm = "toto"
		var targetUserID = "123-456-789"
		var targetGroupName = "customer"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:         &targetRealm,
				GroupName:       &targetGroupName,
				Action:          &deleteUser,
				TargetRealmID:   &targetRealm,
				TargetGroupName: &targetGroupName,
			},
		}, nil)
		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)
		ctx = context.WithValue(ctx, cs.CtContextUserID, targetUserID)

		err = authorizationManager.CheckAuthorizationOnSelfUser(ctx, "DeleteUser")
		assert.Equal(t, ForbiddenError{}, err)
	})
}

func TestGetAuthorization(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockAuthorizationDBReader = mock.NewAuthorizationDBReader(mockCtrl)

	var master = "master"
	var dep = "DEP"
	var toe = "toe"
	var getUsers = "GetUsers"
	var createUser = "CreateUser"
	var any = "*"
	var integratorAgent = "integrator_agent"
	var integratorManager = "integrator_manager"
	var l2SupportAgent = "l2_support_agent"
	var l1SupportManager = "l1_support_manager"
	var l1SupportAgent = "l1_support_agent"
	var productAdministrator = "product_administrator"
	var endUser = "end_user"

	{
		var authorizations = []configuration.Authorization{
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &getUsers,
				TargetRealmID:   &master,
				TargetGroupName: &any,
			},
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &createUser,
				TargetRealmID:   &master,
				TargetGroupName: &integratorAgent,
			},
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &createUser,
				TargetRealmID:   &master,
				TargetGroupName: &integratorManager,
			},
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &createUser,
				TargetRealmID:   &master,
				TargetGroupName: &l2SupportAgent,
			},
			{
				RealmID:   &master,
				GroupName: &toe,
				Action:    &l2SupportAgent,
			},
			{
				RealmID:         &dep,
				GroupName:       &productAdministrator,
				Action:          &getUsers,
				TargetRealmID:   &dep,
				TargetGroupName: &any,
			},
			{
				RealmID:         &dep,
				GroupName:       &productAdministrator,
				Action:          &createUser,
				TargetRealmID:   &dep,
				TargetGroupName: &l1SupportManager,
			},
			{
				RealmID:         &dep,
				GroupName:       &l1SupportManager,
				Action:          &getUsers,
				TargetRealmID:   &dep,
				TargetGroupName: &l1SupportAgent,
			},
			{
				RealmID:         &dep,
				GroupName:       &l1SupportManager,
				Action:          &getUsers,
				TargetRealmID:   &dep,
				TargetGroupName: &endUser,
			},
		}

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return(authorizations, nil)

		var authorizationManager, err = NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		var ctx = context.Background()
		ctx = context.WithValue(ctx, cs.CtContextRealm, "master")
		ctx = context.WithValue(ctx, cs.CtContextGroups, []string{"toe"})
		ctx = context.WithValue(ctx, cs.CtContextUsername, "toe_user")

		rights := authorizationManager.GetRightsOfCurrentUser(ctx)

		_, ok := rights["toe"]["GetUsers"]["master"]["*"]
		assert.Equal(t, true, ok)

		_, ok = rights["toe"]["CreateUser"]["master"]
		assert.Equal(t, true, ok)

	}
}

func TestSuggestForbiddenError(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockDetailedError = mock.NewDetailedError(mockCtrl)

	assert.Equal(t, ForbiddenError{}, suggestForbiddenError(errors.New("any error")))

	mockDetailedError.EXPECT().Status().Return(http.StatusBadRequest)
	assert.Equal(t, ForbiddenError{}, suggestForbiddenError(mockDetailedError))

	mockDetailedError.EXPECT().Status().Return(http.StatusUnauthorized)
	assert.Equal(t, mockDetailedError, suggestForbiddenError(mockDetailedError))
}

func TestReloadAuthorizations(t *testing.T) {

	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockAuthorizationDBReader = mock.NewAuthorizationDBReader(mockCtrl)

	var master = "master"
	var toe = "toe"
	var any = "*"
	var getUsers = "GetUsers"

	var ctx = context.Background()
	ctx = context.WithValue(ctx, cs.CtContextRealm, "master")
	ctx = context.WithValue(ctx, cs.CtContextGroups, []string{"toe"})
	ctx = context.WithValue(ctx, cs.CtContextUsername, "toe_user")

	{
		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, errors.New("Error"))

		_, err := NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.NotNil(t, err)
	}

	{
		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, nil)

		authorizationManager, err := NewAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, errors.New("Error"))
		err = authorizationManager.ReloadAuthorizations(ctx)
		assert.NotNil(t, err)
		rights := authorizationManager.GetRightsOfCurrentUser(ctx)
		_, ok := rights["toe"]["GetUsers"]["master"]["*"]
		assert.Equal(t, false, ok)

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &getUsers,
				TargetRealmID:   &master,
				TargetGroupName: &any,
			},
		}, nil)
		err = authorizationManager.ReloadAuthorizations(ctx)
		assert.Nil(t, err)

		rights = authorizationManager.GetRightsOfCurrentUser(ctx)
		_, ok = rights["toe"]["GetUsers"]["master"]["*"]
		assert.Equal(t, true, ok)
	}
}
