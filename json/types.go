package json

import (
	"encoding/json"
	"fmt"
)

// OptionalString is a struct that represents a JSON string that can be
// undefined (Defined == false), null (Value == nil && Defined == true) or
// defined with a string value
type OptionalString struct {
	Defined bool
	Value   *string
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// When called, it means that the value is defined in the JSON payload.
func (os *OptionalString) UnmarshalJSON(data []byte) error {
	// UnmarshalJSON is called only if the key is present
	os.Defined = true
	return json.Unmarshal(data, &os.Value)
}

// MarshalJSON implements the json.Marshaler interface
func (os OptionalString) MarshalJSON() ([]byte, error) {
	// omitempty has no effect here
	if os.Defined && os.Value != nil {
		return []byte(fmt.Sprintf(`"%s"`, *os.Value)), nil
	}
	return []byte("null"), nil
}

// ToValue returns nil if not defined, the associated Value if set or defaultValue if defined as nil
func (os OptionalString) ToValue(defaultValue string) *string {
	return os.ToValueExt(nil, &defaultValue)
}

// ToValueExt returns the undefinedValue if not defined, the associated Value if set or defaultValue if defined as nil
func (os OptionalString) ToValueExt(undefinedValue *string, defaultValue *string) *string {
	if !os.Defined {
		return undefinedValue
	}
	if os.Value == nil {
		return defaultValue
	}
	return os.Value
}

// StringPtrToOptional creates an OptionalString using an input string pointer
func StringPtrToOptional(value *string) OptionalString {
	if value == nil {
		return OptionalString{Defined: false}
	}
	return OptionalString{
		Defined: true,
		Value:   value,
	}
}

// StringToOptional creates an OptionalString using an input string
func StringToOptional(value string) OptionalString {
	return OptionalString{
		Defined: true,
		Value:   &value,
	}
}
