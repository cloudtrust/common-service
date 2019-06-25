package database

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	cs "github.com/cloudtrust/common-service"
)

// Define an internal structure to manage DB versions
type dbVersion struct {
	major int
	minor int
}

func newDbVersion(version string) (*dbVersion, error) {
	var r = regexp.MustCompile(`^(\d+)\.(\d+)$`)
	var match = r.FindStringSubmatch(version)
	if match == nil {
		return nil, fmt.Errorf("version %s does not match the required format", version)
	}
	// We don't test the Atoi errors as version matches the regexp
	// Major/Minor version numbers are stored in groups 1 and 2. Group 0 is the whole matched valued
	var maj, _ = strconv.Atoi(match[1])
	var min, _ = strconv.Atoi(match[2])
	return &dbVersion{
		major: maj,
		minor: min,
	}, nil
}

func (v *dbVersion) matchesRequired(required *dbVersion) bool {
	return !(v.major < required.major || (v.major == required.major && v.minor < required.minor))
}

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
	HostPort         string
	Username         string
	Password         string
	Database         string
	Protocol         string
	Parameters       string
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  int
	Noop             bool
	MigrationEnabled bool
	MigrationVersion string
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
	v.SetDefault(prefix+"-parameters", "")
	v.SetDefault(prefix+"-max-open-conns", 10)
	v.SetDefault(prefix+"-max-idle-conns", 2)
	v.SetDefault(prefix+"-conn-max-lifetime", 3600)
	v.SetDefault(prefix+"-migration", false)
	v.SetDefault(prefix+"-migration-version", "")

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
		cfg.Parameters = v.GetString(prefix + "-parameters")
		cfg.MaxOpenConns = v.GetInt(prefix + "-max-open-conns")
		cfg.MaxIdleConns = v.GetInt(prefix + "-max-idle-conns")
		cfg.ConnMaxLifetime = v.GetInt(prefix + "-conn-max-lifetime")
		cfg.MigrationEnabled = v.GetBool(prefix + "-migration")
		cfg.MigrationVersion = v.GetString(prefix + "-migration-version")
	}

	return &cfg
}

func (cfg *DbConfig) getDbConnectionString() string {
	var separ = ""
	if len(cfg.Parameters) > 0 {
		separ = "?"
	}
	return fmt.Sprintf("%s:%s@%s(%s)/%s%s%s", cfg.Username, cfg.Password, cfg.Protocol, cfg.HostPort, cfg.Database, separ, cfg.Parameters)
}

// OpenDatabase gets an access to a database
// If cfg.Noop is true, a Noop access will be provided
func (cfg *DbConfig) OpenDatabase() (CloudtrustDB, error) {
	if cfg.Noop {
		return &NoopDB{}, nil
	}

	dbConn, err := sql.Open("mysql", cfg.getDbConnectionString())
	if err != nil {
		return nil, err
	}

	// DB migration version
	// checking that the flyway_schema_history has the minimum imposed migration version
	if cfg.MigrationEnabled {
		if cfg.MigrationVersion == "" {
			// DB schema versioning is enabled but no minimum version was given
			return nil, errors.New("Check of database schema is enabled, but no minimum version provided")
		}
		err = cfg.checkMigrationVersion(dbConn)
		if err != nil {
			dbConn.Close()
			dbConn = nil
		}
	}

	// the config of the DB should have a max_connections > SetMaxOpenConns
	if err == nil {
		dbConn.SetMaxOpenConns(cfg.MaxOpenConns)
		dbConn.SetMaxIdleConns(cfg.MaxIdleConns)
		dbConn.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}

	return dbConn, err
}

func (cfg *DbConfig) checkMigrationVersion(conn CloudtrustDB) error {
	var flywayVersion string
	row := conn.QueryRow(`SELECT version FROM flyway_schema_history ORDER BY installed_rank DESC LIMIT 1`)
	err := row.Scan(&flywayVersion)

	//flyway version and db migration version must match the format x.y where x and y are integers
	var requiredDbVersion, flywayDbVersion *dbVersion
	if err == nil {
		requiredDbVersion, err = newDbVersion(cfg.MigrationVersion)
	}
	if err == nil {
		flywayDbVersion, err = newDbVersion(flywayVersion)
	}

	// compare the two versions of the type x.y (major.minor)

	// it is required for the last script version of the flyway is to be "bigger" than the required version
	if err == nil && !flywayDbVersion.matchesRequired(requiredDbVersion) {
		err = fmt.Errorf("Database schema not up-to-date (current: %s, required: %s)", flywayVersion, cfg.MigrationVersion)
	}

	return err
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
