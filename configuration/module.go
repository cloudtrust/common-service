package configuration

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/cloudtrust/common-service/database/sqltypes"
	"github.com/cloudtrust/common-service/log"
)

const (
	selectConfigStmt   = `SELECT configuration FROM realm_configuration WHERE (realm_id = ?)`
	selectAllAuthzStmt = `SELECT realm_id, group_name, action, target_realm_id, target_group_name FROM authorizations;`
)

type ConfigurationReaderDBModule struct {
	db     sqltypes.CloudtrustDB
	logger log.Logger
}

// NewConfigurationDBModule returns a ConfigurationDB module.
func NewConfigurationReaderDBModule(db sqltypes.CloudtrustDB, logger log.Logger) *ConfigurationReaderDBModule {
	return &ConfigurationReaderDBModule{
		db:     db,
		logger: logger,
	}
}

func (c *ConfigurationReaderDBModule) GetConfiguration(ctx context.Context, realmID string) (RealmConfiguration, error) {
	var configJSON string
	var config = RealmConfiguration{}
	row := c.db.QueryRow(selectConfigStmt, realmID)

	switch err := row.Scan(&configJSON); err {
	case sql.ErrNoRows:
		c.logger.Warn(ctx, "msg", "Realm Configuration not found in DB", "error", err.Error())
		return config, err

	default:
		if err != nil {
			return config, err
		}

		err = json.Unmarshal([]byte(configJSON), &config)
		return config, err
	}
}

func (c *ConfigurationReaderDBModule) GetAuthorizations(ctx context.Context) ([]Authorization, error) {
	// Get Authorizations from DB
	rows, err := c.db.Query(selectAllAuthzStmt)
	if err != nil {
		c.logger.Warn(ctx, "msg", "Can't get authorizations", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var authz Authorization
	var res = make([]Authorization, 0)
	for rows.Next() {
		authz, err = c.scanAuthorization(rows)
		if err != nil {
			c.logger.Warn(ctx, "msg", "Can't get authorizations. Scan failed", "error", err.Error())
			return nil, err
		}
		res = append(res, authz)
	}

	return res, nil
}

func (c *ConfigurationReaderDBModule) scanAuthorization(scanner sqltypes.SQLRow) (Authorization, error) {
	var (
		realmID         string
		groupName       string
		action          string
		targetGroupName sql.NullString
		targetRealmID   sql.NullString
	)

	err := scanner.Scan(&realmID, &groupName, &action, &targetRealmID, &targetGroupName)
	if err != nil {
		return Authorization{}, err
	}

	var authz = Authorization{
		RealmID:   &realmID,
		GroupName: &groupName,
		Action:    &action,
	}

	if targetRealmID.Valid {
		authz.TargetRealmID = &targetRealmID.String
	}

	if targetGroupName.Valid {
		authz.TargetGroupName = &targetGroupName.String
	}

	return authz, nil
}
