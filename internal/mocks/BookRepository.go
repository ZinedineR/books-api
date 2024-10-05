// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	entity "books-api/internal/entity"
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	model "books-api/internal/model"
)

// BookRepository is an autogenerated mock type for the BookRepository type
type BookRepository struct {
	mock.Mock
}

// CreateTx provides a mock function with given fields: ctx, tx, data
func (_m *BookRepository) CreateTx(ctx context.Context, tx *gorm.DB, data *entity.Book) error {
	ret := _m.Called(ctx, tx, data)

	if len(ret) == 0 {
		panic("no return value specified for CreateTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *entity.Book) error); ok {
		r0 = rf(ctx, tx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByIDTx provides a mock function with given fields: ctx, tx, id
func (_m *BookRepository) DeleteByIDTx(ctx context.Context, tx *gorm.DB, id string) error {
	ret := _m.Called(ctx, tx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteByIDTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, string) error); ok {
		r0 = rf(ctx, tx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByID provides a mock function with given fields: ctx, tx, id
func (_m *BookRepository) FindByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Book, error) {
	ret := _m.Called(ctx, tx, id)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 *entity.Book
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, string) (*entity.Book, error)); ok {
		return rf(ctx, tx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, string) *entity.Book); ok {
		r0 = rf(ctx, tx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Book)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, string) error); ok {
		r1 = rf(ctx, tx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByName provides a mock function with given fields: ctx, tx, column, value
func (_m *BookRepository) FindByName(ctx context.Context, tx *gorm.DB, column string, value string) (*entity.Book, error) {
	ret := _m.Called(ctx, tx, column, value)

	if len(ret) == 0 {
		panic("no return value specified for FindByName")
	}

	var r0 *entity.Book
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, string, string) (*entity.Book, error)); ok {
		return rf(ctx, tx, column, value)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, string, string) *entity.Book); ok {
		r0 = rf(ctx, tx, column, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Book)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, string, string) error); ok {
		r1 = rf(ctx, tx, column, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByPagination provides a mock function with given fields: ctx, tx, page, order, filter
func (_m *BookRepository) FindByPagination(ctx context.Context, tx *gorm.DB, page model.PaginationParam, order model.OrderParam, filter model.FilterParams) (*model.PaginationData[entity.Book], error) {
	ret := _m.Called(ctx, tx, page, order, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindByPagination")
	}

	var r0 *model.PaginationData[entity.Book]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, model.PaginationParam, model.OrderParam, model.FilterParams) (*model.PaginationData[entity.Book], error)); ok {
		return rf(ctx, tx, page, order, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, model.PaginationParam, model.OrderParam, model.FilterParams) *model.PaginationData[entity.Book]); ok {
		r0 = rf(ctx, tx, page, order, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PaginationData[entity.Book])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, model.PaginationParam, model.OrderParam, model.FilterParams) error); ok {
		r1 = rf(ctx, tx, page, order, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTx provides a mock function with given fields: ctx, tx, data
func (_m *BookRepository) UpdateTx(ctx context.Context, tx *gorm.DB, data *entity.Book) error {
	ret := _m.Called(ctx, tx, data)

	if len(ret) == 0 {
		panic("no return value specified for UpdateTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *entity.Book) error); ok {
		r0 = rf(ctx, tx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewBookRepository creates a new instance of BookRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBookRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *BookRepository {
	mock := &BookRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
