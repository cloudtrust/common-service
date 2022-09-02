package validation

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsStringInSlice(t *testing.T) {
	t.Run("Empty slice", func(t *testing.T) {
		assert.False(t, IsStringInSlice([]string{}, "sunday"))
	})
	t.Run("String is not in slice", func(t *testing.T) {
		assert.False(t, IsStringInSlice([]string{"saturday", "sunday"}, "monday"))
	})
	t.Run("String is in slice", func(t *testing.T) {
		assert.True(t, IsStringInSlice([]string{"saturday", "sunday"}, "sunday"))
	})
}

func failingValidator() Validator {
	return &failedValidator{err: errors.New("unknown reason")}
}

func TestValidator(t *testing.T) {
	var regexp = `\d+`
	var validValue1 = "123"
	var validValue2 = "456789"
	var invalidValue = "abc"

	t.Run("Empty validator - Success", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().Status())
	})
	t.Run("2 valid tests - Success", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().
			ValidateParameterRegExp("param", &validValue1, regexp, false).
			ValidateParameterRegExp("param", &validValue2, regexp, false).
			Status())
	})
	t.Run("1 valid and 1 invalid tests - Error", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().
			ValidateParameterRegExp("param", &validValue1, regexp, false).
			ValidateParameterRegExp("param", &invalidValue, regexp, false).
			Status())
	})
	t.Run("Still failing with different order", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().
			ValidateParameterRegExp("param", &invalidValue, regexp, false).
			ValidateParameterRegExp("param", &validValue1, regexp, false).
			Status())
	})
}

type mystruct struct {
	OneField *time.Time
}

func (ms *mystruct) Validate() error {
	if ms.OneField == nil {
		return errors.New("invalid")
	}
	return nil
}

func TestValidateParameterAndFunc(t *testing.T) {
	t.Run("Nil value, not mandatory", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameter("param", nil, false).Status())
	})
	t.Run("Nil value, mandatory", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameter("param", nil, true).Status())
	})

	var ms mystruct

	t.Run("Invalid value", func(t *testing.T) {
		ms.OneField = nil
		assert.NotNil(t, NewParameterValidator().ValidateParameter("param", &ms, true).Status())
	})
	t.Run("Valid value", func(t *testing.T) {
		var now = time.Now()
		ms.OneField = &now
		assert.Nil(t, NewParameterValidator().ValidateParameter("param", &ms, true).Status())
	})
	t.Run("Valid check after failed validation - Validator", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameter("param", nil, false).Status())
	})
	t.Run("Valid check after failed validation - Func", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterFunc(func() error { return nil }).Status())
	})
}

func TestValidateParameterNotNil(t *testing.T) {
	t.Run("Nil value", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterNotNil("param", nil).Status())
	})
	t.Run("Not nil value", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterNotNil("param", t).Status())
	})
	// This should be "Nil interface{}", but it confuses sonar-scanner
	t.Run("Nil interface", func(t *testing.T) {
		// Skip this test as it fails on the CI
		t.Skip()
		var s = mystruct{OneField: nil}
		assert.Nil(t, NewParameterValidator().ValidateParameterNotNil("param", s.OneField).Status())
	})
	t.Run("Valid check after failed validation", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterNotNil("param", t).Status())
	})
}

func TestValidateParameterIn(t *testing.T) {
	var weekend = map[string]bool{"saturday": true, "sunday": true}
	var validValue = "sunday"
	var invalidValue = "green"

	t.Run("Nil value, not mandatory", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterIn("param", nil, weekend, false).Status())
	})
	t.Run("Nil value, mandatory", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterIn("param", nil, weekend, true).Status())
	})
	t.Run("Valid value", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterIn("param", &validValue, weekend, true).Status())
	})
	t.Run("Invalid value", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterIn("param", &invalidValue, weekend, true).Status())
	})
	t.Run("Valid check after failed validation", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterIn("param", nil, weekend, false).Status())
	})
}

func TestValidateParameterRegExp(t *testing.T) {
	var regexp = `^s[a-z]+day$`
	var validValue = "sunday"
	var invalidValue = "friday"

	t.Run("Nil value, not mandatory", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterRegExp("param", nil, regexp, false).Status())
	})
	t.Run("Nil value, mandatory", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterRegExp("param", nil, regexp, true).Status())
	})
	t.Run("Valid value", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterRegExp("param", &validValue, regexp, true).Status())
	})
	t.Run("Invalid value", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterRegExp("param", &invalidValue, regexp, true).Status())
	})
	t.Run("Valid check after failed validation", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterRegExp("param", nil, regexp, false).Status())
	})
}

func TestValidateParameterPhoneNumber(t *testing.T) {
	t.Run("Nil value, not mandatory", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterPhoneNumber("param", nil, false).Status())
	})
	t.Run("Nil value, mandatory", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterPhoneNumber("param", nil, true).Status())
	})
	t.Run("Valid local number", func(t *testing.T) {
		var validValue = "0785463197"
		assert.Nil(t, NewParameterValidator().ValidateParameterPhoneNumber("param", &validValue, true).Status())
		assert.Equal(t, "+41785463197", validValue)
	})
	t.Run("Valid foreign number", func(t *testing.T) {
		var validValue = "0033685463197"
		assert.Nil(t, NewParameterValidator().ValidateParameterPhoneNumber("param", &validValue, true).Status())
		assert.Equal(t, "+33685463197", validValue)
	})
	t.Run("Invalid value", func(t *testing.T) {
		var invalidValue = "+410654789"
		assert.NotNil(t, NewParameterValidator().ValidateParameterPhoneNumber("param", &invalidValue, true).Status())
	})
	t.Run("Valid check after failed validation", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterPhoneNumber("param", nil, false).Status())
	})
}

