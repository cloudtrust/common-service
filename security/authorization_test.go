package security

//go:generate mockgen -destination=./mock/keycloak_client.go -package=mock -mock_names=KeycloakClient=KeycloakClient github.com/cloudtrust/common-service/security KeycloakClient

import (
	"context"
	"fmt"
	"testing"

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

	var accessToken = "TOKEN=="
	var groups = []string{"toe", "svc"}

	// Authorized for all realm (test wildcard)
	{
		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"GetRealm": {"*": {} }} }}`)
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, "master")

		err = authorizationManager.CheckAuthorizationOnTargetRealm(ctx, "GetRealm", "master")

		assert.Nil(t, err)
	}

	// Authorized for non admin realm
	{
		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"GetRealm": {"/": {} }} }}`)
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
		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"GetRealm": {"master": {} }} }}`)
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
		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"CreateUser": {"master": {} }} }}`)
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
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)

	var accessToken = "TOKEN=="
	var groups = []string{"toe", "svc"}
	var realm = "master"

	// Authorized for all groups (test wildcard)
	{
		var targetRealm = "master"
		var targetGroupID = "123-456-789"

		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"DeleteUser": {"master": { "*": {} } }} }}`)
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupName(accessToken, targetRealm, targetGroupID).Return(groupName, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetGroupID(ctx, "DeleteUser", "master", targetGroupID)
		assert.Nil(t, err)
	}

	// Error management
	{
		var targetRealm = "master"
		var targetGroupID = "123-456-789"

		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"DeleteUser": {"master": { "*": {} } }} }}`)
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		mockKeycloakClient.EXPECT().GetGroupName(accessToken, targetRealm, targetGroupID).Return("", fmt.Errorf("ERROR")).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetGroupID(ctx, "DeleteUser", "master", targetGroupID)
		assert.Equal(t, ForbiddenError{}, err)
	}

	// No group found with this ID
	{
		var targetRealm = "master"
		var targetGroupID = "123-456-789"

		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"DeleteUser": {"master": { "*": {} } }} }}`)
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		mockKeycloakClient.EXPECT().GetGroupName(accessToken, targetRealm, targetGroupID).Return("", nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetGroupID(ctx, "DeleteUser", "master", targetGroupID)
		assert.Equal(t, ForbiddenError{}, err)
	}

}

func TestCheckAuthorizationOnTargetUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)

	var accessToken = "TOKEN=="
	var groups = []string{"toe", "svc"}
	var realm = "master"

	// Authorized for all groups (test wildcard)
	{
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"DeleteUser": {"master": { "*": {} } }} }}`)
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var userID = "123-456-789"
		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(accessToken, targetRealm, targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", "master", userID)
		assert.Nil(t, err)
	}

	// Test no groups assigned to targetUser
	{
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"DeleteUser": {"master": { "*": {} } }} }}`)
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var userID = "123-456-789"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(accessToken, targetRealm, targetUserID).Return([]string{}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", "master", userID)
		assert.Equal(t, ForbiddenError{}, err)
	}

	// Test allowed only for non master realm
	{
		var targetRealm = "toto"
		var targetUserID = "123-456-789"

		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"DeleteUser": {"/": { "*": {} } }} }}`)
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(accessToken, targetRealm, targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", targetRealm, targetUserID)
		assert.Nil(t, err)

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(accessToken, "master", targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", "master", targetUserID)
		assert.Equal(t, ForbiddenError{}, err)
	}

	// Authorized for all realms (test wildcard) and all groups
	{
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"DeleteUser": {"*": { "*": {} } }} }}`)
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(accessToken, targetRealm, targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", "master", targetUserID)
		assert.Nil(t, err)
	}

	// Test cannot GetUser infos
	{
		var targetRealm = "master"
		var targetUserID = "123-456-789"

		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"DeleteUser": {"*": { "*": {} } }} }}`)
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(accessToken, targetRealm, targetUserID).Return([]string{}, fmt.Errorf("Error")).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", "master", targetUserID)
		assert.Equal(t, ForbiddenError{}, err)
	}

	// Test for a specific target group
	{
		var targetRealm = "toto"
		var targetUserID = "123-456-789"

		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"DeleteUser": {"toto": { "customer": {} } }} }}`)
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(accessToken, targetRealm, targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", targetRealm, targetUserID)
		assert.Nil(t, err)
	}

	// Deny
	{
		var targetRealm = "toto"
		var targetUserID = "123-456-789"

		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), `{"master": {"toe": {"DeleteUser": {}} }}`)
		assert.Nil(t, err)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realm)

		var groupName = "customer"

		mockKeycloakClient.EXPECT().GetGroupNamesOfUser(accessToken, targetRealm, targetUserID).Return([]string{
			groupName,
		}, nil).Times(1)

		err = authorizationManager.CheckAuthorizationOnTargetUser(ctx, "DeleteUser", targetRealm, targetUserID)
		assert.Equal(t, ForbiddenError{}, err)
	}
}

func TestLoadAuthorizationsFromMissingFile(t *testing.T) {
	_, err := NewAuthorizationManagerFromFile(nil, nil, "missing.file")
	assert.NotNil(t, err)
}

func TestLoadAuthorizations(t *testing.T) {
	// Empty file
	{
		var jsonAuthz = ""
		_, err := loadAuthorizations(jsonAuthz)
		assert.NotNil(t, err)
		assert.Equal(t, "JSON structure expected", err.Error())

		_, err = NewAuthorizationManager(nil, log.NewNopLogger(), jsonAuthz)
		assert.NotNil(t, err)
		assert.Equal(t, "JSON structure expected", err.Error())
	}

	// Empty JSON
	{
		var jsonAuthz = "{}"
		_, err := loadAuthorizations(jsonAuthz)
		assert.Nil(t, err)
	}

	// Wrong format
	{
		var jsonAuthz = "{sdf}ref"
		_, err := loadAuthorizations(jsonAuthz)
		assert.NotNil(t, err)
	}

	// Correct format
	{
		var jsonAuthz = `{
			"master":{
			  "toe_administrator":{
				"GetUsers": {
				  "master": {
					"*": {}
				  }
				},
				"CreateUser": {
				  "master": {
					"integrator_manager": {},
					"integrator_agent": {},
					"l2_support_manager": {},
					"l2_support_agent": {},
					"l3_support_manager": {},
					"l3_support_agent": {}
				  }
				}
			  },
			  "l3_support_agent": {}
			},
			"DEP":{
			  "product_administrator":{
				"GetUsers": {
				  "DEP": {
					"*": {}
				  }
				},
				"CreateUser": {
				  "DEP": {
					"l1_support_manager": {}
				  }
				}
			  },
			  "l1_support_manager": {
				"GetUsers": {
				  "DEP": {
					"l1_support_agent": {},
					"end_user": {}
				  }
				}
			  }
			}
		  }`

		authorizations, err := loadAuthorizations(jsonAuthz)
		assert.Nil(t, err)

		_, ok := authorizations["master"]["toe_administrator"]["GetUsers"]["master"]["*"]
		assert.Equal(t, true, ok)

		_, ok = authorizations["master"]["toe_administrator"]["GetUsers"]["master"]
		assert.Equal(t, true, ok)

		_, ok = authorizations["master"]["l3_support_agent"]
		assert.Equal(t, true, ok)

		_, ok = authorizations["master"]["l3_support_agent"]["GetUsers"]["master"]
		assert.Equal(t, false, ok)

		_, ok = authorizations["DEP"]["l1_support_manager"]["GetUsers"]["DEP"]
		assert.Equal(t, true, ok)

		_, ok = authorizations["DEP"]["l1_support_manager"]["GetUsers"]["DEP"]["end_user"]
		assert.Equal(t, true, ok)

		_, ok = authorizations["DEP"]["l1_support_manager"]["GetUsers"]["DEP"]["end_user2"]
		assert.Equal(t, false, ok)
	}
}

func TestGetAuthorization(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)

	{
		var jsonAuthz = `{
			"master":{
			  "toe_administrator":{
				"GetUsers": {
				  "master": {
					"*": {}
				  }
				},
				"CreateUser": {
				  "master": {
					"integrator_manager": {},
					"integrator_agent": {},
					"l2_support_manager": {},
					"l2_support_agent": {},
					"l3_support_manager": {},
					"l3_support_agent": {}
				  }
				}
			  },
			  "l3_support_agent": {}
			},
			"DEP":{
			  "product_administrator":{
				"GetUsers": {
				  "DEP": {
					"*": {}
				  }
				},
				"CreateUser": {
				  "DEP": {
					"l1_support_manager": {}
				  }
				}
			  },
			  "l1_support_manager": {
				"GetUsers": {
				  "DEP": {
					"l1_support_agent": {},
					"end_user": {}
				  }
				}
			  }
			}
		  }`

		var authorizationManager, err = NewAuthorizationManager(mockKeycloakClient, log.NewNopLogger(), jsonAuthz)
		assert.Nil(t, err)

		var ctx = context.Background()
		ctx = context.WithValue(ctx, cs.CtContextRealm, "master")
		ctx = context.WithValue(ctx, cs.CtContextGroups, []string{"toe_administrator"})
		ctx = context.WithValue(ctx, cs.CtContextUsername, "toe")

		rights := authorizationManager.GetRightsOfCurrentUser(ctx)

		_, ok := rights["toe_administrator"]["GetUsers"]["master"]["*"]
		assert.Equal(t, true, ok)

		_, ok = rights["toe_administrator"]["CreateUser"]["master"]
		assert.Equal(t, true, ok)

	}
}
