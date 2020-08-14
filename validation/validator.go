package validation

import (
	"encoding/base64"
	"reflect"
	"regexp"
	"time"

	cerrors "github.com/cloudtrust/common-service/errors"
	"github.com/nyaruka/phonenumbers"
)

// Validatable interface
type Validatable interface {
	Validate() error
}

// Validator interface
type Validator interface {
	ValidateParameter(prmName string, validatable Validatable, mandatory bool) Validator
	ValidateParameterFunc(validator func() error) Validator
	ValidateParameterNotNil(prmName string, value interface{}) Validator
	ValidateParameterIn(prmName string, value *string, allowedValues map[string]bool, mandatory bool) Validator
	ValidateParameterRegExp(prmName string, value *string, regExp string, mandatory bool) Validator
	ValidateParameterPhoneNumber(prmName string, value *string, mandatory bool) Validator
	ValidateParameterLength(prmName string, value *string, min, max int, mandatory bool) Validator
	ValidateParameterIntBetween(prmName string, value *int, min, max int, mandatory bool) Validator
	ValidateParameterDate(prmName string, value *string, dateLayout string, mandatory bool) Validator
	ValidateParameterDateMultipleLayout(prmName string, value *string, dateLayout []string, mandatory bool) Validator
	ValidateParameterDateAfter(prmName string, value *string, dateLayout string, reference time.Time, mandatory bool) Validator
	ValidateParameterDateBefore(prmName string, value *string, dateLayout string, reference time.Time, mandatory bool) Validator
	ValidateParameterDateBetween(prmName string, value *string, dateLayout string, refAfter time.Time, refBetween time.Time, mandatory bool) Validator
	ValidateParameterLargeDuration(prmName string, value *string, mandatory bool) Validator
	ValidateParameterBase64(prmName string, value *string, mandatory bool) Validator
	Status() error
}

type successValidator struct {
}

type failedValidator struct {
	err error
}

// IsStringInSlice tells if a searched value is part of a slice or not
func IsStringInSlice(slice []string, searched string) bool {
	for _, value := range slice {
		if value == searched {
			return true
		}
	}
	return false
}

// NewParameterValidator creates a validator ready to check multiple parameters
func NewParameterValidator() Validator {
	return &successValidator{}
}

func (v *successValidator) ValidateParameter(prmName string, validatable Validatable, mandatory bool) Validator {
	if validatable == nil {
		if mandatory {
			return &failedValidator{err: cerrors.CreateMissingParameterError(prmName)}
		}
	} else {
		return v.ValidateParameterFunc(validatable.Validate)
	}
	return v
}

func (v *successValidator) ValidateParameterFunc(validator func() error) Validator {
	if err := validator(); err != nil {
		return &failedValidator{err: err}
	}
	return v
}

func (v *successValidator) ValidateParameterNotNil(prmName string, value interface{}) Validator {
	if value == nil {
		return &failedValidator{err: cerrors.CreateMissingParameterError(prmName)}
	}
	t := reflect.TypeOf(value)
	if reflect.ValueOf(value) == reflect.Zero(t) {
		return &failedValidator{err: cerrors.CreateMissingParameterError(prmName)}
	}
	return v
}

func (v *successValidator) ValidateParameterIn(prmName string, value *string, allowedValues map[string]bool, mandatory bool) Validator {
	if value == nil {
		if mandatory {
			return &failedValidator{err: cerrors.CreateMissingParameterError(prmName)}
		}
	} else {
		if _, ok := allowedValues[*value]; !ok {
			return &failedValidator{err: cerrors.CreateBadRequestError(cerrors.MsgErrInvalidParam + "." + prmName)}
		}
	}
	return v
}

func (v *successValidator) ValidateParameterRegExp(prmName string, value *string, regExp string, mandatory bool) Validator {
	if value == nil {
		if mandatory {
			return &failedValidator{err: cerrors.CreateMissingParameterError(prmName)}
		}
	} else {
		res, _ := regexp.MatchString(regExp, *value)
		if !res {
			return &failedValidator{err: cerrors.CreateBadRequestError(cerrors.MsgErrInvalidParam + "." + prmName)}
		}
	}
	return v
}

func (v *successValidator) ValidateParameterPhoneNumber(prmName string, value *string, mandatory bool) Validator {
	if value == nil {
		if mandatory {
			return &failedValidator{err: cerrors.CreateMissingParameterError(prmName)}
		}
	} else {
		var metadata, err = phonenumbers.Parse(*value, "CH")
		if err != nil || !phonenumbers.IsPossibleNumber(metadata) {
			return &failedValidator{err: cerrors.CreateBadRequestError(cerrors.MsgErrInvalidParam + "." + prmName)}
		}
	}
	return v
}

func (v *successValidator) ValidateParameterLength(prmName string, value *string, min, max int, mandatory bool) Validator {
	var intValue *int
	if value != nil {
		var length = len(*value)
		intValue = &length
	}
	return v.ValidateParameterIntBetween(prmName, intValue, min, max, mandatory)
}

