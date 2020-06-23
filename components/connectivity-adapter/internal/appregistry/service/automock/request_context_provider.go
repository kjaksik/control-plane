// Code generated by mockery v1.1.1. DO NOT EDIT.

package automock

import (
	http "net/http"

	service "github.com/kyma-incubator/compass/components/connectivity-adapter/internal/appregistry/service"
	mock "github.com/stretchr/testify/mock"
)

// RequestContextProvider is an autogenerated mock type for the RequestContextProvider type
type RequestContextProvider struct {
	mock.Mock
}

// ForRequest provides a mock function with given fields: r
func (_m *RequestContextProvider) ForRequest(r *http.Request) (service.RequestContext, error) {
	ret := _m.Called(r)

	var r0 service.RequestContext
	if rf, ok := ret.Get(0).(func(*http.Request) service.RequestContext); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Get(0).(service.RequestContext)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
