package validation

import (
	"regexp"
	"strconv"
	"time"

	errorhandler "github.com/cloudtrust/common-service/v2/errors"
)

type largeDuration struct {
	years  int
	months int
	days   int
}

const (
	// RegExpLargeDuration is a regular expression used to parse large durations
	// LargeDuration is a duration in days, months, years... Distinguish it from time.Duration which maximum is hours
	RegExpLargeDuration = `([1-9]\d*)([ymwd])`
)

var (
	regExpLargeDuration = regexp.MustCompile(RegExpLargeDuration)
)

func parseLargeDuration(value string) (largeDuration, error) {
	var values = regExpLargeDuration.FindAllStringSubmatch(value, -1)
	var duration largeDuration
	var totalLen = len(value)
	if totalLen == 0 {
		return largeDuration{}, errorhandler.CreateMissingParameterError("duration")
	}
	for _, extract := range values {
		totalLen -= len(extract[0])
		var number, _ = strconv.Atoi(extract[1])
		switch extract[2] {
		case "y":
			duration.years += number
			break
		case "m":
			duration.months += number
			break
		case "w":
			duration.days += 7 * number
			break
		case "d":
			duration.days += number
			break
		}
	}
	if totalLen != 0 {
		return largeDuration{}, errorhandler.CreateBadRequestError(errorhandler.MsgErrInvalidParam + ".duration")
	}
	return duration, nil
}

// IsValidLargeDuration tells whether a duration provided as a string is valid or not
func IsValidLargeDuration(value string) bool {
	var _, err = parseLargeDuration(value)
	if err != nil {
		return false
	}
	return true
}

// AddLargeDuration adds a specified duration value to the given date. In case of any error, the provided date is returned untouched.
func AddLargeDuration(date time.Time, value string) time.Time {
	var duration, _ = AddLargeDurationE(date, value)
	return duration
}

// AddLargeDurationE adds a specified duration value to the given date.
func AddLargeDurationE(date time.Time, value string) (time.Time, error) {
	var duration, err = parseLargeDuration(value)
	if err != nil {
		return date, err
	}
	return date.AddDate(duration.years, duration.months, duration.days), nil
}

// SubstractLargeDuration substracts a specified duration value from the given date. In case of any error, the provided date is returned untouched.
func SubstractLargeDuration(date time.Time, value string) time.Time {
	var duration, _ = SubstractLargeDurationE(date, value)
	return duration
}

// SubstractLargeDurationE substracts a specified duration value from the given date.
func SubstractLargeDurationE(date time.Time, value string) (time.Time, error) {
	var duration, err = parseLargeDuration(value)
	if err != nil {
		return date, err
	}
	return date.AddDate(-duration.years, -duration.months, -duration.days), nil
}
