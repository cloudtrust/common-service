package healthcheck

import "time"

// HealthDatabase allows to execute a simple one-row-result SQL query
type HealthDatabase interface {
	Ping() error
}

type databaseChecker struct {
	alias    string
	dbase    HealthDatabase
	response HealthStatus
}

// newDatabaseChecker creates a database health checker
func newDatabaseChecker(alias string, dbase HealthDatabase, cacheDuration time.Duration, timeProvider TimeProvider) BasicChecker {
	var database = "database"
	return &databaseChecker{
		alias:    alias,
		dbase:    dbase,
		response: HealthStatus{Name: &alias, Type: &database, CacheDuration: cacheDuration, TimeProvider: timeProvider},
	}
}

func (dbc *databaseChecker) CheckStatus() HealthStatus {
	if !dbc.response.hasExpired() {
		return dbc.response
	}

	err := dbc.dbase.Ping()

	if err != nil {
		dbc.response.stateDown(err.Error())
	} else {
		dbc.response.connection("established")
		dbc.response.stateUp()
	}

	dbc.response.touch()
	return dbc.response
}
