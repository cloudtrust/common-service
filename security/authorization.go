package security

import (
	"context"
	"encoding/json"
	"strings"

	cs "github.com/cloudtrust/common-service"
	"github.com/cloudtrust/common-service/configuration"
	"github.com/cloudtrust/common-service/log"
)

func (am *authorizationManager) CheckAuthorizationOnTargetUser(ctx context.Context, action, targetRealm, userID string) error {
	var accessToken = ctx.Value(cs.CtContextAccessToken).(string)

	infos, _ := json.Marshal(map[string]string{
		"ThrownBy":    "CheckAuthorizationOnTargetUser",
		"Action":      action,
		"targetRealm": targetRealm,
		"userID":      userID,
	})

	// Retrieve the group of the target user

	var groupsRep []string
	var err error
	if groupsRep, err = am.keycloakClient.GetGroupNamesOfUser(ctx, accessToken, targetRealm, userID); err != nil {
		am.logger.Info(ctx, "msg", "ForbiddenError: "+err.Error(), "infos", string(infos))
		return ForbiddenError{}
	}

	if groupsRep == nil || len(groupsRep) == 0 {
		// No groups assigned, nothing allowed
		am.logger.Info(ctx, "msg", "ForbiddenError: No groups assigned to this user, nothing allowed", "infos", string(infos))
		return ForbiddenError{}
	}

	for _, targetGroup := range groupsRep {
		if am.CheckAuthorizationOnTargetGroup(ctx, action, targetRealm, targetGroup) == nil {
			return nil
		}
	}

	am.logger.Info(ctx, "msg", "ForbiddenError: Not allowed to perform the action on user with such groups", "infos", string(infos))
	return ForbiddenError{}
}

func (am *authorizationManager) CheckAuthorizationOnTargetGroupID(ctx context.Context, action, targetRealm, targetGroupID string) error {
	var accessToken = ctx.Value(cs.CtContextAccessToken).(string)
	var currentRealm = ctx.Value(cs.CtContextRealm).(string)
	var currentGroups = ctx.Value(cs.CtContextGroups).([]string)

	infos, _ := json.Marshal(map[string]string{
		"ThrownBy":      "CheckAuthorizationOnTargetGroupID",
		"Action":        action,
		"targetRealm":   targetRealm,
		"targetGroupID": targetGroupID,
		"currentRealm":  currentRealm,
		"currentGroups": strings.Join(currentGroups, "|"),
	})

	// Retrieve the name of the target group
	var err error
	var targetGroup string
	if targetGroup, err = am.keycloakClient.GetGroupName(ctx, accessToken, targetRealm, targetGroupID); err != nil {
		am.logger.Info(ctx, "msg", "ForbiddenError: "+err.Error(), "infos", string(infos))
		return ForbiddenError{}
	}

	if targetGroup == "" {
		am.logger.Info(ctx, "msg", "ForbiddenError: Group not found", "infos", string(infos))
		return ForbiddenError{}
	}

	return am.CheckAuthorizationOnTargetGroup(ctx, action, targetRealm, targetGroup)
}

func (am *authorizationManager) CheckAuthorizationOnTargetGroup(ctx context.Context, action, targetRealm, targetGroup string) error {
	var currentRealm = ctx.Value(cs.CtContextRealm).(string)
	var currentGroups = ctx.Value(cs.CtContextGroups).([]string)

	infos, _ := json.Marshal(map[string]string{
		"ThrownBy":      "CheckAuthorizationOnTargetGroup",
		"Action":        action,
		"targetRealm":   targetRealm,
		"targetGroup":   targetGroup,
		"currentRealm":  currentRealm,
		"currentGroups": strings.Join(currentGroups, "|"),
	})

	for _, group := range currentGroups {
		targetGroupAllowed, wildcard := (*am.authorizations)[currentRealm][group][action]["*"]

		if wildcard {
			_, allGroupsAllowed := targetGroupAllowed["*"]
			_, groupAllowed := targetGroupAllowed[targetGroup]

			if allGroupsAllowed || groupAllowed {
				return nil
			}
		}

		targetGroupAllowed, nonMasterRealmAllowed := (*am.authorizations)[currentRealm][group][action]["/"]

		if targetRealm != "master" && nonMasterRealmAllowed {
			_, allGroupsAllowed := targetGroupAllowed["*"]
			_, groupAllowed := targetGroupAllowed[targetGroup]

			if allGroupsAllowed || groupAllowed {
				return nil
			}
		}

		targetGroupAllowed, realmAllowed := (*am.authorizations)[currentRealm][group][action][targetRealm]

		if realmAllowed {
			_, allGroupsAllowed := targetGroupAllowed["*"]
			_, groupAllowed := targetGroupAllowed[targetGroup]

			if allGroupsAllowed || groupAllowed {
				return nil
			}
		}
	}

	am.logger.Info(ctx, "msg", "ForbiddenError: Not allowed to perform the action on this group", "infos", string(infos))
	return ForbiddenError{}
}

