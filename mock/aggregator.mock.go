// Code generated by MockGen. DO NOT EDIT.
// Source: aggregator.go
//
// Generated by this command:
//
//	mockgen -source=aggregator.go -destination=../mock/aggregator.mock.go -package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	types "github.com/chenmingyong0423/go-mongox/types"
	gomock "go.uber.org/mock/gomock"
)

// MockiAggregator is a mock of iAggregator interface.
type MockiAggregator[T any] struct {
	ctrl     *gomock.Controller
	recorder *MockiAggregatorMockRecorder[T]
}

// MockiAggregatorMockRecorder is the mock recorder for MockiAggregator.
type MockiAggregatorMockRecorder[T any] struct {
	mock *MockiAggregator[T]
}

// NewMockiAggregator creates a new mock instance.
func NewMockiAggregator[T any](ctrl *gomock.Controller) *MockiAggregator[T] {
	mock := &MockiAggregator[T]{ctrl: ctrl}
	mock.recorder = &MockiAggregatorMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockiAggregator[T]) EXPECT() *MockiAggregatorMockRecorder[T] {
	return m.recorder
}

// Aggregation mocks base method.
func (m *MockiAggregator[T]) Aggregation(ctx context.Context) ([]*T, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Aggregation", ctx)
	ret0, _ := ret[0].([]*T)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Aggregation indicates an expected call of Aggregation.
func (mr *MockiAggregatorMockRecorder[T]) Aggregation(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Aggregation", reflect.TypeOf((*MockiAggregator[T])(nil).Aggregation), ctx)
}

// AggregationWithCallback mocks base method.
func (m *MockiAggregator[T]) AggregationWithCallback(ctx context.Context, handler types.ResultHandler) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AggregationWithCallback", ctx, handler)
	ret0, _ := ret[0].(error)
	return ret0
}

// AggregationWithCallback indicates an expected call of AggregationWithCallback.
func (mr *MockiAggregatorMockRecorder[T]) AggregationWithCallback(ctx, handler any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AggregationWithCallback", reflect.TypeOf((*MockiAggregator[T])(nil).AggregationWithCallback), ctx, handler)
}
