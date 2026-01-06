package security

import (
	"context"
	"encoding/json"
	"slices"
	"strings"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/cloudtrust/common-service/v2/configuration"
	"github.com/cloudtrust/common-service/v2/log"
)

// IdentificationKeycloakClient interface
type IdentificationKeycloakClient interface {
	GetRoleNamesOfUser(ctx context.Context, accessToken string, realmName, userID string) ([]string, error)
}

// IdentificationAuthorizationDBReader interface
type IdentificationAuthorizationDBReader interface {
	GetAdminConfiguration(ctx context.Context, realmID string) (configuration.RealmAdminConfiguration, error)
}

// IdentificationAuthorizationManager interface
type IdentificationAuthorizationManager interface {
	CheckRoleAuthorizationOnTargetUser(ctx context.Context, action string, targetRealm string, userID string) error
	CheckRoleAuthorizationOnSelfUser(ctx context.Context, action string) error
}

type identificationAuthorizationManager struct {
	authorizationDBReader IdentificationAuthorizationDBReader
	keycloakClient        IdentificationKeycloakClient
	logger                log.Logger
}

func NewIdentificationAuthorizationManager(authorizationDBReader IdentificationAuthorizationDBReader, keycloakClient IdentificationKeycloakClient, logger log.Logger) (IdentificationAuthorizationManager, error) {
	var manager = &identificationAuthorizationManager{
		authorizationDBReader: authorizationDBReader,
		keycloakClient:        keycloakClient,
		logger:                logger,
	}

	return manager, nil
}

// CheckRoleAuthorizationOnTargetUser checks if the target user has the required role to init identification
func (iam *identificationAuthorizationManager) CheckRoleAuthorizationOnTargetUser(ctx context.Context, action string, targetRealm string, userID string) error {
	var accessToken = ctx.Value(cs.CtContextAccessToken).(string)

	adminConfig, err := iam.authorizationDBReader.GetAdminConfiguration(ctx, targetRealm)
	if err != nil {
		iam.logger.Info(ctx, "msg", "ForbiddenError: "+err.Error(), "realm", targetRealm)
		return suggestForbiddenError(err)
	}

	userRoles, err := iam.keycloakClient.GetRoleNamesOfUser(ctx, accessToken, targetRealm, userID)
	if err != nil {
		iam.logger.Info(ctx, "msg", "ForbiddenError: "+err.Error(), "realm", targetRealm, "userID", userID)
		return suggestForbiddenError(err)
	}

	if checkRoleAuthorization(action, adminConfig, userRoles) {
		return nil
	}

	infos, _ := json.Marshal(map[string]string{
		"ThrownBy":    "CheckRoleAuthorizationOnTargetUser",
		"Action":      action,
		"targetRealm": targetRealm,
		"userRoles":   strings.Join(userRoles, "|"),
		"userID":      userID,
	})
	iam.logger.Info(ctx, "msg", "ForbiddenError: Not allowed to init identification", "infos", string(infos))
	return ForbiddenError{}
}

// CheckRoleAuthorizationOnSelfUser checks if the current user has the required role to init identification
func (iam *identificationAuthorizationManager) CheckRoleAuthorizationOnSelfUser(ctx context.Context, action string) error {
	currentRealm := ctx.Value(cs.CtContextRealm).(string)
	currentRoles, ok := ctx.Value(cs.CtContextRoles).([]string)
	if !ok {
		currentRoles = []string{}
	}

	adminConfig, err := iam.authorizationDBReader.GetAdminConfiguration(ctx, currentRealm)
	if err != nil {
		iam.logger.Info(ctx, "msg", "ForbiddenError: "+err.Error(), "realm", currentRealm)
		return suggestForbiddenError(err)
	}

	if checkRoleAuthorization(action, adminConfig, currentRoles) {
		return nil
	}

	infos, _ := json.Marshal(map[string]string{
		"ThrownBy":    "CheckRoleAuthorizationOnSelfUser",
		"Action":      action,
		"targetRealm": currentRealm,
		"userRoles":   strings.Join(currentRoles, "|"),
	})
	iam.logger.Info(ctx, "msg", "ForbiddenError: Not allowed to init identification", "infos", string(infos))
	return ForbiddenError{}
}

func checkRoleAuthorization(action string, adminConfig configuration.RealmAdminConfiguration, userRoles []string) bool {
	allowedRoles := []string{}
	switch action {
	case IDNVideoIdentInit.String():
		allowedRoles = adminConfig.VideoIdentificationAllowedRoles
	case IDNAuxiliaryVideoIdentInit.String():
		allowedRoles = adminConfig.AuxiliaryVideoIdentificationAllowedRoles
	case IDNAutoIdentInit.String():
		allowedRoles = adminConfig.AutoIdentificationAllowedRoles
	case KYCGetUser.String(), KYCValidateUser.String():
		allowedRoles = adminConfig.PhysicalIdentificationAllowedRoles
	}

	if len(allowedRoles) == 0 {
		// For retrocompatibility, if no restrictions are defined, all roles are allowed
		return true
	}

	for _, userRole := range userRoles {
		if slices.Contains(allowedRoles, userRole) {
			return true
		}
	}

	return false
}