func (am *authorizationManager) CheckAuthorizationOnTargetRealm(ctx context.Context, action, targetRealm string) error {
	var currentRealm = ctx.Value(cs.CtContextRealm).(string)
	var currentGroups = ctx.Value(cs.CtContextGroups).([]string)

	infos, _ := json.Marshal(map[string]string{
		"ThrownBy":      "CheckAuthorizationOnTargetRealm",
		"Action":        action,
		"targetRealm":   targetRealm,
		"currentRealm":  currentRealm,
		"currentGroups": strings.Join(currentGroups, "|"),
	})

	for _, group := range currentGroups {
		_, wildcard := (*am.authorizations)[currentRealm][group][action]["*"]
		_, nonMasterRealmAllowed := (*am.authorizations)[currentRealm][group][action]["/"]
		_, realmAllowed := (*am.authorizations)[currentRealm][group][action][targetRealm]

		if wildcard || realmAllowed || (targetRealm != "master" && nonMasterRealmAllowed) {
			return nil
		}
	}

	am.logger.Info(ctx, "msg", "ForbiddenError: Not allowed to perform the action on this realm", "infos", string(infos))

	return ForbiddenError{}
}

// GetRightsOfCurrentUser returns the matrix rights of the current user
func (am *authorizationManager) GetRightsOfCurrentUser(ctx context.Context) map[string]map[string]map[string]map[string]struct{} {
	var currentRealm string
	var currentGroups = []string{}
	var currentRealmVal = ctx.Value(cs.CtContextRealm)
	var currentGroupsVal = ctx.Value(cs.CtContextGroups)

	if currentRealmVal != nil {
		currentRealm = currentRealmVal.(string)
	}

	if currentGroupsVal != nil {
		currentGroups = currentGroupsVal.([]string)
	}

	//3 dimensions table to express authorizations (group_of_user, action, target_realm) -> target_group for which the action is allowed
	// We keep group_of_user as a user may be part of multiple groups
	var rights = map[string]map[string]map[string]map[string]struct{}{}

	for _, group := range currentGroups {
		rightsForGroup, exist := (*am.authorizations)[currentRealm][group]

		if exist {
			rights[group] = rightsForGroup
		}
	}

	return rights
}

// ForbiddenError when an operation is not permitted.
type ForbiddenError struct{}

func (e ForbiddenError) Error() string {
	return "ForbiddenError: Operation not permitted"
}

// AuthorizationsMatrix data structure
// 4 dimensions table to express authorizations (realm_of_user, group_of_user, action, target_realm) -> target_group for which the action is allowed
type AuthorizationsMatrix map[string]map[string]map[string]map[string]map[string]struct{}

