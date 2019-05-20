package database

import (
	"database/sql"
	"fmt"
	"time"

	cs "github.com/cloudtrust/common-service"
)

// CloudtrustDB interface
type CloudtrustDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
}

// DbConfig Db configuration parameters
type DbConfig struct {
	HostPort        string
	Username        string
	Password        string
	Database        string
	Protocol        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	Noop            bool
}

// ConfigureDbDefault configure default database parameters for a given prefix
// Parameters are built with the given prefix, then a dash symbol, then one of these suffixes:
// host-port, username, password, database, protocol, max-open-conns, max-idle-conns, conn-max-lifetime
// If a parameter exists only named with the given prefix and if its value if false, the database connection
// will be a Noop one
func ConfigureDbDefault(v cs.Configuration, prefix, envUser, envPasswd string) {
	v.SetDefault(prefix+"-host-port", "")
	v.SetDefault(prefix+"-username", "")
	v.SetDefault(prefix+"-password", "")
	v.SetDefault(prefix+"-database", "")
	v.SetDefault(prefix+"-protocol", "")
	v.SetDefault(prefix+"-max-open-conns", 10)
	v.SetDefault(prefix+"-max-idle-conns", 2)
	v.SetDefault(prefix+"-conn-max-lifetime", 3600)

	v.BindEnv(prefix+"-username", envUser)
	v.BindEnv(prefix+"-password", envPasswd)
}

// GetDbConfig reads db configuration parameters
func GetDbConfig(v cs.Configuration, prefix string, noop bool) *DbConfig {
	var cfg DbConfig

	cfg.Noop = noop
	if !noop {
		cfg.HostPort = v.GetString(prefix + "-host-port")
		cfg.Username = v.GetString(prefix + "-username")
		cfg.Password = v.GetString(prefix + "-password")
		cfg.Database = v.GetString(prefix + "-database")
		cfg.Protocol = v.GetString(prefix + "-protocol")
		cfg.MaxOpenConns = v.GetInt(prefix + "-max-open-conns")
		cfg.MaxIdleConns = v.GetInt(prefix + "-max-idle-conns")
		cfg.ConnMaxLifetime = v.GetInt(prefix + "-conn-max-lifetime")
	}

	return &cfg
}

// OpenDatabase gets an access to a database
// If cfg.Noop is true, a Noop access will be provided
func (cfg *DbConfig) OpenDatabase() (CloudtrustDB, error) {
	if cfg.Noop {
		return &NoopDB{}, nil
	}

	dbConn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s)/%s", cfg.Username, cfg.Password, cfg.Protocol, cfg.HostPort, cfg.Database))
	if err != nil {
		return &NoopDB{}, err
	}

	// the config of the DB should have a max_connections > SetMaxOpenConns
	if err == nil {
		dbConn.SetMaxOpenConns(cfg.MaxOpenConns)
		dbConn.SetMaxIdleConns(cfg.MaxIdleConns)
		dbConn.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}

	return dbConn, err
}

// NoopDB is a database client that does nothing.
type NoopDB struct{}

// Exec does nothing.
func (db *NoopDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return NoopResult{}, nil
}

// Query does nothing.
func (db *NoopDB) Query(query string, args ...interface{}) (*sql.Rows, error) { return nil, nil }

// QueryRow does nothing.
func (db *NoopDB) QueryRow(query string, args ...interface{}) *sql.Row { return nil }

// SetMaxOpenConns does nothing.
func (db *NoopDB) SetMaxOpenConns(n int) {}

// SetMaxIdleConns does nothing.
func (db *NoopDB) SetMaxIdleConns(n int) {}

// SetConnMaxLifetime does nothing.
func (db *NoopDB) SetConnMaxLifetime(d time.Duration) {}

// NoopResult is a sql.Result that does nothing.
type NoopResult struct{}

// LastInsertId does nothing.
func (NoopResult) LastInsertId() (int64, error) { return 0, nil }

// RowsAffected does nothing.
func (NoopResult) RowsAffected() (int64, error) { return 0, nil }
