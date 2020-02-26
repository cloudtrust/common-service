package validation

import (
	"reflect"
	"regexp"
	"time"

	cerrors "github.com/cloudtrust/common-service/errors"
	"github.com/nyaruka/phonenumbers"
)

// Validator interface
type Validator interface {
	ValidateParameterNotNil(prmName string, value interface{}) Validator
	ValidateParameterIn(prmName string, value *string, allowedValues map[string]bool, mandatory bool) Validator
	ValidateParameterRegExp(prmName string, value *string, regExp string, mandatory bool) Validator
	ValidateParameterPhoneNumber(prmName string, value *string, mandatory bool) Validator
	ValidateParameterDate(prmName string, value *string, dateLayout string, mandatory bool) Validator
	ValidateParameterDateMultipleLayout(prmName string, value *string, dateLayout []string, mandatory bool) Validator
	ValidateParameterDateAfter(prmName string, value *string, dateLayout string, reference time.Time, mandatory bool) Validator
	ValidateParameterDateBefore(prmName string, value *string, dateLayout string, reference time.Time, mandatory bool) Validator
	ValidateParameterDateBetween(prmName string, value *string, dateLayout string, refAfter time.Time, refBetween time.Time, mandatory bool) Validator
	Status() error
}

type successValidator struct {
}

type failedValidator struct {
	err error
}

// NewParameterValidator creates a validator ready to check multiple parameters
func NewParameterValidator() Validator {
	return &successValidator{}
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

func (v *failedValidator) Status() error {
	return v.err
}
