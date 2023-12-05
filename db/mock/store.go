// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/naderSameh/billing_system/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	db "github.com/naderSameh/billing_system/db/sqlc"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AddBundleToCustomer mocks base method.
func (m *MockStore) AddBundleToCustomer(arg0 context.Context, arg1 db.AddBundleToCustomerParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBundleToCustomer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBundleToCustomer indicates an expected call of AddBundleToCustomer.
func (mr *MockStoreMockRecorder) AddBundleToCustomer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBundleToCustomer", reflect.TypeOf((*MockStore)(nil).AddBundleToCustomer), arg0, arg1)
}

// AddToDue mocks base method.
func (m *MockStore) AddToDue(arg0 context.Context, arg1 db.AddToDueParams) (db.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToDue", arg0, arg1)
	ret0, _ := ret[0].(db.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddToDue indicates an expected call of AddToDue.
func (mr *MockStoreMockRecorder) AddToDue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToDue", reflect.TypeOf((*MockStore)(nil).AddToDue), arg0, arg1)
}

// AddToPaid mocks base method.
func (m *MockStore) AddToPaid(arg0 context.Context, arg1 db.AddToPaidParams) (db.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToPaid", arg0, arg1)
	ret0, _ := ret[0].(db.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddToPaid indicates an expected call of AddToPaid.
func (mr *MockStoreMockRecorder) AddToPaid(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToPaid", reflect.TypeOf((*MockStore)(nil).AddToPaid), arg0, arg1)
}

// CreateBatch mocks base method.
func (m *MockStore) CreateBatch(arg0 context.Context, arg1 db.CreateBatchParams) (db.Batch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBatch", arg0, arg1)
	ret0, _ := ret[0].(db.Batch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBatch indicates an expected call of CreateBatch.
func (mr *MockStoreMockRecorder) CreateBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBatch", reflect.TypeOf((*MockStore)(nil).CreateBatch), arg0, arg1)
}

// CreateBundle mocks base method.
func (m *MockStore) CreateBundle(arg0 context.Context, arg1 db.CreateBundleParams) (db.Bundle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBundle", arg0, arg1)
	ret0, _ := ret[0].(db.Bundle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBundle indicates an expected call of CreateBundle.
func (mr *MockStoreMockRecorder) CreateBundle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBundle", reflect.TypeOf((*MockStore)(nil).CreateBundle), arg0, arg1)
}

// CreateCustomer mocks base method.
func (m *MockStore) CreateCustomer(arg0 context.Context, arg1 string) (db.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCustomer", arg0, arg1)
	ret0, _ := ret[0].(db.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCustomer indicates an expected call of CreateCustomer.
func (mr *MockStoreMockRecorder) CreateCustomer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCustomer", reflect.TypeOf((*MockStore)(nil).CreateCustomer), arg0, arg1)
}

// CreateOrder mocks base method.
func (m *MockStore) CreateOrder(arg0 context.Context, arg1 db.CreateOrderParams) (db.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", arg0, arg1)
	ret0, _ := ret[0].(db.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockStoreMockRecorder) CreateOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockStore)(nil).CreateOrder), arg0, arg1)
}

// CreatePayment mocks base method.
func (m *MockStore) CreatePayment(arg0 context.Context, arg1 db.CreatePaymentParams) (db.PaymentLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePayment", arg0, arg1)
	ret0, _ := ret[0].(db.PaymentLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePayment indicates an expected call of CreatePayment.
func (mr *MockStoreMockRecorder) CreatePayment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePayment", reflect.TypeOf((*MockStore)(nil).CreatePayment), arg0, arg1)
}

// DeleteBatch mocks base method.
func (m *MockStore) DeleteBatch(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBatch", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBatch indicates an expected call of DeleteBatch.
func (mr *MockStoreMockRecorder) DeleteBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBatch", reflect.TypeOf((*MockStore)(nil).DeleteBatch), arg0, arg1)
}

// DeleteBundle mocks base method.
func (m *MockStore) DeleteBundle(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBundle", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBundle indicates an expected call of DeleteBundle.
func (mr *MockStoreMockRecorder) DeleteBundle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBundle", reflect.TypeOf((*MockStore)(nil).DeleteBundle), arg0, arg1)
}

// DeleteOrder mocks base method.
func (m *MockStore) DeleteOrder(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrder indicates an expected call of DeleteOrder.
func (mr *MockStoreMockRecorder) DeleteOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrder", reflect.TypeOf((*MockStore)(nil).DeleteOrder), arg0, arg1)
}

// DeletePayment mocks base method.
func (m *MockStore) DeletePayment(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePayment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePayment indicates an expected call of DeletePayment.
func (mr *MockStoreMockRecorder) DeletePayment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePayment", reflect.TypeOf((*MockStore)(nil).DeletePayment), arg0, arg1)
}

// GetBatchByName mocks base method.
func (m *MockStore) GetBatchByName(arg0 context.Context, arg1 string) (db.Batch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatchByName", arg0, arg1)
	ret0, _ := ret[0].(db.Batch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatchByName indicates an expected call of GetBatchByName.
func (mr *MockStoreMockRecorder) GetBatchByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatchByName", reflect.TypeOf((*MockStore)(nil).GetBatchByName), arg0, arg1)
}

// GetBatchForUpdate mocks base method.
func (m *MockStore) GetBatchForUpdate(arg0 context.Context, arg1 int64) (db.Batch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatchForUpdate", arg0, arg1)
	ret0, _ := ret[0].(db.Batch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatchForUpdate indicates an expected call of GetBatchForUpdate.
func (mr *MockStoreMockRecorder) GetBatchForUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatchForUpdate", reflect.TypeOf((*MockStore)(nil).GetBatchForUpdate), arg0, arg1)
}

// GetBundleByID mocks base method.
func (m *MockStore) GetBundleByID(arg0 context.Context, arg1 int64) (db.Bundle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBundleByID", arg0, arg1)
	ret0, _ := ret[0].(db.Bundle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBundleByID indicates an expected call of GetBundleByID.
func (mr *MockStoreMockRecorder) GetBundleByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBundleByID", reflect.TypeOf((*MockStore)(nil).GetBundleByID), arg0, arg1)
}

// GetCustomerID mocks base method.
func (m *MockStore) GetCustomerID(arg0 context.Context, arg1 string) (db.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCustomerID", arg0, arg1)
	ret0, _ := ret[0].(db.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCustomerID indicates an expected call of GetCustomerID.
func (mr *MockStoreMockRecorder) GetCustomerID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomerID", reflect.TypeOf((*MockStore)(nil).GetCustomerID), arg0, arg1)
}

// GetOrderByID mocks base method.
func (m *MockStore) GetOrderByID(arg0 context.Context, arg1 int64) (db.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByID", arg0, arg1)
	ret0, _ := ret[0].(db.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByID indicates an expected call of GetOrderByID.
func (mr *MockStoreMockRecorder) GetOrderByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByID", reflect.TypeOf((*MockStore)(nil).GetOrderByID), arg0, arg1)
}

// GetPaymentForUpdate mocks base method.
func (m *MockStore) GetPaymentForUpdate(arg0 context.Context, arg1 int64) (db.PaymentLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPaymentForUpdate", arg0, arg1)
	ret0, _ := ret[0].(db.PaymentLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPaymentForUpdate indicates an expected call of GetPaymentForUpdate.
func (mr *MockStoreMockRecorder) GetPaymentForUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPaymentForUpdate", reflect.TypeOf((*MockStore)(nil).GetPaymentForUpdate), arg0, arg1)
}

// ListAllActiveOrders mocks base method.
func (m *MockStore) ListAllActiveOrders(arg0 context.Context) ([]db.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllActiveOrders", arg0)
	ret0, _ := ret[0].([]db.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllActiveOrders indicates an expected call of ListAllActiveOrders.
func (mr *MockStoreMockRecorder) ListAllActiveOrders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllActiveOrders", reflect.TypeOf((*MockStore)(nil).ListAllActiveOrders), arg0)
}

// ListAllBatches mocks base method.
func (m *MockStore) ListAllBatches(arg0 context.Context, arg1 db.ListAllBatchesParams) ([]db.Batch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllBatches", arg0, arg1)
	ret0, _ := ret[0].([]db.Batch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllBatches indicates an expected call of ListAllBatches.
func (mr *MockStoreMockRecorder) ListAllBatches(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllBatches", reflect.TypeOf((*MockStore)(nil).ListAllBatches), arg0, arg1)
}

// ListAllBundles mocks base method.
func (m *MockStore) ListAllBundles(arg0 context.Context) ([]db.Bundle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllBundles", arg0)
	ret0, _ := ret[0].([]db.Bundle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllBundles indicates an expected call of ListAllBundles.
func (mr *MockStoreMockRecorder) ListAllBundles(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllBundles", reflect.TypeOf((*MockStore)(nil).ListAllBundles), arg0)
}

// ListAllCharges mocks base method.
func (m *MockStore) ListAllCharges(arg0 context.Context, arg1 sql.NullString) ([]db.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllCharges", arg0, arg1)
	ret0, _ := ret[0].([]db.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllCharges indicates an expected call of ListAllCharges.
func (mr *MockStoreMockRecorder) ListAllCharges(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllCharges", reflect.TypeOf((*MockStore)(nil).ListAllCharges), arg0, arg1)
}

// ListBundlesByCustomerID mocks base method.
func (m *MockStore) ListBundlesByCustomerID(arg0 context.Context, arg1 int64) ([]db.Bundle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBundlesByCustomerID", arg0, arg1)
	ret0, _ := ret[0].([]db.Bundle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBundlesByCustomerID indicates an expected call of ListBundlesByCustomerID.
func (mr *MockStoreMockRecorder) ListBundlesByCustomerID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBundlesByCustomerID", reflect.TypeOf((*MockStore)(nil).ListBundlesByCustomerID), arg0, arg1)
}

// ListOrdersByBatchID mocks base method.
func (m *MockStore) ListOrdersByBatchID(arg0 context.Context, arg1 int64) ([]db.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOrdersByBatchID", arg0, arg1)
	ret0, _ := ret[0].([]db.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOrdersByBatchID indicates an expected call of ListOrdersByBatchID.
func (mr *MockStoreMockRecorder) ListOrdersByBatchID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOrdersByBatchID", reflect.TypeOf((*MockStore)(nil).ListOrdersByBatchID), arg0, arg1)
}

// ListOrdersByBundleID mocks base method.
func (m *MockStore) ListOrdersByBundleID(arg0 context.Context, arg1 int64) ([]db.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOrdersByBundleID", arg0, arg1)
	ret0, _ := ret[0].([]db.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOrdersByBundleID indicates an expected call of ListOrdersByBundleID.
func (mr *MockStoreMockRecorder) ListOrdersByBundleID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOrdersByBundleID", reflect.TypeOf((*MockStore)(nil).ListOrdersByBundleID), arg0, arg1)
}

// ListPayments mocks base method.
func (m *MockStore) ListPayments(arg0 context.Context, arg1 db.ListPaymentsParams) ([]db.PaymentLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPayments", arg0, arg1)
	ret0, _ := ret[0].([]db.PaymentLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPayments indicates an expected call of ListPayments.
func (mr *MockStoreMockRecorder) ListPayments(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPayments", reflect.TypeOf((*MockStore)(nil).ListPayments), arg0, arg1)
}

// UpdateBatch mocks base method.
func (m *MockStore) UpdateBatch(arg0 context.Context, arg1 db.UpdateBatchParams) (db.Batch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBatch", arg0, arg1)
	ret0, _ := ret[0].(db.Batch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBatch indicates an expected call of UpdateBatch.
func (mr *MockStoreMockRecorder) UpdateBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBatch", reflect.TypeOf((*MockStore)(nil).UpdateBatch), arg0, arg1)
}

// UpdateOrders mocks base method.
func (m *MockStore) UpdateOrders(arg0 context.Context, arg1 db.UpdateOrdersParams) (db.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrders", arg0, arg1)
	ret0, _ := ret[0].(db.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateOrders indicates an expected call of UpdateOrders.
func (mr *MockStoreMockRecorder) UpdateOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrders", reflect.TypeOf((*MockStore)(nil).UpdateOrders), arg0, arg1)
}

// UpdatePayment mocks base method.
func (m *MockStore) UpdatePayment(arg0 context.Context, arg1 db.UpdatePaymentParams) (db.PaymentLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePayment", arg0, arg1)
	ret0, _ := ret[0].(db.PaymentLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePayment indicates an expected call of UpdatePayment.
func (mr *MockStoreMockRecorder) UpdatePayment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePayment", reflect.TypeOf((*MockStore)(nil).UpdatePayment), arg0, arg1)
}
