// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	db "github.com/farismfirdaus/plant-nursery/utils/db"
	mock "github.com/stretchr/testify/mock"
)

// TrxSupportRepo is an autogenerated mock type for the TrxSupportRepo type
type TrxSupportRepo struct {
	mock.Mock
}

// Begin provides a mock function with given fields:
func (_m *TrxSupportRepo) Begin() (db.TrxObj, error) {
	ret := _m.Called()

	var r0 db.TrxObj
	var r1 error
	if rf, ok := ret.Get(0).(func() (db.TrxObj, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() db.TrxObj); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.TrxObj)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTrxSupportRepo creates a new instance of TrxSupportRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTrxSupportRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *TrxSupportRepo {
	mock := &TrxSupportRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}