// Code generated by mockery v2.20.0. DO NOT EDIT.

package example

import mock "github.com/stretchr/testify/mock"

// MockDB is an autogenerated mock type for the DB type
type MockDB[K comparable, V interface{}] struct {
	mock.Mock
}

type MockDB_Expecter[K comparable, V interface{}] struct {
	mock *mock.Mock
}

func (_m *MockDB[K, V]) EXPECT() *MockDB_Expecter[K, V] {
	return &MockDB_Expecter[K, V]{mock: &_m.Mock}
}

// Get provides a mock function with given fields: k
func (_m *MockDB[K, V]) Get(k K) (V, error) {
	ret := _m.Called(k)

	var r0 V
	var r1 error
	if rf, ok := ret.Get(0).(func(K) (V, error)); ok {
		return rf(k)
	}
	if rf, ok := ret.Get(0).(func(K) V); ok {
		r0 = rf(k)
	} else {
		r0 = ret.Get(0).(V)
	}

	if rf, ok := ret.Get(1).(func(K) error); ok {
		r1 = rf(k)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDB_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockDB_Get_Call[K comparable, V interface{}] struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - k K
func (_e *MockDB_Expecter[K, V]) Get(k interface{}) *MockDB_Get_Call[K, V] {
	return &MockDB_Get_Call[K, V]{Call: _e.mock.On("Get", k)}
}

func (_c *MockDB_Get_Call[K, V]) Run(run func(k K)) *MockDB_Get_Call[K, V] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(K))
	})
	return _c
}

func (_c *MockDB_Get_Call[K, V]) Return(v V, ok error) *MockDB_Get_Call[K, V] {
	_c.Call.Return(v, ok)
	return _c
}

func (_c *MockDB_Get_Call[K, V]) RunAndReturn(run func(K) (V, error)) *MockDB_Get_Call[K, V] {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: k, v
func (_m *MockDB[K, V]) Set(k K, v V) error {
	ret := _m.Called(k, v)

	var r0 error
	if rf, ok := ret.Get(0).(func(K, V) error); ok {
		r0 = rf(k, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDB_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type MockDB_Set_Call[K comparable, V interface{}] struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - k K
//   - v V
func (_e *MockDB_Expecter[K, V]) Set(k interface{}, v interface{}) *MockDB_Set_Call[K, V] {
	return &MockDB_Set_Call[K, V]{Call: _e.mock.On("Set", k, v)}
}

func (_c *MockDB_Set_Call[K, V]) Run(run func(k K, v V)) *MockDB_Set_Call[K, V] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(K), args[1].(V))
	})
	return _c
}

func (_c *MockDB_Set_Call[K, V]) Return(_a0 error) *MockDB_Set_Call[K, V] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDB_Set_Call[K, V]) RunAndReturn(run func(K, V) error) *MockDB_Set_Call[K, V] {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockDB interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDB creates a new instance of MockDB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDB[K comparable, V interface{}](t mockConstructorTestingTNewMockDB) *MockDB[K, V] {
	mock := &MockDB[K, V]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
