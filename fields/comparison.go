package fields

import (
	"strings"

	csjson "github.com/cloudtrust/common-service/v2/json"
)

// FieldsComparator interface
type FieldsComparator interface {
	CaseSensitive(caseSensitive bool) FieldsComparator

	CompareValues(field Field, newValue, oldValue *string) FieldsComparator
	CompareValuesForUpdate(field Field, newValue, oldValue *string) FieldsComparator
	CompareValueAndFunction(field Field, newValue *string, oldValueFunc func(Field) []string) FieldsComparator
	CompareValueAndFunctionForUpdate(field Field, newValue *string, oldValueFunc func(Field) []string) FieldsComparator
	CompareOptionalAndFunction(field Field, newValue csjson.OptionalString, oldValueFunc func(Field) []string) FieldsComparator

	IsAnyFieldUpdated(fields ...Field) bool
	IsFieldUpdated(field Field) bool
	AreAllFieldsUpdated(fields ...Field) bool
	UpdatedFields() []string
}

type fieldsComparator struct {
	updatedFields map[Field]bool
	caseSensitive bool
}

// NewFieldsComparator create an instance of FieldsComparator
func NewFieldsComparator() FieldsComparator {
	return &fieldsComparator{updatedFields: make(map[Field]bool), caseSensitive: true}
}

func (fc *fieldsComparator) CaseSensitive(caseSensitive bool) FieldsComparator {
	fc.caseSensitive = caseSensitive
	return fc
}

func (fc *fieldsComparator) compareStrings(value1, value2 string) bool {
	if fc.caseSensitive {
		return value1 == value2
	}
	return strings.EqualFold(value1, value2)
}

func (fc *fieldsComparator) CompareValues(field Field, newValue, oldValue *string) FieldsComparator {
	if (newValue == nil && oldValue != nil) ||
		(newValue != nil && oldValue == nil) ||
		(newValue != nil && oldValue != nil && !fc.compareStrings(*newValue, *oldValue)) {
		fc.updatedFields[field] = true
	}
	return fc
}

func (fc *fieldsComparator) CompareValuesForUpdate(field Field, newValue, oldValue *string) FieldsComparator {
	// When specifying ForUpdate, we suppose that no update is asked if second value is nil
	if newValue == nil {
		return fc
	}
	return fc.CompareValues(field, newValue, oldValue)
}

func (fc *fieldsComparator) CompareValueAndFunction(field Field, newValue *string, oldValueFunc func(Field) []string) FieldsComparator {
	var oldValues = oldValueFunc(field)
	var l = len(oldValues)
	if l == 0 {
		return fc.CompareValues(field, newValue, nil)
	} else if l == 1 {
		return fc.CompareValues(field, newValue, &oldValues[0])
	}
	fc.updatedFields[field] = true
	return fc
}

func (fc *fieldsComparator) CompareValueAndFunctionForUpdate(field Field, newValue *string, oldValueFunc func(Field) []string) FieldsComparator {
	if newValue == nil {
		// Nothing to update
		return fc
	}
	return fc.CompareValueAndFunction(field, newValue, oldValueFunc)
}

func (fc *fieldsComparator) CompareOptionalAndFunction(field Field, newValue csjson.OptionalString, oldValueFunc func(Field) []string) FieldsComparator {
	if !newValue.Defined {
		// Nothing to update
		return fc
	}
	return fc.CompareValueAndFunction(field, newValue.Value, oldValueFunc)
}

// Return true if one of the given fields has been updated...
// If no field is provided in parameters, just check if any field has been updated
func (fc *fieldsComparator) IsAnyFieldUpdated(fields ...Field) bool {
	if len(fields) == 0 {
		return len(fc.updatedFields) > 0
	}
	for _, field := range fields {
		if _, ok := fc.updatedFields[field]; ok {
			return true
		}
	}
	return false
}

func (fc *fieldsComparator) IsFieldUpdated(field Field) bool {
	_, ok := fc.updatedFields[field]
	return ok
}

func (fc *fieldsComparator) AreAllFieldsUpdated(fields ...Field) bool {
	for _, field := range fields {
		if _, ok := fc.updatedFields[field]; !ok {
			return false
		}
	}
	return true
}

func (fc *fieldsComparator) UpdatedFields() []string {
	var fields []string
	for field := range fc.updatedFields {
		fields = append(fields, field.Key())
	}
	return fields
}
