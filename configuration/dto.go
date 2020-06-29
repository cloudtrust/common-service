package configuration

import "encoding/json"

// Constants
const (
	CheckKeyIDNow    = "IDNow"
	CheckKeyPhysical = "physical-check"
)

var (
	// AvailableCheckKeys lists all available check keys for RealmAdminConfiguration
	AvailableCheckKeys = []string{CheckKeyIDNow, CheckKeyPhysical}
)

// RealmConfiguration struct. APISelfAccountEditingEnabled replaces former field APISelfMailEditingEnabled
type RealmConfiguration struct {
	DefaultClientID                     *string   `json:"default_client_id,omitempty"`
	DefaultRedirectURI                  *string   `json:"default_redirect_uri,omitempty"`
	APISelfAuthenticatorDeletionEnabled *bool     `json:"api_self_authenticator_deletion_enabled,omitempty"`
	APISelfPasswordChangeEnabled        *bool     `json:"api_self_password_change_enabled,omitempty"`
	DeprecatedAPISelfMailEditingEnabled *bool     `json:"api_self_mail_editing_enabled,omitempty"`
	APISelfAccountEditingEnabled        *bool     `json:"api_self_account_editing_enabled,omitempty"`
	APISelfAccountDeletionEnabled       *bool     `json:"api_self_account_deletion_enabled,omitempty"`
	ShowAuthenticatorsTab               *bool     `json:"show_authenticators_tab,omitempty"`
	ShowPasswordTab                     *bool     `json:"show_password_tab,omitempty"`
	ShowProfileTab                      *bool     `json:"show_profile_tab,omitempty"`
	ShowMailEditing                     *bool     `json:"show_mail_editing,omitempty"`
	ShowAccountDeletionButton           *bool     `json:"show_account_deletion_button,omitempty"`
	RegisterExecuteActions              *[]string `json:"register_execute_actions,omitempty"`
	RedirectCancelledRegistrationURL    *string   `json:"redirect_cancelled_registration_url,omitempty"`
	RedirectSuccessfulRegistrationURL   *string   `json:"redirect_successful_registration_url,omitempty"`
	BarcodeType                         *string   `json:"barcode_type"`
}

// RealmAdminConfiguration struct
type RealmAdminConfiguration struct {
	Mode            *string                   `json:"mode"`
	AvailableChecks map[string]bool           `json:"available-checks,omitempty"`
	Accreditations  []RealmAdminAccreditation `json:"accreditations,omitempty"`
}

// RealmAdminAccreditation struct
type RealmAdminAccreditation struct {
	Type      *string `json:"type,omitempty"`
	Validity  *string `json:"validity,omitempty"`
	Condition *string `json:"condition,omitempty"`
}

// Authorization struct
type Authorization struct {
	RealmID         *string `json:"realm_id"`
	GroupName       *string `json:"group_id"`
	Action          *string `json:"action"`
	TargetRealmID   *string `json:"target_realm_id,omitempty"`
	TargetGroupName *string `json:"target_group_name,omitempty"`
}

// NewRealmConfiguration returns the realm configuration from its JSON representation
func NewRealmConfiguration(confJSON string) (RealmConfiguration, error) {
	var conf RealmConfiguration
	var err = json.Unmarshal([]byte(confJSON), &conf)
	if err != nil {
		return conf, err
	}
	if conf.DeprecatedAPISelfMailEditingEnabled != nil && conf.APISelfAccountEditingEnabled == nil {
		conf.APISelfAccountEditingEnabled = conf.DeprecatedAPISelfMailEditingEnabled
	}
	conf.DeprecatedAPISelfMailEditingEnabled = nil
	return conf, nil
}

// NewRealmAdminConfiguration returns the realm admin configuration from its JSON representation
func NewRealmAdminConfiguration(configJSON string) (RealmAdminConfiguration, error) {
	var conf RealmAdminConfiguration
	var err = json.Unmarshal([]byte(configJSON), &conf)
	return conf, err
}