// LoadAuthorizations loads the authorization JSON into the data structure
// Authorization matrix is a 4 dimensions table :
//   - realm_of_user
//   - role_of_user
//   - action
//   - target_realm
// -> target_groups for which the action is allowed
//
// Note:
//   '*' can be used to express all target realms
//   '/' can be used to express all non master realms
//   '*' can be used to express all target groups are allowed
func (am *authorizationManager) ReloadAuthorizations(ctx context.Context) error {
	authorizations, err := am.authorizationDBReader.GetAuthorizations(context.Background())
	if err != nil {
		return err
	}

	var matrix = make(AuthorizationsMatrix)

	for _, authz := range authorizations {
		// Realm of user
		if _, ok := matrix[*authz.RealmID]; !ok {
			matrix[*authz.RealmID] = make(map[string]map[string]map[string]map[string]struct{})
		}

		// Group of user
		if _, ok := matrix[*authz.RealmID][*authz.GroupName]; !ok {
			matrix[*authz.RealmID][*authz.GroupName] = make(map[string]map[string]map[string]struct{})
		}

		// Action
		if _, ok := matrix[*authz.RealmID][*authz.GroupName][*authz.Action]; !ok {
			matrix[*authz.RealmID][*authz.GroupName][*authz.Action] = make(map[string]map[string]struct{})
		}

		// Target Realm
		if authz.TargetRealmID == nil {
			continue
		}

		if _, ok := matrix[*authz.RealmID][*authz.GroupName][*authz.Action][*authz.TargetRealmID]; !ok {
			matrix[*authz.RealmID][*authz.GroupName][*authz.Action][*authz.TargetRealmID] = make(map[string]struct{})
		}

		// Target Group
		if authz.TargetGroupName == nil {
			continue
		}

		matrix[*authz.RealmID][*authz.GroupName][*authz.Action][*authz.TargetRealmID][*authz.TargetGroupName] = struct{}{}
	}

	am.authorizations = &matrix

	return nil
}

type authorizationManager struct {
	authorizations        *AuthorizationsMatrix
	authorizationDBReader AuthorizationDBReader
	keycloakClient        KeycloakClient
	logger                log.Logger
}

// KeycloakClient is the minimum interface required to access Keycloak
type KeycloakClient interface {
	GetGroupNamesOfUser(ctx context.Context, accessToken string, realmName, userID string) ([]string, error)
	GetGroupName(ctx context.Context, accessToken string, realmName, groupID string) (string, error)
}

// AuthorizationDBReader interface
type AuthorizationDBReader interface {
	GetAuthorizations(context.Context) ([]configuration.Authorization, error)
}

// AuthorizationManager interface
type AuthorizationManager interface {
	CheckAuthorizationOnTargetRealm(ctx context.Context, action, targetRealm string) error
	CheckAuthorizationOnTargetGroup(ctx context.Context, action, targetRealm, targetGroup string) error
	CheckAuthorizationOnTargetGroupID(ctx context.Context, action, targetRealm, targetGroupID string) error
	CheckAuthorizationOnTargetUser(ctx context.Context, action, targetRealm, userID string) error
	GetRightsOfCurrentUser(ctx context.Context) map[string]map[string]map[string]map[string]struct{}
	ReloadAuthorizations(ctx context.Context) error
}

// Authorizations data structure

// NewAuthorizationManager loads the authorization from DB into a cache data structure and create an AuthorizationManager instance.
// Authorization matrix is a 4 dimensions table :
//   - realm_of_user
//   - role_of_user
//   - action
//   - target_realm
// -> target_groups for which the action is allowed
//
// Note:
//   '*' can be used to express all target realms
//   '/' can be used to express all non master realms
//   '*' can be used to express all target groups are allowed
func NewAuthorizationManager(authorizationDBReader AuthorizationDBReader, keycloakClient KeycloakClient, logger log.Logger) (AuthorizationManager, error) {
	var manager = &authorizationManager{
		authorizationDBReader: authorizationDBReader,
		keycloakClient:        keycloakClient,
		logger:                logger,
	}

	err := manager.ReloadAuthorizations(context.Background())
	if err != nil {
		return nil, err
	}

	return manager, nil
}

// Action type
type Action struct {
	Name  string
	Scope Scope
}

// Scope type
type Scope string

// Scope values
var (
	ScopeGlobal = Scope("global")
	ScopeRealm  = Scope("realm")
	ScopeGroup  = Scope("group")
)

func (a Action) String() string {
	return a.Name
}
