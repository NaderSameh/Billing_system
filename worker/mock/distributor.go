// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/naderSameh/billing_system/worker (interfaces: TaskDistributor)

// Package mockwk is a generated GoMock package.
package mockwk

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	asynq "github.com/hibiken/asynq"
	worker "github.com/naderSameh/billing_system/worker"
)

// MockTaskDistributor is a mock of TaskDistributor interface.
type MockTaskDistributor struct {
	ctrl     *gomock.Controller
	recorder *MockTaskDistributorMockRecorder
}

// MockTaskDistributorMockRecorder is the mock recorder for MockTaskDistributor.
type MockTaskDistributorMockRecorder struct {
	mock *MockTaskDistributor
}

// NewMockTaskDistributor creates a new mock instance.
func NewMockTaskDistributor(ctrl *gomock.Controller) *MockTaskDistributor {
	mock := &MockTaskDistributor{ctrl: ctrl}
	mock.recorder = &MockTaskDistributorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskDistributor) EXPECT() *MockTaskDistributorMockRecorder {
	return m.recorder
}

// NewEmailDeliveryTask mocks base method.
func (m *MockTaskDistributor) NewEmailDeliveryTask(arg0 *worker.PayloadSendEmail, arg1 ...asynq.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewEmailDeliveryTask", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewEmailDeliveryTask indicates an expected call of NewEmailDeliveryTask.
func (mr *MockTaskDistributorMockRecorder) NewEmailDeliveryTask(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewEmailDeliveryTask", reflect.TypeOf((*MockTaskDistributor)(nil).NewEmailDeliveryTask), varargs...)
}
