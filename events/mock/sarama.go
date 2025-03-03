// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/IBM/sarama (interfaces: SyncProducer)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=./mock/sarama.go -package=mock -mock_names=SyncProducer=SyncProducer github.com/IBM/sarama SyncProducer
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	sarama "github.com/IBM/sarama"
	gomock "go.uber.org/mock/gomock"
)

// SyncProducer is a mock of SyncProducer interface.
type SyncProducer struct {
	ctrl     *gomock.Controller
	recorder *SyncProducerMockRecorder
	isgomock struct{}
}

// SyncProducerMockRecorder is the mock recorder for SyncProducer.
type SyncProducerMockRecorder struct {
	mock *SyncProducer
}

// NewSyncProducer creates a new mock instance.
func NewSyncProducer(ctrl *gomock.Controller) *SyncProducer {
	mock := &SyncProducer{ctrl: ctrl}
	mock.recorder = &SyncProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *SyncProducer) EXPECT() *SyncProducerMockRecorder {
	return m.recorder
}

// AbortTxn mocks base method.
func (m *SyncProducer) AbortTxn() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AbortTxn")
	ret0, _ := ret[0].(error)
	return ret0
}

// AbortTxn indicates an expected call of AbortTxn.
func (mr *SyncProducerMockRecorder) AbortTxn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AbortTxn", reflect.TypeOf((*SyncProducer)(nil).AbortTxn))
}

// AddMessageToTxn mocks base method.
func (m *SyncProducer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMessageToTxn", msg, groupId, metadata)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMessageToTxn indicates an expected call of AddMessageToTxn.
func (mr *SyncProducerMockRecorder) AddMessageToTxn(msg, groupId, metadata any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMessageToTxn", reflect.TypeOf((*SyncProducer)(nil).AddMessageToTxn), msg, groupId, metadata)
}

// AddOffsetsToTxn mocks base method.
func (m *SyncProducer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOffsetsToTxn", offsets, groupId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOffsetsToTxn indicates an expected call of AddOffsetsToTxn.
func (mr *SyncProducerMockRecorder) AddOffsetsToTxn(offsets, groupId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOffsetsToTxn", reflect.TypeOf((*SyncProducer)(nil).AddOffsetsToTxn), offsets, groupId)
}

// BeginTxn mocks base method.
func (m *SyncProducer) BeginTxn() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTxn")
	ret0, _ := ret[0].(error)
	return ret0
}

// BeginTxn indicates an expected call of BeginTxn.
func (mr *SyncProducerMockRecorder) BeginTxn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTxn", reflect.TypeOf((*SyncProducer)(nil).BeginTxn))
}

// Close mocks base method.
func (m *SyncProducer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *SyncProducerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*SyncProducer)(nil).Close))
}

// CommitTxn mocks base method.
func (m *SyncProducer) CommitTxn() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitTxn")
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitTxn indicates an expected call of CommitTxn.
func (mr *SyncProducerMockRecorder) CommitTxn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitTxn", reflect.TypeOf((*SyncProducer)(nil).CommitTxn))
}

// IsTransactional mocks base method.
func (m *SyncProducer) IsTransactional() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTransactional")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTransactional indicates an expected call of IsTransactional.
func (mr *SyncProducerMockRecorder) IsTransactional() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTransactional", reflect.TypeOf((*SyncProducer)(nil).IsTransactional))
}

// SendMessage mocks base method.
func (m *SyncProducer) SendMessage(msg *sarama.ProducerMessage) (int32, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", msg)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SendMessage indicates an expected call of SendMessage.
func (mr *SyncProducerMockRecorder) SendMessage(msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*SyncProducer)(nil).SendMessage), msg)
}

// SendMessages mocks base method.
func (m *SyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessages", msgs)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessages indicates an expected call of SendMessages.
func (mr *SyncProducerMockRecorder) SendMessages(msgs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessages", reflect.TypeOf((*SyncProducer)(nil).SendMessages), msgs)
}

// TxnStatus mocks base method.
func (m *SyncProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxnStatus")
	ret0, _ := ret[0].(sarama.ProducerTxnStatusFlag)
	return ret0
}

// TxnStatus indicates an expected call of TxnStatus.
func (mr *SyncProducerMockRecorder) TxnStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxnStatus", reflect.TypeOf((*SyncProducer)(nil).TxnStatus))
}
