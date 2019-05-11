package database

import (
	"context"
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

// EventsDBModule is the interface of the audit events module.
type EventsDBModule interface {
	Store(context.Context, map[string]string) error
	ReportEvent(ctx context.Context, apiCall string, origin string, values ...string) error
}

type eventsDBModule struct {
	db CloudtrustDB
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
	if m["ct_event_type"] == "" {
		return nil
	}

	// the event was already formatted according to the DB structure already at the component level

	//auditTime - time of the event
	auditTime := m["audit_time"]
	// origin - the component that initiated the event
	origin := m["origin"]
	// realmName - realm name of the user that is impacted by the action
	realmName := m["realm_name"]
	//agentUserID - userId of who is performing an action
	agentUserID := m["agent_user_id"]
	//agentUsername - username of who is performing an action
	agentUsername := m["agent_username"]
	//agentRealmName - realm of who is performing an action
	agentRealmName := m["agent_realm_name"]
	//userID - ID of the user that is impacted by the action
	userID := m["user_id"]
	//username - username of the user that is impacted by the action
	username := m["username"]
	// ctEventType that  is established before at the component level
	ctEventType := m["ct_event_type"]
	// kcEventType corresponds to keycloak event type
	kcEventType := m["kc_event_type"]
	// kcOperationType - operation type of the event that comes from Keycloak
	kcOperationType := m["kc_operation_type"]
	// Id of the client
	clientID := m["client_id"]
	//additional_info - all the rest of the information from the event
	additionalInfo := m["additional_info"]

	//store the event in the DB
	_, err := cm.db.Exec(insertEvent, auditTime, origin, realmName, agentUserID, agentUsername, agentRealmName, userID, username, ctEventType, kcEventType, kcOperationType, clientID, additionalInfo)

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
	event.details["ct_event_type"] = apiCall
	event.details["origin"] = origin
	event.details["audit_time"] = time.Now().UTC().Format(timeFormat)

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
	//retrieve agent username
	er.details["agent_username"] = ctx.Value(cs.CtContextUsername).(string)
	//retrieve agent user id - not yet implemented
	//to be uncommented once the ctx contains the userId value
	//er.details["userId"] = ctx.Value(cs.CtContextUserID).(string)
	//retrieve agent realm
	er.details["agent_realm_name"] = ctx.Value(cs.CtContextRealm).(string)
}
