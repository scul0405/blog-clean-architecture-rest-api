// Code generated by MockGen. DO NOT EDIT.
// Source: pg_repo.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	models "github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	utils "github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepository) Create(ctx context.Context, blog *models.Blog) (*models.BlogBase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, blog)
	ret0, _ := ret[0].(*models.BlogBase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(ctx, blog interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), ctx, blog)
}

// Delete mocks base method.
func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), ctx, id)
}

// GetByID mocks base method.
func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.BlogBase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.BlogBase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), ctx, id)
}

// List mocks base method.
func (m *MockRepository) List(ctx context.Context, pq *utils.PaginationQuery) (*models.BlogsList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, pq)
	ret0, _ := ret[0].(*models.BlogsList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRepositoryMockRecorder) List(ctx, pq interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepository)(nil).List), ctx, pq)
}

// Update mocks base method.
func (m *MockRepository) Update(ctx context.Context, blog *models.BlogBase) (*models.BlogBase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, blog)
	ret0, _ := ret[0].(*models.BlogBase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(ctx, blog interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), ctx, blog)
}
