// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	entity "books-api/internal/entity"
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// CreateTx provides a mock function with given fields: ctx, tx, data
func (_m *UserRepository) CreateTx(ctx context.Context, tx *gorm.DB, data *entity.User) error {
	ret := _m.Called(ctx, tx, data)

	if len(ret) == 0 {
		panic("no return value specified for CreateTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *entity.User) error); ok {
		r0 = rf(ctx, tx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByName provides a mock function with given fields: ctx, tx, column, value
func (_m *UserRepository) FindByName(ctx context.Context, tx *gorm.DB, column string, value string) (*entity.User, error) {
	ret := _m.Called(ctx, tx, column, value)

	if len(ret) == 0 {
		panic("no return value specified for FindByName")
	}

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, string, string) (*entity.User, error)); ok {
		return rf(ctx, tx, column, value)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, string, string) *entity.User); ok {
		r0 = rf(ctx, tx, column, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, string, string) error); ok {
		r1 = rf(ctx, tx, column, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
