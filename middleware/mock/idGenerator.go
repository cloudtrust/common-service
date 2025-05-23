// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cloudtrust/common-service/v2/idgenerator (interfaces: IDGenerator)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=./mock/idGenerator.go -package=mock -mock_names=IDGenerator=IDGenerator github.com/cloudtrust/common-service/v2/idgenerator IDGenerator
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// IDGenerator is a mock of IDGenerator interface.
type IDGenerator struct {
	ctrl     *gomock.Controller
	recorder *IDGeneratorMockRecorder
	isgomock struct{}
}

// IDGeneratorMockRecorder is the mock recorder for IDGenerator.
type IDGeneratorMockRecorder struct {
	mock *IDGenerator
}

// NewIDGenerator creates a new mock instance.
func NewIDGenerator(ctrl *gomock.Controller) *IDGenerator {
	mock := &IDGenerator{ctrl: ctrl}
	mock.recorder = &IDGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *IDGenerator) EXPECT() *IDGeneratorMockRecorder {
	return m.recorder
}

// NextID mocks base method.
func (m *IDGenerator) NextID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextID")
	ret0, _ := ret[0].(string)
	return ret0
}

// NextID indicates an expected call of NextID.
func (mr *IDGeneratorMockRecorder) NextID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextID", reflect.TypeOf((*IDGenerator)(nil).NextID))
}
