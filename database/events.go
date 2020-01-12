package database

import (
	"context"
	"encoding/json"
	"time"

	cs "github.com/cloudtrust/common-service"
)

const (
	timeFormat  = "2006-01-02 15:04:05.000"
	insertEvent = `INSERT INTO audit (
		audit_time,
		origin,
		realm_name,
		agent_user_id,
		agent_username,
		agent_realm_name,
		user_id,
		username,
		ct_event_type,
		kc_event_type,
		kc_operation_type,
		client_id,
		additional_info) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
)

// Defines event information constants
const (
	CtEventType            = "ct_event_type"
	CtEventAgentUsername   = "agent_username"
	CtEventAgentRealmName  = "agent_realm_name"
	CtEventUserID          = "user_id"
	CtEventGroupID         = "group_id"
	CtEventGroupName       = "group_name"
	CtEventOrigin          = "origin"
	CtEventAuditTime       = "audit_time"
	CtEventRealmName       = "realm_name"
	CtEventAgentUserID     = "agent_user_id"
	CtEventUsername        = "username"
	CtEventKcEventType     = "kc_event_type"
	CtEventKcOperationType = "kc_operation_type"
	CtEventClientID        = "client_id"
	CtEventAdditionalInfo  = "additional_info"
)

var ctEventColumns = []string{
	CtEventType, CtEventAgentUsername, CtEventAgentRealmName, CtEventUserID, CtEventOrigin, CtEventAuditTime, CtEventRealmName,
	CtEventAgentUserID, CtEventUsername, CtEventKcEventType, CtEventKcOperationType, CtEventClientID, CtEventAdditionalInfo}

// EventsDBModule is the interface of the audit events module.
type EventsDBModule interface {
	Store(context.Context, map[string]string) error
	ReportEvent(ctx context.Context, apiCall string, origin string, values ...string) error
}

type eventsDBModule struct {
	db CloudtrustDB
}

func isInArray(array []string, value string) bool {
	for _, e := range array {
		if e == value {
			return true
		}
	}
	return false
}

func checkNull(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

// CreateAdditionalInfo creates the additional info value
func CreateAdditionalInfo(values ...string) string {
	var nfo = make(map[string]string)
	for i := 0; i+1 < len(values); i += 2 {
		nfo[values[i]] = values[i+1]
	}
	addInfo, _ := json.Marshal(nfo)
	return string(addInfo)
}

// NewEventsDBModule returns a Console module.
func NewEventsDBModule(db CloudtrustDB) EventsDBModule {
	//db.Exec(createTable)
	return &eventsDBModule{
		db: db,
	}
}

func (cm *eventsDBModule) Store(_ context.Context, m map[string]string) error {
	// if ctEventType is not "", then record the events in MariaDB
	// otherwise, do nothing
	if m[CtEventType] == "" {
		return nil
	}

	// the event was already formatted according to the DB structure already at the component level

	//auditTime - time of the event
	auditTime := m[CtEventAuditTime]
	// origin - the component that initiated the event
	origin := m[CtEventOrigin]
	// realmName - realm name of the user that is impacted by the action
	realmName := m[CtEventRealmName]
	//agentUserID - userId of who is performing an action
	agentUserID := m[CtEventAgentUserID]
	//agentUsername - username of who is performing an action
	agentUsername := m[CtEventAgentUsername]
	//agentRealmName - realm of who is performing an action
	agentRealmName := m[CtEventAgentRealmName]
	//userID - ID of the user that is impacted by the action
	userID := m[CtEventUserID]
	//username - username of the user that is impacted by the action
	username := m[CtEventUsername]
	// ctEventType that  is established before at the component level
	ctEventType := m[CtEventType]
	// kcEventType corresponds to keycloak event type
	kcEventType := m[CtEventKcEventType]
	// kcOperationType - operation type of the event that comes from Keycloak
	kcOperationType := m[CtEventKcOperationType]
	// Id of the client
	clientID := m[CtEventClientID]
	//additional_info - all the rest of the information from the event
	additionalInfo := m[CtEventAdditionalInfo]
	if additionalInfo == "" {
		var addNfo = make(map[string]string)
		for k, v := range m {
			if !isInArray(ctEventColumns, k) {
				addNfo[k] = v
			}
		}
		if additionalInfoBytes, err := json.Marshal(addNfo); err == nil && len(addNfo) > 0 {
			additionalInfo = string(additionalInfoBytes)
		}
	}

	//store the event in the DB
	_, err := cm.db.Exec(insertEvent, auditTime, origin, checkNull(realmName), checkNull(agentUserID), checkNull(agentUsername),
		checkNull(agentRealmName), checkNull(userID), checkNull(username), checkNull(ctEventType), checkNull(kcEventType),
		checkNull(kcOperationType), checkNull(clientID), checkNull(additionalInfo))

	return err
}

// ReportEvent Report the event into the specified eventStorer
func (cm *eventsDBModule) ReportEvent(ctx context.Context, apiCall string, origin string, values ...string) error {
	event := CreateEvent(apiCall, origin)
	event.AddAgentDetails(ctx)
	event.AddEventValues(values...)
	return cm.Store(ctx, event.details)
}

// EventStorer interface of a
type EventStorer interface {
	Store(context.Context, map[string]string) error
}

// ReportEventDetails information of an event to be reported
type ReportEventDetails struct {
	details map[string]string
}

// CreateEvent create the generic event that contains the ct_event_type, origin and audit_time
func CreateEvent(apiCall string, origin string) ReportEventDetails {
	var event ReportEventDetails
	event.details = make(map[string]string)
	event.details[CtEventType] = apiCall
	event.details[CtEventOrigin] = origin
	event.details[CtEventAuditTime] = time.Now().UTC().Format(timeFormat)

	return event
}

// AddEventValues enhance the event with more information
func (er *ReportEventDetails) AddEventValues(values ...string) {
	//add information to the event
	noTuples := len(values)
	for i := 0; i+1 < noTuples; i = i + 2 {
		er.details[values[i]] = values[i+1]
	}
}

// AddAgentDetails add details from the context
func (er *ReportEventDetails) AddAgentDetails(ctx context.Context) {
	var mapper = map[cs.CtContext]string{
		cs.CtContextUsername: CtEventAgentUsername,
		cs.CtContextUserID:   CtEventUserID,
		cs.CtContextRealm:    CtEventAgentRealmName,
	}
	for keyFrom, keyTo := range mapper {
		var value = ctx.Value(keyFrom)
		if value != nil {
			er.details[keyTo] = value.(string)
		}
	}
}
