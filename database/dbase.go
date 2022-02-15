package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/cloudtrust/common-service/v2/database/sqltypes"
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

type basicCloudtrustDB struct {
	dbConn            *sql.DB
	pingTimeoutMillis time.Duration
}

func (db *basicCloudtrustDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (sqltypes.Transaction, error) {
	var tx, err = db.dbConn.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return NewTransaction(tx), nil
}

func (db *basicCloudtrustDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.dbConn.Exec(query, args...)
}

func (db *basicCloudtrustDB) Query(query string, args ...interface{}) (sqltypes.SQLRows, error) {
	return db.dbConn.Query(query, args...)
}

func (db *basicCloudtrustDB) QueryRow(query string, args ...interface{}) sqltypes.SQLRow {
	return db.dbConn.QueryRow(query, args...)
}

func (db *basicCloudtrustDB) Ping() error {
	if db.pingTimeoutMillis > 0 {
		var ctxTimeout, cancelTimeout = context.WithTimeout(context.Background(), time.Millisecond*db.pingTimeoutMillis)
		defer cancelTimeout()

		return db.dbConn.PingContext(ctxTimeout)
	}
	return db.dbConn.Ping()
}

func (db *basicCloudtrustDB) Close() error {
	return db.dbConn.Close()
}

func (db *basicCloudtrustDB) Stats() sql.DBStats {
	return db.dbConn.Stats()
}

// DbConfig Db configuration parameters
type DbConfig struct {
	HostPort          string
	Username          string
	Password          string
	Database          string
	Protocol          string
	Parameters        string
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   int
	Noop              bool
	MigrationEnabled  bool
	MigrationVersion  string
	ConnectionCheck   bool
	PingTimeoutMillis int
}

// ConfigureDbDefault configure default database parameters for a given prefix
// Parameters are built with the given prefix, then a dash symbol, then one of these suffixes:
// host-port, username, password, database, protocol, max-open-conns, max-idle-conns, conn-max-lifetime
// If a parameter exists only named with the given prefix and if its value if false, the database connection
// will be a Noop one
func ConfigureDbDefault(v cs.Configuration, prefix, envUser, envPasswd string) {
	v.SetDefault(prefix+"-enabled", true)
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
	v.SetDefault(prefix+"-connection-check", true)
	v.SetDefault(prefix+"-ping-timeout-ms", 1500)

	_ = v.BindEnv(prefix+"-username", envUser)
	_ = v.BindEnv(prefix+"-password", envPasswd)
}

// GetDbConfig reads db configuration parameters
// Check the parameter {prefix}-enabled to know if the connection to the database should be enabled
func GetDbConfig(v cs.Configuration, prefix string) *DbConfig {
	return GetDbConfigExt(v, prefix, !v.GetBool(prefix+"-enabled"))
}

// GetDbConfigExt is an extension of GetDbConfig for cases where we want to force the connection to be enabled/disabled
func GetDbConfigExt(v cs.Configuration, prefix string, noop bool) *DbConfig {
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
		cfg.ConnectionCheck = v.GetBool(prefix + "-connection-check")
		cfg.PingTimeoutMillis = v.GetInt(prefix + "-ping-timeout-ms")
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
func (cfg *DbConfig) OpenDatabase() (sqltypes.CloudtrustDB, error) {
	if cfg.Noop {
		return &NoopDB{}, nil
	}

	sqlConn, err := sql.Open("mysql", cfg.getDbConnectionString())
	if err != nil {
		return nil, err
	}
	dbConn := &basicCloudtrustDB{dbConn: sqlConn, pingTimeoutMillis: time.Duration(cfg.PingTimeoutMillis)}

	// DB migration version
	// checking that the flyway_schema_history has the minimum imposed migration version
	if cfg.MigrationEnabled {
		if cfg.MigrationVersion == "" {
			// DB schema versioning is enabled but no minimum version was given
			return nil, errors.New("Check of database schema is enabled, but no minimum version provided")
		}
		err = cfg.checkMigrationVersion(dbConn)
		if err != nil {
			_ = dbConn.Close()
			dbConn = nil
		}
	} else if cfg.ConnectionCheck {
		// Executes a simple query to check that the connection is valid
		err = dbConn.Ping()
	}

	// the config of the DB should have a max_connections > SetMaxOpenConns
	if err == nil {
		sqlConn.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlConn.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlConn.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}

	return dbConn, err
}

func (cfg *DbConfig) checkMigrationVersion(conn sqltypes.CloudtrustDB) error {
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

// ReconnectableCloudtrustDB implements an auto-reconnect mechanism
type ReconnectableCloudtrustDB struct {
	dbConnFactory sqltypes.CloudtrustDBFactory
	connection    sqltypes.CloudtrustDB
	mutex         *sync.Mutex
}

// NewReconnectableCloudtrustDB opens a connection to a database. This connection will be renewed if necessary
func NewReconnectableCloudtrustDB(dbConnFactory sqltypes.CloudtrustDBFactory) (sqltypes.CloudtrustDB, error) {
	dbConn, err := dbConnFactory.OpenDatabase()
	if err != nil {
		return nil, err
	}
	return &ReconnectableCloudtrustDB{
		dbConnFactory: dbConnFactory,
		connection:    dbConn,
		mutex:         &sync.Mutex{},
	}, nil
}

func (rcdb *ReconnectableCloudtrustDB) getActiveConnection() (sqltypes.CloudtrustDB, error) {
	var err error
	if rcdb.connection == nil {
		rcdb.mutex.Lock()
		// Ensure connection has not already been reopened by another thread
		if rcdb.connection == nil {
			rcdb.connection, err = rcdb.dbConnFactory.OpenDatabase()
		}
		rcdb.mutex.Unlock()
	}

	return rcdb.connection, err
}

func (rcdb *ReconnectableCloudtrustDB) resetConnection(reconnect bool) error {
	var err error

	if rcdb.connection != nil {
		rcdb.mutex.Lock()
		if rcdb.connection != nil {
			err = rcdb.connection.Close()
			rcdb.connection = nil
			if reconnect {
				// Reconnect later
				defer rcdb.asyncReconnect()
			}
		}
		rcdb.mutex.Unlock()
	}

	return err
}

func (rcdb *ReconnectableCloudtrustDB) asyncReconnect() {
	// Try to reconnect in asynchronously
	go rcdb.getActiveConnection()
}

func (rcdb *ReconnectableCloudtrustDB) checkError(err error) {
	connection := rcdb.connection
	if err != nil && connection != nil {
		switch err {
		case sql.ErrNoRows:
			return
		}
		if connection.Ping() != nil {
			_ = rcdb.resetConnection(true)
		}
	}
}

// BeginTx creates a transaction
func (rcdb *ReconnectableCloudtrustDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (sqltypes.Transaction, error) {
	dbConn, err := rcdb.getActiveConnection()
	if err != nil {
		return nil, err
	}
	var tx sqltypes.Transaction
	tx, err = dbConn.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// Exec an SQL query
func (rcdb *ReconnectableCloudtrustDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	dbConn, err := rcdb.getActiveConnection()
	if err != nil {
		return nil, err
	}

	var res sql.Result
	res, err = dbConn.Exec(query, args...)
	rcdb.checkError(err)

	return res, err
}

// Query a multiple-rows SQL result
func (rcdb *ReconnectableCloudtrustDB) Query(query string, args ...interface{}) (sqltypes.SQLRows, error) {
	dbConn, err := rcdb.getActiveConnection()
	if err != nil {
		return nil, err
	}

	var res sqltypes.SQLRows
	res, err = dbConn.Query(query, args...)
	rcdb.checkError(err)

	return res, err
}

// QueryRow a single-row SQL result
func (rcdb *ReconnectableCloudtrustDB) QueryRow(query string, args ...interface{}) sqltypes.SQLRow {
	dbConn, err := rcdb.getActiveConnection()
	if err != nil {
		return sqltypes.NewSQLRowError(err)
	}
	return dbConn.QueryRow(query, args...)
}

// Ping check the connection with the database
func (rcdb *ReconnectableCloudtrustDB) Ping() error {
	dbConn, err := rcdb.getActiveConnection()
	if err != nil {
		return err
	}
	err = dbConn.Ping()
	if err != nil {
		_ = rcdb.resetConnection(true)
	}
	return err
}

// Close the connection with the database
func (rcdb *ReconnectableCloudtrustDB) Close() error {
	return rcdb.resetConnection(false)
}

// Stats returns database statistics
func (rcdb *ReconnectableCloudtrustDB) Stats() sql.DBStats {
	dbConn, err := rcdb.getActiveConnection()
	if err != nil {
		return sql.DBStats{}
	}

	return dbConn.Stats()
}
