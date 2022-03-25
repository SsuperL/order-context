// Code generated by MockGen. DO NOT EDIT.
// Source: acl/ports/repositories/invoice.go

// Package mock is a generated GoMock package.
package mock

import (
	pl "order-context/acl/adapters/pl"
	aggregate "order-context/domain/aggregate"
	pl0 "order-context/ohs/local/pl"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockInvoiceRepository is a mock of InvoiceRepository interface.
type MockInvoiceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockInvoiceRepositoryMockRecorder
}

// MockInvoiceRepositoryMockRecorder is the mock recorder for MockInvoiceRepository.
type MockInvoiceRepositoryMockRecorder struct {
	mock *MockInvoiceRepository
}

// NewMockInvoiceRepository creates a new mock instance.
func NewMockInvoiceRepository(ctrl *gomock.Controller) *MockInvoiceRepository {
	mock := &MockInvoiceRepository{ctrl: ctrl}
	mock.recorder = &MockInvoiceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInvoiceRepository) EXPECT() *MockInvoiceRepositoryMockRecorder {
	return m.recorder
}

// CheckInvoiceExists mocks base method.
func (m *MockInvoiceRepository) CheckInvoiceExists(invoiceID, siteCode string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckInvoiceExists", invoiceID, siteCode)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckInvoiceExists indicates an expected call of CheckInvoiceExists.
func (mr *MockInvoiceRepositoryMockRecorder) CheckInvoiceExists(invoiceID, siteCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckInvoiceExists", reflect.TypeOf((*MockInvoiceRepository)(nil).CheckInvoiceExists), invoiceID, siteCode)
}

// CreateInvoice mocks base method.
func (m *MockInvoiceRepository) CreateInvoice(arg0 *aggregate.InvoiceAggregate, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateInvoice", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateInvoice indicates an expected call of CreateInvoice.
func (mr *MockInvoiceRepositoryMockRecorder) CreateInvoice(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInvoice", reflect.TypeOf((*MockInvoiceRepository)(nil).CreateInvoice), arg0, arg1)
}

// GetInvoiceDetail mocks base method.
func (m *MockInvoiceRepository) GetInvoiceDetail(invoiceID, siteCode string) (pl.Invoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoiceDetail", invoiceID, siteCode)
	ret0, _ := ret[0].(pl.Invoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvoiceDetail indicates an expected call of GetInvoiceDetail.
func (mr *MockInvoiceRepositoryMockRecorder) GetInvoiceDetail(invoiceID, siteCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoiceDetail", reflect.TypeOf((*MockInvoiceRepository)(nil).GetInvoiceDetail), invoiceID, siteCode)
}

// GetInvoiceList mocks base method.
func (m *MockInvoiceRepository) GetInvoiceList(arg0 pl0.ListInvoiceParams) ([]pl.Invoice, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoiceList", arg0)
	ret0, _ := ret[0].([]pl.Invoice)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetInvoiceList indicates an expected call of GetInvoiceList.
func (mr *MockInvoiceRepositoryMockRecorder) GetInvoiceList(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoiceList", reflect.TypeOf((*MockInvoiceRepository)(nil).GetInvoiceList), arg0)
}

// UpdateInvoice mocks base method.
func (m *MockInvoiceRepository) UpdateInvoice(invoiceID, siteCode string, params pl0.UpdateInvoiceParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInvoice", invoiceID, siteCode, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInvoice indicates an expected call of UpdateInvoice.
func (mr *MockInvoiceRepositoryMockRecorder) UpdateInvoice(invoiceID, siteCode, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInvoice", reflect.TypeOf((*MockInvoiceRepository)(nil).UpdateInvoice), invoiceID, siteCode, params)
}
