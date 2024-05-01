package commonservice

import (
	"fmt"
	"strconv"
	"time"
)

const (
	layout = "02.01.2006"
)

// ToTimePtr converts an 'any' value as a Time pointer. Error is not possible, defaults to nil
func ToTimePtr(value any) *time.Time {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		return &v
	case *time.Time:
		return v
	case string:
		var t, _ = time.Parse(layout, v)
		return &t
	case *string:
		var t, _ = time.Parse(layout, *v)
		return &t
	}
	return nil
}

// ToStringPtr converts an 'any' value as a string pointer
func ToStringPtr(value any) *string {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case *string:
		return v
	}
	var res = fmt.Sprintf("%v", value)
	return &res
}

// ToInt converts an 'any' value as int and returns a default value in case of error
func ToInt(value any, defaultValue int) int {
	switch v := value.(type) {
	case int:
		return v
	case *int:
		return *v
	case int64:
		return int(v)
	case *int64:
		return int(*v)
	case int32:
		return int(v)
	case *int32:
		return int(*v)
	case float32:
		return int(v)
	case *float32:
		return int(*v)
	case float64:
		return int(v)
	case *float64:
		return int(*v)
	case string:
		return atoi(v, defaultValue)
	case *string:
		return atoi(*v, defaultValue)
	default:
		return defaultValue
	}
}

// ToFloat converts an 'any' value as float64 and returns a default value in case of error
func ToFloat(value any, defaultValue float64) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case *int:
		return float64(*v)
	case int64:
		return float64(v)
	case *int64:
		return float64(*v)
	case int32:
		return float64(v)
	case *int32:
		return float64(*v)
	case float32:
		return float64(v)
	case *float32:
		return float64(*v)
	case float64:
		return v
	case *float64:
		return *v
	case string:
		return atof(v, defaultValue)
	case *string:
		return atof(*v, defaultValue)
	default:
		return defaultValue
	}
}

func atoi(value string, defaultValue int) int {
	var res, err = strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return res
}

func atof(value string, defaultValue float64) float64 {
	var res, err = strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return res
}
