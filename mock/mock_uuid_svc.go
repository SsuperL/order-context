// Code generated by MockGen. DO NOT EDIT.
// Source: acl/ports/clients/uuid_svc.go

// Package mock is a generated GoMock package.
package mock

import (
	pl "order-context/acl/adapters/pl"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUUIDClient is a mock of UUIDClient interface.
type MockUUIDClient struct {
	ctrl     *gomock.Controller
	recorder *MockUUIDClientMockRecorder
}

// MockUUIDClientMockRecorder is the mock recorder for MockUUIDClient.
type MockUUIDClientMockRecorder struct {
	mock *MockUUIDClient
}

// NewMockUUIDClient creates a new mock instance.
func NewMockUUIDClient(ctrl *gomock.Controller) *MockUUIDClient {
	mock := &MockUUIDClient{ctrl: ctrl}
	mock.recorder = &MockUUIDClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUUIDClient) EXPECT() *MockUUIDClientMockRecorder {
	return m.recorder
}

// GetUUID mocks base method.
func (m *MockUUIDClient) GetUUID(arg0 int) (pl.UUIDRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUUID", arg0)
	ret0, _ := ret[0].(pl.UUIDRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUUID indicates an expected call of GetUUID.
func (mr *MockUUIDClientMockRecorder) GetUUID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUUID", reflect.TypeOf((*MockUUIDClient)(nil).GetUUID), arg0)
}
