// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import context "context"

import mock "github.com/stretchr/testify/mock"
import model "github.com/kyma-incubator/compass/components/director/internal/model"

// FetchRequestRepository is an autogenerated mock type for the FetchRequestRepository type
type FetchRequestRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, item
func (_m *FetchRequestRepository) Create(ctx context.Context, item *model.FetchRequest) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.FetchRequest) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByReferenceObjectID provides a mock function with given fields: ctx, tenant, objectType, objectID
func (_m *FetchRequestRepository) DeleteByReferenceObjectID(ctx context.Context, tenant string, objectType model.FetchRequestReferenceObjectType, objectID string) error {
	ret := _m.Called(ctx, tenant, objectType, objectID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.FetchRequestReferenceObjectType, string) error); ok {
		r0 = rf(ctx, tenant, objectType, objectID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByReferenceObjectID provides a mock function with given fields: ctx, tenant, objectType, objectID
func (_m *FetchRequestRepository) GetByReferenceObjectID(ctx context.Context, tenant string, objectType model.FetchRequestReferenceObjectType, objectID string) (*model.FetchRequest, error) {
	ret := _m.Called(ctx, tenant, objectType, objectID)

	var r0 *model.FetchRequest
	if rf, ok := ret.Get(0).(func(context.Context, string, model.FetchRequestReferenceObjectType, string) *model.FetchRequest); ok {
		r0 = rf(ctx, tenant, objectType, objectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FetchRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.FetchRequestReferenceObjectType, string) error); ok {
		r1 = rf(ctx, tenant, objectType, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