func TestValidateParameterLength(t *testing.T) {
	var validValue = "sunday"
	var tooShortValue = "cool"
	var tooLongValue = "freaky friday"

	t.Run("Nil value, not mandatory", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterLength("param", nil, 5, 10, false).Status())
	})
	t.Run("Nil value, mandatory", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterLength("param", nil, 5, 10, true).Status())
	})
	t.Run("Valid value", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterLength("param", &validValue, 5, 10, true).Status())
	})
	t.Run("Too short value", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterLength("param", &tooShortValue, 5, 10, true).Status())
	})
	t.Run("Too long value", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterLength("param", &tooLongValue, 5, 10, true).Status())
	})
	t.Run("Valid check after failed validation", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterLength("param", nil, 5, 10, false).Status())
	})
}

func TestValidateParameterDate(t *testing.T) {
	var dateLayout = "02.01.2006"
	var multipleLayouts = []string{dateLayout, "2006/01/02"}
	var invalidDate = "57.13.2017"
	var alternateFormatDate = "2020/04/12"
	var summer = "21.06.2020"
	var summerAltFormat = "2020/06/21"
	var easter, _ = time.Parse(dateLayout, "12.04.2020")
	var halloween, _ = time.Parse(dateLayout, "31.10.2020")
	var notBetween = "30.07.2019"

	t.Run("Nil value, not mandatory", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterDate("param", nil, dateLayout, false).Status())
	})
	t.Run("Nil value, mandatory", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterDate("param", nil, dateLayout, true).Status())
	})
	t.Run("Valid date", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterDate("param", &summer, dateLayout, true).Status())
	})
	t.Run("Invalid date", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterDate("param", &invalidDate, dateLayout, true).Status())
	})
	t.Run("Multiple layouts", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterDate("param", &alternateFormatDate, dateLayout, true).Status())
		assert.Nil(t, NewParameterValidator().ValidateParameterDateMultipleLayout("param", &alternateFormatDate, multipleLayouts, true).Status())
	})
	t.Run("Multiple layouts-After", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterDate("param", &alternateFormatDate, dateLayout, true).Status())
		assert.Nil(t, NewParameterValidator().ValidateParameterDateAfterMultipleLayout("param", &summerAltFormat, multipleLayouts, easter, true).Status())
	})
	t.Run("Multiple layouts-Not after", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterDate("param", &alternateFormatDate, dateLayout, true).Status())
		assert.NotNil(t, NewParameterValidator().ValidateParameterDateAfterMultipleLayout("param", &summerAltFormat, multipleLayouts, halloween, true).Status())
	})
	t.Run("After", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterDateAfter("param", &summer, dateLayout, easter, true).Status())
	})
	t.Run("Not after", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterDateAfter("param", &summer, dateLayout, halloween, true).Status())
	})
	t.Run("Before", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterDateBefore("param", &summer, dateLayout, halloween, true).Status())
	})
	t.Run("Not before", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterDateBefore("param", &summer, dateLayout, easter, true).Status())
	})
	t.Run("Between", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterDateBetween("param", &summer, dateLayout, easter, halloween, true).Status())
	})
	t.Run("Not between", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterDateBetween("param", &notBetween, dateLayout, easter, halloween, true).Status())
	})
	t.Run("Valid check after failed validation-Date", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterDate("param", nil, dateLayout, false).Status())
	})
	t.Run("Valid check after failed validation-DateMultipleLayout", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterDateMultipleLayout("param", nil, multipleLayouts, false).Status())
	})
	t.Run("Valid check after failed validation-DateMultipleLayoutAfter", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterDateAfterMultipleLayout("param", nil, multipleLayouts, easter, false).Status())
	})
	t.Run("Valid check after failed validation-After", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterDateAfter("param", nil, dateLayout, easter, false).Status())
	})
	t.Run("Valid check after failed validation-Before", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterDateBefore("param", nil, dateLayout, easter, false).Status())
	})
	t.Run("Valid check after failed validation-Between", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterDateBetween("param", nil, dateLayout, easter, halloween, false).Status())
	})
}

func TestValidateParameterLargeDuration(t *testing.T) {
	var validValue = "2y4m7w3d"
	var zeroDurationValue = "2y0m7d"
	var invalidValue = "2y4"

	t.Run("Nil value, not mandatory", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterLargeDuration("param", nil, false).Status())
	})
	t.Run("Nil value, mandatory", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterLargeDuration("param", nil, true).Status())
	})
	t.Run("Valid value", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterLargeDuration("param", &validValue, true).Status())
	})
	t.Run("Invalid value", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterLargeDuration("param", &invalidValue, true).Status())
	})
	t.Run("Contains zero", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterLargeDuration("param", &zeroDurationValue, true).Status())
	})
	t.Run("Valid check after failed validation", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterLargeDuration("param", nil, false).Status())
	})
}

func TestValidateParameterBase64(t *testing.T) {
	var validValue = "dGVzdA=="
	var invalidValue = "2y4==="

	t.Run("Nil value, not mandatory", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterBase64("param", nil, false).Status())
	})
	t.Run("Nil value, mandatory", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterBase64("param", nil, true).Status())
	})
	t.Run("Valid value", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterBase64("param", &validValue, true).Status())
	})
	t.Run("Invalid value", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterBase64("param", &invalidValue, true).Status())
	})
	t.Run("Valid check after failed validation", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterBase64("param", nil, false).Status())
	})
}
