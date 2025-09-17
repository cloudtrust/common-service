package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAction(t *testing.T) {
	var action = Action{
		Name:  "test",
		Scope: ScopeGlobal,
	}

	assert.Equal(t, "test", action.String())
}

func TestActionsIndex(t *testing.T) {
	assert.Len(t, Actions.GetActionsForAPIs(BridgeService, ManagementAPI), len(Actions.index[BridgeService][ManagementAPI]))
	assert.Equal(t, Actions.index[BridgeService][ManagementAPI], Actions.GetActionsForAPIs(BridgeService, ManagementAPI))
	assert.Len(t, Actions.GetActionNamesForService(BridgeService), len(Actions.GetActionsForAPIs(BridgeService, CommunicationAPI, EventsAPI, KycAPI, ManagementAPI, StatisticAPI, TaskAPI, IdpAPI, ComponentsAPI)))
}
