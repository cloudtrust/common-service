package commonservice

// CtContext type
type CtContext int

const (
	// CtContextAccessToken is the access token context key
	CtContextAccessToken CtContext = iota
	// CtContextRealm is the realm name context key
	CtContextRealm CtContext = iota
	// CtContextRealmID is the realm id context key
	CtContextRealmID CtContext = iota
	// CtContextUserID is the user id context key
	CtContextUserID CtContext = iota
	// CtContextUsername is the username context key
	CtContextUsername CtContext = iota
	// CtContextGroups is the groups context key
	CtContextGroups CtContext = iota
	// CtContextCorrelationID is the correlation id context key
	CtContextCorrelationID CtContext = iota
	// CtContextIssuerDomain is the issuer domain context key
	CtContextIssuerDomain CtContext = iota
	// CtContextRoles is the roles context key
	CtContextRoles CtContext = iota
)
