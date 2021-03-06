// Code generated by mockery v1.0.0. DO NOT EDIT.

package modelmock

import charon "github.com/piotrkowalczuk/charon"
import context "context"
import mock "github.com/stretchr/testify/mock"
import model "github.com/piotrkowalczuk/charon/internal/model"
import ntypes "github.com/piotrkowalczuk/ntypes"

// GroupProvider is an autogenerated mock type for the GroupProvider type
type GroupProvider struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, createdBy, name, description
func (_m *GroupProvider) Create(ctx context.Context, createdBy int64, name string, description *ntypes.String) (*model.GroupEntity, error) {
	ret := _m.Called(ctx, createdBy, name, description)

	var r0 *model.GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, *ntypes.String) *model.GroupEntity); ok {
		r0 = rf(ctx, createdBy, name, description)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, string, *ntypes.String) error); ok {
		r1 = rf(ctx, createdBy, name, description)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOneByID provides a mock function with given fields: _a0, _a1
func (_m *GroupProvider) DeleteOneByID(_a0 context.Context, _a1 int64) (int64, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: _a0, _a1
func (_m *GroupProvider) Find(_a0 context.Context, _a1 *model.GroupFindExpr) ([]*model.GroupEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*model.GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, *model.GroupFindExpr) []*model.GroupEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.GroupFindExpr) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUserID provides a mock function with given fields: _a0, _a1
func (_m *GroupProvider) FindByUserID(_a0 context.Context, _a1 int64) ([]*model.GroupEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*model.GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64) []*model.GroupEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByID provides a mock function with given fields: _a0, _a1
func (_m *GroupProvider) FindOneByID(_a0 context.Context, _a1 int64) (*model.GroupEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64) *model.GroupEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: _a0, _a1
func (_m *GroupProvider) Insert(_a0 context.Context, _a1 *model.GroupEntity) (*model.GroupEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, *model.GroupEntity) *model.GroupEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.GroupEntity) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsGranted provides a mock function with given fields: _a0, _a1, _a2
func (_m *GroupProvider) IsGranted(_a0 context.Context, _a1 int64, _a2 charon.Permission) (bool, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64, charon.Permission) bool); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, charon.Permission) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPermissions provides a mock function with given fields: _a0, _a1, _a2
func (_m *GroupProvider) SetPermissions(_a0 context.Context, _a1 int64, _a2 ...charon.Permission) (int64, int64, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64, ...charon.Permission) int64); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, int64, ...charon.Permission) int64); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, int64, ...charon.Permission) error); ok {
		r2 = rf(_a0, _a1, _a2...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UpdateOneByID provides a mock function with given fields: _a0, _a1, _a2
func (_m *GroupProvider) UpdateOneByID(_a0 context.Context, _a1 int64, _a2 *model.GroupPatch) (*model.GroupEntity, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64, *model.GroupPatch) *model.GroupEntity); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, *model.GroupPatch) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
