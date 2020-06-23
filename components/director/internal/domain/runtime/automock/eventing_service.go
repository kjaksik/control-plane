// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// EventingService is an autogenerated mock type for the EventingService type
type EventingService struct {
	mock.Mock
}

// GetForRuntime provides a mock function with given fields: ctx, runtimeID
func (_m *EventingService) GetForRuntime(ctx context.Context, runtimeID uuid.UUID) (*model.RuntimeEventingConfiguration, error) {
	ret := _m.Called(ctx, runtimeID)

	var r0 *model.RuntimeEventingConfiguration
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *model.RuntimeEventingConfiguration); ok {
		r0 = rf(ctx, runtimeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RuntimeEventingConfiguration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, runtimeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
