package security

import (
	"context"
	"errors"
	"testing"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/cloudtrust/common-service/v2/configuration"
	"github.com/cloudtrust/common-service/v2/log"
	"github.com/cloudtrust/common-service/v2/security/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCheckRoleAuthorizationOnTargetUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockKeycloakClient := mock.NewKeycloakClient(mockCtrl)
	mockAuthorizationDBReader := mock.NewAuthorizationDBReader(mockCtrl)

	accessToken := "TOKEN=="
	master := "master"
	groups := []string{"toe", "svc"}
	targetRealm := "targetRealm"
	targetUserID := "user-id-123"
	kycAction := KYCGetUser.String()
	allowedRoles := []string{"kyc_officer", "kyc_admin"}
	ctx := context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
	ctx = context.WithValue(ctx, cs.CtContextRealm, master)
	ctx = context.WithValue(ctx, cs.CtContextGroups, groups)

	t.Run("Error when getting admin configuration", func(t *testing.T) {
		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, nil)
		authorizationManager, err := NewIdentificationAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		mockAuthorizationDBReader.EXPECT().GetAdminConfiguration(ctx, targetRealm).Return(configuration.RealmAdminConfiguration{}, errors.New("Error"))

		err = authorizationManager.CheckRoleAuthorizationOnTargetUser(ctx, kycAction, targetRealm, targetUserID)
		assert.Equal(t, ForbiddenError{}, err)
	})

	t.Run("Error when getting user roles from Keycloak", func(t *testing.T) {
		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, nil)
		authorizationManager, err := NewIdentificationAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		mockAuthorizationDBReader.EXPECT().GetAdminConfiguration(ctx, targetRealm).Return(configuration.RealmAdminConfiguration{
			PhysicalIdentificationAllowedRoles: allowedRoles,
		}, nil)
		mockKeycloakClient.EXPECT().GetRoleNamesOfUser(ctx, accessToken, targetRealm, targetUserID).Return([]string{}, errors.New("Error"))

		err = authorizationManager.CheckRoleAuthorizationOnTargetUser(ctx, kycAction, targetRealm, targetUserID)
		assert.Equal(t, ForbiddenError{}, err)
	})

	t.Run("Authorized for KYC action with required role", func(t *testing.T) {
		userRoles := []string{"kyc_officer"}

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, nil)
		authorizationManager, err := NewIdentificationAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		mockAuthorizationDBReader.EXPECT().GetAdminConfiguration(ctx, targetRealm).Return(configuration.RealmAdminConfiguration{
			PhysicalIdentificationAllowedRoles: allowedRoles,
		}, nil)
		mockKeycloakClient.EXPECT().GetRoleNamesOfUser(ctx, accessToken, targetRealm, targetUserID).Return(userRoles, nil)

		err = authorizationManager.CheckRoleAuthorizationOnTargetUser(ctx, kycAction, targetRealm, targetUserID)
		assert.Nil(t, err)
	})

	t.Run("Unauthorized for KYC action without required role", func(t *testing.T) {
		userRoles := []string{"standard_user"}

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, nil)
		authorizationManager, err := NewIdentificationAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		mockAuthorizationDBReader.EXPECT().GetAdminConfiguration(ctx, targetRealm).Return(configuration.RealmAdminConfiguration{
			PhysicalIdentificationAllowedRoles: allowedRoles,
		}, nil)
		mockKeycloakClient.EXPECT().GetRoleNamesOfUser(ctx, accessToken, targetRealm, targetUserID).Return(userRoles, nil)

		err = authorizationManager.CheckRoleAuthorizationOnTargetUser(ctx, kycAction, targetRealm, targetUserID)
		assert.Equal(t, ForbiddenError{}, err)
	})

	t.Run("No restrictions configured - all roles allowed", func(t *testing.T) {
		userRoles := []string{"any_role"}

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, nil)
		authorizationManager, err := NewIdentificationAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		mockAuthorizationDBReader.EXPECT().GetAdminConfiguration(ctx, targetRealm).Return(configuration.RealmAdminConfiguration{
			// Empty PhysicalIdentificationAllowedRoles means all roles are allowed
		}, nil)
		mockKeycloakClient.EXPECT().GetRoleNamesOfUser(ctx, accessToken, targetRealm, targetUserID).Return(userRoles, nil).Times(1)

		err = authorizationManager.CheckRoleAuthorizationOnTargetUser(ctx, kycAction, targetRealm, targetUserID)
		assert.Nil(t, err)
	})
}

func TestCheckRoleAuthorizationOnSelfUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockKeycloakClient := mock.NewKeycloakClient(mockCtrl)
	mockAuthorizationDBReader := mock.NewAuthorizationDBReader(mockCtrl)

	accessToken := "TOKEN=="
	master := "master"
	groups := []string{"toe", "svc"}
	kycAction := KYCGetUser.String()
	allowedRoles := []string{"kyc_officer", "kyc_admin"}
	ctx := context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
	ctx = context.WithValue(ctx, cs.CtContextRealm, master)
	ctx = context.WithValue(ctx, cs.CtContextGroups, groups)

	t.Run("Error when getting admin configuration - empty user roles in context", func(t *testing.T) {
		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, nil)
		authorizationManager, err := NewIdentificationAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		mockAuthorizationDBReader.EXPECT().GetAdminConfiguration(ctx, master).Return(configuration.RealmAdminConfiguration{}, errors.New("Error"))

		err = authorizationManager.CheckRoleAuthorizationOnSelfUser(ctx, kycAction)
		assert.Equal(t, ForbiddenError{}, err)
	})

	t.Run("Authorized for KYC action with required role", func(t *testing.T) {
		ctx := context.WithValue(ctx, cs.CtContextRoles, []string{"kyc_officer"})

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, nil)
		authorizationManager, err := NewIdentificationAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		mockAuthorizationDBReader.EXPECT().GetAdminConfiguration(ctx, master).Return(configuration.RealmAdminConfiguration{
			PhysicalIdentificationAllowedRoles: allowedRoles,
		}, nil)

		err = authorizationManager.CheckRoleAuthorizationOnSelfUser(ctx, kycAction)
		assert.Nil(t, err)
	})

	t.Run("Unauthorized for KYC action without required role", func(t *testing.T) {
		ctx := context.WithValue(ctx, cs.CtContextRoles, []string{"standard_user"})

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, nil)
		authorizationManager, err := NewIdentificationAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		mockAuthorizationDBReader.EXPECT().GetAdminConfiguration(ctx, master).Return(configuration.RealmAdminConfiguration{
			PhysicalIdentificationAllowedRoles: allowedRoles,
		}, nil)

		err = authorizationManager.CheckRoleAuthorizationOnSelfUser(ctx, kycAction)
		assert.Equal(t, ForbiddenError{}, err)
	})

	t.Run("No restrictions configured - all roles allowed", func(t *testing.T) {
		ctx := context.WithValue(ctx, cs.CtContextRoles, []string{"any_role"})

		mockAuthorizationDBReader.EXPECT().GetAuthorizations(gomock.Any()).Return([]configuration.Authorization{}, nil)
		authorizationManager, err := NewIdentificationAuthorizationManager(mockAuthorizationDBReader, mockKeycloakClient, log.NewNopLogger())
		assert.Nil(t, err)

		mockAuthorizationDBReader.EXPECT().GetAdminConfiguration(ctx, master).Return(configuration.RealmAdminConfiguration{
			// Empty PhysicalIdentificationAllowedRoles means all roles are allowed
		}, nil)

		err = authorizationManager.CheckRoleAuthorizationOnSelfUser(ctx, kycAction)
		assert.Nil(t, err)
	})
}

func TestCheckRoleAuthorization(t *testing.T) {
	adminConfig := configuration.RealmAdminConfiguration{
		VideoIdentificationAllowedRoles:          []string{"end_user", "video_user"},
		AuxiliaryVideoIdentificationAllowedRoles: []string{"end_user"},
		AutoIdentificationAllowedRoles:           []string{"end_user"},
		PhysicalIdentificationAllowedRoles:       []string{},
	}

	videoCheck := checkRoleAuthorization(IDNAutoIdentInit.String(), adminConfig, []string{"end_user"})
	assert.True(t, videoCheck)

	auxiliaryVideoCheck := checkRoleAuthorization(IDNAuxiliaryVideoIdentInit.String(), adminConfig, []string{"l1_support_agent"})
	assert.False(t, auxiliaryVideoCheck)

	autoCheck := checkRoleAuthorization(IDNAutoIdentInit.String(), adminConfig, []string{"end_user", "some_other_role"})
	assert.True(t, autoCheck)

	check := checkRoleAuthorization(KYCGetUser.String(), adminConfig, []string{"any_role"})
	assert.True(t, check)
}