func (v *successValidator) ValidateParameterIntBetween(prmName string, value *int, min, max int, mandatory bool) Validator {
	if value == nil {
		if mandatory {
			return &failedValidator{err: cerrors.CreateMissingParameterError(prmName)}
		}
	} else {
		if *value < min || *value > max {
			return &failedValidator{err: cerrors.CreateBadRequestError(cerrors.MsgErrInvalidParam + "." + prmName)}
		}
	}
	return v
}

func (v *successValidator) ValidateParameterDate(prmName string, value *string, dateLayout string, mandatory bool) Validator {
	return v.validateParameterDate(prmName, value, []string{dateLayout}, nil, nil, mandatory)
}

func (v *successValidator) ValidateParameterDateMultipleLayout(prmName string, value *string, dateLayout []string, mandatory bool) Validator {
	return v.validateParameterDate(prmName, value, dateLayout, nil, nil, mandatory)
}

func (v *successValidator) ValidateParameterDateAfter(prmName string, value *string, dateLayout string, reference time.Time, mandatory bool) Validator {
	return v.validateParameterDate(prmName, value, []string{dateLayout}, &reference, nil, mandatory)
}

func (v *successValidator) ValidateParameterDateBefore(prmName string, value *string, dateLayout string, reference time.Time, mandatory bool) Validator {
	return v.validateParameterDate(prmName, value, []string{dateLayout}, nil, &reference, mandatory)
}

func (v *successValidator) ValidateParameterDateBetween(prmName string, value *string, dateLayout string, refAfter, refBefore time.Time, mandatory bool) Validator {
	return v.validateParameterDate(prmName, value, []string{dateLayout}, &refAfter, &refBefore, mandatory)
}

func (v *successValidator) validateParameterDate(prmName string, value *string, dateLayout []string, referenceAfter, referenceBefore *time.Time, mandatory bool) Validator {
	if value == nil {
		if mandatory {
			return &failedValidator{err: cerrors.CreateMissingParameterError(prmName)}
		}
	} else {
		var date, err = v.parseDate(dateLayout, *value)
		if err != nil {
			return &failedValidator{err: cerrors.CreateBadRequestError(cerrors.MsgErrInvalidParam + "." + prmName)}
		}
		if referenceAfter != nil && referenceAfter.After(date) {
			return &failedValidator{err: cerrors.CreateBadRequestError(cerrors.MsgErrInvalidParam + "." + prmName)}
		}
		if referenceBefore != nil && referenceBefore.Before(date) {
			return &failedValidator{err: cerrors.CreateBadRequestError(cerrors.MsgErrInvalidParam + "." + prmName)}
		}
	}
	return v
}

func (v *successValidator) ValidateParameterLargeDuration(prmName string, value *string, mandatory bool) Validator {
	if value == nil {
		if mandatory {
			return &failedValidator{err: cerrors.CreateMissingParameterError(prmName)}
		}
	} else {
		if _, err := parseLargeDuration(*value); err != nil {
			return &failedValidator{err: cerrors.CreateBadRequestError(cerrors.MsgErrInvalidParam + "." + prmName)}
		}
	}
	return v
}

func (v *successValidator) ValidateParameterBase64(prmName string, value *string, mandatory bool) Validator {
	if value == nil {
		if mandatory {
			return &failedValidator{err: cerrors.CreateMissingParameterError(prmName)}
		}
	} else {
		if _, err := base64.StdEncoding.DecodeString(*value); err != nil {
			return &failedValidator{err: cerrors.CreateBadRequestError(cerrors.MsgErrInvalidParam + "." + prmName)}
		}
	}

	return v
}

func (v *successValidator) parseDate(dateLayouts []string, value string) (time.Time, error) {
	var resError error
	for _, layout := range dateLayouts {
		var date, err = time.Parse(layout, value)
		if err == nil {
			return date, err
		}
		if resError == nil {
			resError = err
		}
	}
	return time.Time{}, resError
}

func (v *successValidator) Status() error {
	return nil
}

func (v *failedValidator) ValidateParameter(prmName string, validatable Validatable, mandatory bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterFunc(validator func() error) Validator {
	return v
}

func (v *failedValidator) ValidateParameterNotNil(_ string, _ interface{}) Validator {
	return v
}

func (v *failedValidator) ValidateParameterIn(prmName string, value *string, allowedValues map[string]bool, mandatory bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterRegExp(prmName string, value *string, regExp string, mandatory bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterPhoneNumber(prmName string, value *string, mandatory bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterLength(prmName string, value *string, min, max int, mandatory bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterIntBetween(prmName string, value *int, min, max int, mandatory bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterDate(_ string, _ *string, _ string, _ bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterDateMultipleLayout(prmName string, value *string, dateLayout []string, mandatory bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterDateAfter(_ string, _ *string, _ string, _ time.Time, _ bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterDateBefore(_ string, _ *string, _ string, _ time.Time, _ bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterDateBetween(_ string, _ *string, _ string, _, _ time.Time, _ bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterLargeDuration(prmName string, value *string, mandatory bool) Validator {
	return v
}

func (v *failedValidator) ValidateParameterBase64(prmName string, value *string, mandatory bool) Validator {
	return v
}

func (v *failedValidator) Status() error {
	return v.err
}
