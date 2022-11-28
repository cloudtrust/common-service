package fields

// Field interface
type Field interface {
	Key() string
	AttributeName() string
}

type field struct {
	Field     string
	Attribute string
}

// GetKnownFields gets all known fields
func GetKnownFields() []Field {
	return allFields
}

func genericCreateField(key string, attrb string) Field {
	var res = &field{
		Field:     key,
		Attribute: attrb,
	}
	allFields = append(allFields, res)
	return res
}

func createField(name string) Field {
	return genericCreateField(name, name)
}

func createFieldPII(name string) Field {
	return genericCreateField(name, "ENC_"+name)
}

func createNonAttributeField(name string) Field {
	return genericCreateField(name, "")
}

func (f *field) Key() string {
	return f.Field
}

func (f *field) AttributeName() string {
	return f.Attribute
}

// Fields
var (
	allFields []Field

	Accreditations        = createField("accreditations")
	BirthDate             = createFieldPII("birthDate")
	BirthLocation         = createFieldPII("birthLocation")
	BusinessID            = createField("businessID")
	Email                 = createNonAttributeField("email")
	EmailToValidate       = createField("emailToValidate")
	FirstName             = createNonAttributeField("firstName")
	Gender                = createFieldPII("gender")
	IDDocumentCountry     = createFieldPII("idDocumentCountry")
	IDDocumentExpiration  = createFieldPII("idDocumentExpiration")
	IDDocumentNumber      = createFieldPII("idDocumentNumber")
	IDDocumentType        = createFieldPII("idDocumentType")
	Label                 = createField("label")
	LastName              = createNonAttributeField("lastName")
	Locale                = createField("locale")
	NameID                = createField("saml.persistent.name.id.for.*")
	Nationality           = createFieldPII("nationality")
	OnboardingCompleted   = createField("onboardingCompleted")
	PendingChecks         = createField("pendingChecks")
	PhoneNumber           = createField("phoneNumber")
	PhoneNumberToValidate = createField("phoneNumberToValidate")
	PhoneNumberVerified   = createField("phoneNumberVerified")
	SmsAttempts           = createField("smsAttempts")
	SmsSent               = createField("smsSent")
	Source                = createField("src")
	TrustIDAuthToken      = createField("trustIDAuthToken")
	TrustIDGroups         = createField("trustIDGroups")
)
