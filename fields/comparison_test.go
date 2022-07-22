package fields

import (
	"strings"
	"testing"

	csjson "github.com/cloudtrust/common-service/v2/json"
	"github.com/stretchr/testify/assert"
)

func TestFieldsComparator(t *testing.T) {
	var fieldEmail = createField("EMAIL")
	var value1 = "value1"
	var value2 = "value2"

	t.Run("CompareValues", func(t *testing.T) {
		t.Run("Both values are nil", func(t *testing.T) {
			assert.Len(t, NewFieldsComparator().CompareValues(fieldEmail, nil, nil).UpdatedFields(), 0)
			assert.Len(t, NewFieldsComparator().CompareValuesForUpdate(fieldEmail, nil, nil).UpdatedFields(), 0)
		})
		t.Run("First value is nil", func(t *testing.T) {
			assert.Len(t, NewFieldsComparator().CompareValues(fieldEmail, nil, &value1).UpdatedFields(), 1)
			assert.Len(t, NewFieldsComparator().CompareValuesForUpdate(fieldEmail, nil, &value1).UpdatedFields(), 0)
		})
		t.Run("Second value is nil", func(t *testing.T) {
			assert.Len(t, NewFieldsComparator().CompareValues(fieldEmail, &value1, nil).UpdatedFields(), 1)
			assert.Len(t, NewFieldsComparator().CompareValuesForUpdate(fieldEmail, &value1, nil).UpdatedFields(), 1)
		})
		t.Run("Values are equal", func(t *testing.T) {
			assert.Len(t, NewFieldsComparator().CompareValues(fieldEmail, &value1, &value1).UpdatedFields(), 0)
			assert.Len(t, NewFieldsComparator().CompareValuesForUpdate(fieldEmail, &value1, &value1).UpdatedFields(), 0)
		})
		t.Run("Values are different", func(t *testing.T) {
			assert.Len(t, NewFieldsComparator().CompareValues(fieldEmail, &value1, &value2).UpdatedFields(), 1)
			assert.Len(t, NewFieldsComparator().CompareValuesForUpdate(fieldEmail, &value1, &value2).UpdatedFields(), 1)
		})
		t.Run("Case sensitivity", func(t *testing.T) {
			var uppercaseValue1 = strings.ToUpper(value1)
			assert.True(t, NewFieldsComparator().CaseSensitive(true).CompareValues(fieldEmail, &value1, &uppercaseValue1).IsAnyFieldUpdated())
			assert.False(t, NewFieldsComparator().CaseSensitive(false).CompareValues(fieldEmail, &value1, &uppercaseValue1).IsAnyFieldUpdated())
		})
	})
	t.Run("CompareValueAndFunction", func(t *testing.T) {
		var oldValue = "old value"
		var oldValueFunc = func(f Field) []string {
			return []string{oldValue}
		}
		var oldDualValueFunc = func(f Field) []string {
			return []string{oldValue, oldValue}
		}
		t.Run("Value should not be updated", func(t *testing.T) {
			assert.Len(t, NewFieldsComparator().CompareValueAndFunction(fieldEmail, nil, oldValueFunc).UpdatedFields(), 1)
			assert.Len(t, NewFieldsComparator().CompareValueAndFunctionForUpdate(fieldEmail, nil, oldValueFunc).UpdatedFields(), 0)
			assert.Len(t, NewFieldsComparator().CompareValueAndFunctionForUpdate(fieldEmail, &oldValue, oldValueFunc).UpdatedFields(), 0)
		})
		t.Run("Value should be updated", func(t *testing.T) {
			var valueChanged = "new value"
			assert.Len(t, NewFieldsComparator().CompareValueAndFunctionForUpdate(fieldEmail, &valueChanged, oldValueFunc).UpdatedFields(), 1)
			assert.Len(t, NewFieldsComparator().CompareValueAndFunctionForUpdate(fieldEmail, &oldValue, oldDualValueFunc).UpdatedFields(), 1)
		})
	})
	t.Run("CompareValueAndFunctionForUpdate", func(t *testing.T) {
		var oldValue = "old value"
		var emptyFunc = func(f Field) []string {
			return nil
		}
		var oldValueFunc = func(f Field) []string {
			return []string{oldValue}
		}
		var remove = csjson.OptionalString{Defined: true, Value: nil}
		t.Run("Value should not be updated", func(t *testing.T) {
			var unchange = csjson.OptionalString{Defined: false}
			var setSameValue = csjson.OptionalString{Defined: true, Value: &oldValue}
			assert.Len(t, NewFieldsComparator().CompareOptionalAndFunction(fieldEmail, unchange, emptyFunc).UpdatedFields(), 0)
			assert.Len(t, NewFieldsComparator().CompareOptionalAndFunction(fieldEmail, remove, emptyFunc).UpdatedFields(), 0)
			assert.Len(t, NewFieldsComparator().CompareOptionalAndFunction(fieldEmail, setSameValue, oldValueFunc).UpdatedFields(), 0)
		})
		t.Run("Value should be removed", func(t *testing.T) {
			assert.Len(t, NewFieldsComparator().CompareOptionalAndFunction(fieldEmail, remove, oldValueFunc).UpdatedFields(), 1)
		})
		t.Run("Value should be updated", func(t *testing.T) {
			var valueChanged = "new value"
			var change = csjson.OptionalString{Defined: true, Value: &valueChanged}
			assert.Len(t, NewFieldsComparator().CompareOptionalAndFunction(fieldEmail, change, oldValueFunc).UpdatedFields(), 1)
		})
	})
}

func TestUpdateFunctions(t *testing.T) {
	var comparator FieldsComparator = &fieldsComparator{updatedFields: map[Field]bool{Email: true, FirstName: true}}
	t.Run("AreAllFieldsUpdated", func(t *testing.T) {
		assert.False(t, comparator.AreAllFieldsUpdated(BirthDate, BusinessID))
		assert.False(t, comparator.AreAllFieldsUpdated(BirthDate, Email, BusinessID))
		assert.True(t, comparator.AreAllFieldsUpdated(FirstName, Email))
	})
	t.Run("IsAnyFieldUpdated", func(t *testing.T) {
		assert.False(t, comparator.IsAnyFieldUpdated(BirthDate, BusinessID))
		assert.True(t, comparator.IsAnyFieldUpdated(BirthDate, Email, BusinessID))
		assert.True(t, comparator.IsAnyFieldUpdated(FirstName, Email))
	})
	t.Run("UpdatedFields", func(t *testing.T) {
		assert.Len(t, comparator.UpdatedFields(), 2)
		assert.Contains(t, comparator.UpdatedFields(), FirstName.Key())
		assert.Contains(t, comparator.UpdatedFields(), Email.Key())
		assert.NotContains(t, comparator.UpdatedFields(), BirthDate.Key())
	})
}
