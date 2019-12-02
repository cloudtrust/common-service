package security

//go:generate mockgen -destination=./mock/keycloak_client.go -package=mock -mock_names=KeycloakClient=KeycloakClient github.com/cloudtrust/common-service/security KeycloakClient
//go:generate mockgen -destination=./mock/authentication_db_reader.go -package=mock -mock_names=AuthorizationDBReader=AuthorizationDBReader github.com/cloudtrust/common-service/security AuthorizationDBReader

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudtrust/common-service/configuration"

	cs "github.com/cloudtrust/common-service"
	"github.com/cloudtrust/common-service/log"
	"github.com/cloudtrust/common-service/security/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCheckAuthorizationOnRealm(t *testing.T) {
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
			configuration.Authorization{
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
			configuration.Authorization{
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
			configuration.Authorization{
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
			configuration.Authorization{
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
			configuration.Authorization{
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
			configuration.Authorization{
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
			configuration.Authorization{
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

	// Authorized for all groups (test wildcard)
	{
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			configuration.Authorization{
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
	}

	// Test no groups assigned to targetUser
	{
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			configuration.Authorization{
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
	}

	// Test allowed only for non master realm
	{
		var targetRealm = "toto"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			configuration.Authorization{
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
	}

	// Authorized for all realms (test wildcard) and all groups
	{
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			configuration.Authorization{
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
	}

	// Test cannot GetUser infos
	{
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			configuration.Authorization{
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
	}

	// Test for a specific target group
	{
		var targetRealm = "toto"
		var targetUserID = "123-456-789"
		var targetGroupName = "customer"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			configuration.Authorization{
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
	}

	// Deny
	{
		var targetRealm = "toto"
		var targetUserID = "123-456-789"

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{
			configuration.Authorization{
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
	}
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
			configuration.Authorization{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &getUsers,
				TargetRealmID:   &master,
				TargetGroupName: &any,
			},
			configuration.Authorization{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &createUser,
				TargetRealmID:   &master,
				TargetGroupName: &integratorAgent,
			},
			configuration.Authorization{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &createUser,
				TargetRealmID:   &master,
				TargetGroupName: &integratorManager,
			},
			configuration.Authorization{
				RealmID:         &master,
				GroupName:       &toe,
				Action:          &createUser,
				TargetRealmID:   &master,
				TargetGroupName: &l2SupportAgent,
			},
			configuration.Authorization{
				RealmID:   &master,
				GroupName: &toe,
				Action:    &l2SupportAgent,
			},
			configuration.Authorization{
				RealmID:         &dep,
				GroupName:       &productAdministrator,
				Action:          &getUsers,
				TargetRealmID:   &dep,
				TargetGroupName: &any,
			},
			configuration.Authorization{
				RealmID:         &dep,
				GroupName:       &productAdministrator,
				Action:          &createUser,
				TargetRealmID:   &dep,
				TargetGroupName: &l1SupportManager,
			},
			configuration.Authorization{
				RealmID:         &dep,
				GroupName:       &l1SupportManager,
				Action:          &getUsers,
				TargetRealmID:   &dep,
				TargetGroupName: &l1SupportAgent,
			},
			configuration.Authorization{
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
			configuration.Authorization{
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

func TestAction(t *testing.T) {
	var action = Action{
		Id:    1,
		Name:  "test",
		Scope: ScopeGlobal,
	}

	assert.Equal(t, "test", action.String())
}
