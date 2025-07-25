package configuration

import "encoding/json"

// Constants
const (
	CheckKeyIDNow          = "IDNow"
	CheckKeyIDNowAutoIdent = "IDNowAutoIdent"
	CheckKeyPhysical       = "physical-check"
)

var (
	// AvailableCheckKeys lists all available check keys for RealmAdminConfiguration
	AvailableCheckKeys = []string{CheckKeyIDNow, CheckKeyPhysical, CheckKeyIDNowAutoIdent}
)

// RealmConfiguration struct. APISelfAccountEditingEnabled replaces former field APISelfMailEditingEnabled
type RealmConfiguration struct {
	DefaultClientID                     *string  `json:"default_client_id,omitempty"`
	DefaultRedirectURI                  *string  `json:"default_redirect_uri,omitempty"`
	APISelfAuthenticatorDeletionEnabled *bool    `json:"api_self_authenticator_deletion_enabled,omitempty"`
	APISelfPasswordChangeEnabled        *bool    `json:"api_self_password_change_enabled,omitempty"`
	DeprecatedAPISelfMailEditingEnabled *bool    `json:"api_self_mail_editing_enabled,omitempty"`
	APISelfAccountEditingEnabled        *bool    `json:"api_self_account_editing_enabled,omitempty"`
	APISelfAccountDeletionEnabled       *bool    `json:"api_self_account_deletion_enabled,omitempty"`
	APISelfIDPLinksManagementEnabled    *bool    `json:"api_self_idplinks_management_enabled,omitempty"`
	ShowAuthenticatorsTab               *bool    `json:"show_authenticators_tab,omitempty"`
	ShowPasswordTab                     *bool    `json:"show_password_tab,omitempty"`
	ShowProfileTab                      *bool    `json:"show_profile_tab,omitempty"`
	ShowMailEditing                     *bool    `json:"show_mail_editing,omitempty"`
	ShowAccountDeletionButton           *bool    `json:"show_account_deletion_button,omitempty"`
	ShowIDPLinksTab                     *bool    `json:"show_idplinks_tab,omitempty"`
	SelfServiceDefaultTab               *string  `json:"self_service_default_tab,omitempty"`
	AllowedBackURL                      *string  `json:"allowed_back_url,omitempty"` // DEPRECATED
	AllowedBackURLs                     []string `json:"allowed_back_urls,omitempty"`

	RedirectCancelledRegistrationURL  *string   `json:"redirect_cancelled_registration_url,omitempty"`
	RedirectSuccessfulRegistrationURL *string   `json:"redirect_successful_registration_url,omitempty"`
	OnboardingRedirectURI             *string   `json:"onboarding_redirect_uri,omitempty"`
	OnboardingClientID                *string   `json:"onboarding_client_id,omitempty"`
	SelfRegisterGroupNames            *[]string `json:"self_register_group_names,omitempty"`

	BarcodeType *string `json:"barcode_type"`
}

// RealmAdminConfiguration struct
type RealmAdminConfiguration struct {
	Mode                                  *string         `json:"mode"`
	AvailableChecks                       map[string]bool `json:"available-checks,omitempty"`
	SelfRegisterEnabled                   *bool           `json:"self_register_enabled"`
	RegisterTheme                         *string         `json:"register_theme,omitempty"`
	SseTheme                              *string         `json:"sse_theme,omitempty"`
	BoTheme                               *string         `json:"bo_theme,omitempty"`
	SignerTheme                           *string         `json:"signer_theme,omitempty"`
	NeedVerifiedContact                   *bool           `json:"need_verified_contact,omitempty"`
	ConsentRequiredSocial                 *bool           `json:"consent_required,omitempty"`
	ConsentRequiredCorporate              *bool           `json:"consent_required_corp,omitempty"`
	ShowGlnEditing                        *bool           `json:"show_gln_editing,omitempty"`
	VideoIdentificationVoucherEnabled     *bool           `json:"video_identification_voucher_enabled"`
	VideoIdentificationAccountingEnabled  *bool           `json:"video_identification_accounting_enabled"`
	VideoIdentificationPrepaymentRequired *bool           `json:"video_identification_prepayment_required"`
	AutoIdentificationVoucherEnabled      *bool           `json:"auto_identification_voucher_enabled"`
	AutoIdentificationAccountingEnabled   *bool           `json:"auto_identification_accounting_enabled"`
	AutoIdentificationPrepaymentRequired  *bool           `json:"auto_identification_prepayment_required"`
	OnboardingStatusEnabled               *bool           `json:"onboarding_status_enabled"`
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
