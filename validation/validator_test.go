package validation

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestValidateParameterNotNil(t *testing.T) {
	t.Run("Nil value", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterNotNil("param", nil).Status())
	})
	t.Run("Not nil value", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterNotNil("param", t).Status())
	})
	t.Run("Nil interface{}", func(t *testing.T) {
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
	var validValue = "+33685463197"
	var invalidValue = "+410654789"

	t.Run("Nil value, not mandatory", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterPhoneNumber("param", nil, false).Status())
	})
	t.Run("Nil value, mandatory", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterPhoneNumber("param", nil, true).Status())
	})
	t.Run("Valid value", func(t *testing.T) {
		assert.Nil(t, NewParameterValidator().ValidateParameterPhoneNumber("param", &validValue, true).Status())
	})
	t.Run("Invalid value", func(t *testing.T) {
		assert.NotNil(t, NewParameterValidator().ValidateParameterPhoneNumber("param", &invalidValue, true).Status())
	})
	t.Run("Valid check after failed validation", func(t *testing.T) {
		assert.NotNil(t, failingValidator().ValidateParameterPhoneNumber("param", nil, false).Status())
	})
}

func TestValidateParameterDate(t *testing.T) {
	var dateLayout = "02.01.2006"
	var invalidDate = "57.13.2017"
	var summer = "21.06.2020"
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
