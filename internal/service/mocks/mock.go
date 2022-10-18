// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/madhurikadam/app-transcation/internal/domain"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockRepo) CreateAccount(ctx context.Context, account domain.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockRepoMockRecorder) CreateAccount(ctx, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockRepo)(nil).CreateAccount), ctx, account)
}

// CreateCreditTranscation mocks base method.
func (m *MockRepo) CreateCreditTranscation(ctx context.Context, transcation domain.Transcation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCreditTranscation", ctx, transcation)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCreditTranscation indicates an expected call of CreateCreditTranscation.
func (mr *MockRepoMockRecorder) CreateCreditTranscation(ctx, transcation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCreditTranscation", reflect.TypeOf((*MockRepo)(nil).CreateCreditTranscation), ctx, transcation)
}

// CreateDebitTranscation mocks base method.
func (m *MockRepo) CreateDebitTranscation(ctx context.Context, transcation domain.Transcation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDebitTranscation", ctx, transcation)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateDebitTranscation indicates an expected call of CreateDebitTranscation.
func (mr *MockRepoMockRecorder) CreateDebitTranscation(ctx, transcation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDebitTranscation", reflect.TypeOf((*MockRepo)(nil).CreateDebitTranscation), ctx, transcation)
}

// GetAccount mocks base method.
func (m *MockRepo) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, id)
	ret0, _ := ret[0].(*domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockRepoMockRecorder) GetAccount(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockRepo)(nil).GetAccount), ctx, id)
}
