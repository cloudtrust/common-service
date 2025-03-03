// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cloudtrust/common-service/v2/database/sqltypes (interfaces: CloudtrustDB,SQLRow,SQLRows)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=./mock/cloudtrustdb.go -package=mock -mock_names=CloudtrustDB=CloudtrustDB,SQLRow=SQLRow,SQLRows=SQLRows github.com/cloudtrust/common-service/v2/database/sqltypes CloudtrustDB,SQLRow,SQLRows
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	sqltypes "github.com/cloudtrust/common-service/v2/database/sqltypes"
	gomock "go.uber.org/mock/gomock"
)

// CloudtrustDB is a mock of CloudtrustDB interface.
type CloudtrustDB struct {
	ctrl     *gomock.Controller
	recorder *CloudtrustDBMockRecorder
	isgomock struct{}
}

// CloudtrustDBMockRecorder is the mock recorder for CloudtrustDB.
type CloudtrustDBMockRecorder struct {
	mock *CloudtrustDB
}

// NewCloudtrustDB creates a new mock instance.
func NewCloudtrustDB(ctrl *gomock.Controller) *CloudtrustDB {
	mock := &CloudtrustDB{ctrl: ctrl}
	mock.recorder = &CloudtrustDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *CloudtrustDB) EXPECT() *CloudtrustDBMockRecorder {
	return m.recorder
}

// BeginTx mocks base method.
func (m *CloudtrustDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (sqltypes.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx, opts)
	ret0, _ := ret[0].(sqltypes.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *CloudtrustDBMockRecorder) BeginTx(ctx, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*CloudtrustDB)(nil).BeginTx), ctx, opts)
}

// Close mocks base method.
func (m *CloudtrustDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *CloudtrustDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*CloudtrustDB)(nil).Close))
}

// Exec mocks base method.
func (m *CloudtrustDB) Exec(query string, args ...any) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *CloudtrustDBMockRecorder) Exec(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*CloudtrustDB)(nil).Exec), varargs...)
}

// Ping mocks base method.
func (m *CloudtrustDB) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *CloudtrustDBMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*CloudtrustDB)(nil).Ping))
}

// Query mocks base method.
func (m *CloudtrustDB) Query(query string, args ...any) (sqltypes.SQLRows, error) {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(sqltypes.SQLRows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *CloudtrustDBMockRecorder) Query(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*CloudtrustDB)(nil).Query), varargs...)
}

// QueryRow mocks base method.
func (m *CloudtrustDB) QueryRow(query string, args ...any) sqltypes.SQLRow {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(sqltypes.SQLRow)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *CloudtrustDBMockRecorder) QueryRow(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*CloudtrustDB)(nil).QueryRow), varargs...)
}

// Stats mocks base method.
func (m *CloudtrustDB) Stats() sql.DBStats {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats")
	ret0, _ := ret[0].(sql.DBStats)
	return ret0
}

// Stats indicates an expected call of Stats.
func (mr *CloudtrustDBMockRecorder) Stats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*CloudtrustDB)(nil).Stats))
}

// SQLRow is a mock of SQLRow interface.
type SQLRow struct {
	ctrl     *gomock.Controller
	recorder *SQLRowMockRecorder
	isgomock struct{}
}

// SQLRowMockRecorder is the mock recorder for SQLRow.
type SQLRowMockRecorder struct {
	mock *SQLRow
}

// NewSQLRow creates a new mock instance.
func NewSQLRow(ctrl *gomock.Controller) *SQLRow {
	mock := &SQLRow{ctrl: ctrl}
	mock.recorder = &SQLRowMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *SQLRow) EXPECT() *SQLRowMockRecorder {
	return m.recorder
}

// Scan mocks base method.
func (m *SQLRow) Scan(dest ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range dest {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Scan", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *SQLRowMockRecorder) Scan(dest ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*SQLRow)(nil).Scan), dest...)
}

// SQLRows is a mock of SQLRows interface.
type SQLRows struct {
	ctrl     *gomock.Controller
	recorder *SQLRowsMockRecorder
	isgomock struct{}
}

// SQLRowsMockRecorder is the mock recorder for SQLRows.
type SQLRowsMockRecorder struct {
	mock *SQLRows
}

// NewSQLRows creates a new mock instance.
func NewSQLRows(ctrl *gomock.Controller) *SQLRows {
	mock := &SQLRows{ctrl: ctrl}
	mock.recorder = &SQLRowsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *SQLRows) EXPECT() *SQLRowsMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *SQLRows) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *SQLRowsMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*SQLRows)(nil).Close))
}

// Err mocks base method.
func (m *SQLRows) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *SQLRowsMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*SQLRows)(nil).Err))
}

// Next mocks base method.
func (m *SQLRows) Next() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Next indicates an expected call of Next.
func (mr *SQLRowsMockRecorder) Next() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*SQLRows)(nil).Next))
}

// NextResultSet mocks base method.
func (m *SQLRows) NextResultSet() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextResultSet")
	ret0, _ := ret[0].(bool)
	return ret0
}

// NextResultSet indicates an expected call of NextResultSet.
func (mr *SQLRowsMockRecorder) NextResultSet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextResultSet", reflect.TypeOf((*SQLRows)(nil).NextResultSet))
}

// Scan mocks base method.
func (m *SQLRows) Scan(dest ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range dest {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Scan", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *SQLRowsMockRecorder) Scan(dest ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*SQLRows)(nil).Scan), dest...)
}
