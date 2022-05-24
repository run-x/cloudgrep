// Code generated by mockery v2.12.2. DO NOT EDIT.

package provider

import (
	mapper "github.com/run-x/cloudgrep/pkg/provider/mapper"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Provider is an autogenerated mock type for the Provider type
type Provider struct {
	mock.Mock
}

type Provider_Expecter struct {
	mock *mock.Mock
}

func (_m *Provider) EXPECT() *Provider_Expecter {
	return &Provider_Expecter{mock: &_m.Mock}
}

// GetMapper provides a mock function with given fields:
func (_m *Provider) GetMapper() mapper.Mapper {
	ret := _m.Called()

	var r0 mapper.Mapper
	if rf, ok := ret.Get(0).(func() mapper.Mapper); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(mapper.Mapper)
	}

	return r0
}

// Provider_GetMapper_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMapper'
type Provider_GetMapper_Call struct {
	*mock.Call
}

// GetMapper is a helper method to define mock.On call
func (_e *Provider_Expecter) GetMapper() *Provider_GetMapper_Call {
	return &Provider_GetMapper_Call{Call: _e.mock.On("GetMapper")}
}

func (_c *Provider_GetMapper_Call) Run(run func()) *Provider_GetMapper_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Provider_GetMapper_Call) Return(_a0 mapper.Mapper) *Provider_GetMapper_Call {
	_c.Call.Return(_a0)
	return _c
}

// Region provides a mock function with given fields:
func (_m *Provider) Region() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Provider_Region_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Region'
type Provider_Region_Call struct {
	*mock.Call
}

// Region is a helper method to define mock.On call
func (_e *Provider_Expecter) Region() *Provider_Region_Call {
	return &Provider_Region_Call{Call: _e.mock.On("Region")}
}

func (_c *Provider_Region_Call) Run(run func()) *Provider_Region_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Provider_Region_Call) Return(_a0 string) *Provider_Region_Call {
	_c.Call.Return(_a0)
	return _c
}

// NewProvider creates a new instance of Provider. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewProvider(t testing.TB) *Provider {
	mock := &Provider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
