// Code generated by mockery v1.0.0. DO NOT EDIT.

package modelmock

import context "context"
import mock "github.com/stretchr/testify/mock"
import model "github.com/piotrkowalczuk/charon/internal/model"

// RefreshTokenProvider is an autogenerated mock type for the RefreshTokenProvider type
type RefreshTokenProvider struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *RefreshTokenProvider) Create(_a0 context.Context, _a1 *model.RefreshTokenEntity) (*model.RefreshTokenEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.RefreshTokenEntity
	if rf, ok := ret.Get(0).(func(context.Context, *model.RefreshTokenEntity) *model.RefreshTokenEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RefreshTokenEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.RefreshTokenEntity) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: _a0, _a1
func (_m *RefreshTokenProvider) Find(_a0 context.Context, _a1 *model.RefreshTokenFindExpr) ([]*model.RefreshTokenEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*model.RefreshTokenEntity
	if rf, ok := ret.Get(0).(func(context.Context, *model.RefreshTokenFindExpr) []*model.RefreshTokenEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.RefreshTokenEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.RefreshTokenFindExpr) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByToken provides a mock function with given fields: _a0, _a1
func (_m *RefreshTokenProvider) FindOneByToken(_a0 context.Context, _a1 string) (*model.RefreshTokenEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.RefreshTokenEntity
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.RefreshTokenEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RefreshTokenEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByTokenAndUserID provides a mock function with given fields: ctx, token, userID
func (_m *RefreshTokenProvider) FindOneByTokenAndUserID(ctx context.Context, token string, userID int64) (*model.RefreshTokenEntity, error) {
	ret := _m.Called(ctx, token, userID)

	var r0 *model.RefreshTokenEntity
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) *model.RefreshTokenEntity); ok {
		r0 = rf(ctx, token, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RefreshTokenEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int64) error); ok {
		r1 = rf(ctx, token, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOneByToken provides a mock function with given fields: _a0, _a1, _a2
func (_m *RefreshTokenProvider) UpdateOneByToken(_a0 context.Context, _a1 string, _a2 *model.RefreshTokenPatch) (*model.RefreshTokenEntity, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.RefreshTokenEntity
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.RefreshTokenPatch) *model.RefreshTokenEntity); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RefreshTokenEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *model.RefreshTokenPatch) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}