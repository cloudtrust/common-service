// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cloudtrust/common-service/v2 (interfaces: Configuration)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=./mock/configuration.go -package=mock -mock_names=Configuration=Configuration github.com/cloudtrust/common-service/v2 Configuration
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// Configuration is a mock of Configuration interface.
type Configuration struct {
	ctrl     *gomock.Controller
	recorder *ConfigurationMockRecorder
	isgomock struct{}
}

// ConfigurationMockRecorder is the mock recorder for Configuration.
type ConfigurationMockRecorder struct {
	mock *Configuration
}

// NewConfiguration creates a new mock instance.
func NewConfiguration(ctrl *gomock.Controller) *Configuration {
	mock := &Configuration{ctrl: ctrl}
	mock.recorder = &ConfigurationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Configuration) EXPECT() *ConfigurationMockRecorder {
	return m.recorder
}

// BindEnv mocks base method.
func (m *Configuration) BindEnv(input ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range input {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BindEnv", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// BindEnv indicates an expected call of BindEnv.
func (mr *ConfigurationMockRecorder) BindEnv(input ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BindEnv", reflect.TypeOf((*Configuration)(nil).BindEnv), input...)
}

// Get mocks base method.
func (m *Configuration) Get(key string) any {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(any)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *ConfigurationMockRecorder) Get(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*Configuration)(nil).Get), key)
}

// GetBool mocks base method.
func (m *Configuration) GetBool(key string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBool", key)
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetBool indicates an expected call of GetBool.
func (mr *ConfigurationMockRecorder) GetBool(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBool", reflect.TypeOf((*Configuration)(nil).GetBool), key)
}

// GetDuration mocks base method.
func (m *Configuration) GetDuration(key string) time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDuration", key)
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// GetDuration indicates an expected call of GetDuration.
func (mr *ConfigurationMockRecorder) GetDuration(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDuration", reflect.TypeOf((*Configuration)(nil).GetDuration), key)
}

// GetFloat64 mocks base method.
func (m *Configuration) GetFloat64(key string) float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFloat64", key)
	ret0, _ := ret[0].(float64)
	return ret0
}

// GetFloat64 indicates an expected call of GetFloat64.
func (mr *ConfigurationMockRecorder) GetFloat64(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFloat64", reflect.TypeOf((*Configuration)(nil).GetFloat64), key)
}

// GetInt mocks base method.
func (m *Configuration) GetInt(key string) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt", key)
	ret0, _ := ret[0].(int)
	return ret0
}

// GetInt indicates an expected call of GetInt.
func (mr *ConfigurationMockRecorder) GetInt(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt", reflect.TypeOf((*Configuration)(nil).GetInt), key)
}

// GetInt32 mocks base method.
func (m *Configuration) GetInt32(key string) int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt32", key)
	ret0, _ := ret[0].(int32)
	return ret0
}

// GetInt32 indicates an expected call of GetInt32.
func (mr *ConfigurationMockRecorder) GetInt32(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt32", reflect.TypeOf((*Configuration)(nil).GetInt32), key)
}

// GetInt64 mocks base method.
func (m *Configuration) GetInt64(key string) int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt64", key)
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetInt64 indicates an expected call of GetInt64.
func (mr *ConfigurationMockRecorder) GetInt64(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt64", reflect.TypeOf((*Configuration)(nil).GetInt64), key)
}

// GetString mocks base method.
func (m *Configuration) GetString(key string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", key)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetString indicates an expected call of GetString.
func (mr *ConfigurationMockRecorder) GetString(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*Configuration)(nil).GetString), key)
}

// GetStringSlice mocks base method.
func (m *Configuration) GetStringSlice(key string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStringSlice", key)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetStringSlice indicates an expected call of GetStringSlice.
func (mr *ConfigurationMockRecorder) GetStringSlice(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStringSlice", reflect.TypeOf((*Configuration)(nil).GetStringSlice), key)
}

// GetTime mocks base method.
func (m *Configuration) GetTime(key string) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTime", key)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetTime indicates an expected call of GetTime.
func (mr *ConfigurationMockRecorder) GetTime(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTime", reflect.TypeOf((*Configuration)(nil).GetTime), key)
}

// Set mocks base method.
func (m *Configuration) Set(key string, value any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", key, value)
}

// Set indicates an expected call of Set.
func (mr *ConfigurationMockRecorder) Set(key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*Configuration)(nil).Set), key, value)
}

// SetDefault mocks base method.
func (m *Configuration) SetDefault(key string, value any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDefault", key, value)
}

// SetDefault indicates an expected call of SetDefault.
func (mr *ConfigurationMockRecorder) SetDefault(key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDefault", reflect.TypeOf((*Configuration)(nil).SetDefault), key, value)
}