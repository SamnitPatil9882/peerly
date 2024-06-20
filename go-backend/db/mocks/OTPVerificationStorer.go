// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	db "joshsoftware/peerly/db"

	dto "joshsoftware/peerly/pkg/dto"

	mock "github.com/stretchr/testify/mock"
)

// OTPVerificationStorer is an autogenerated mock type for the OTPVerificationStorer type
type OTPVerificationStorer struct {
	mock.Mock
}

// ChangeIsVerifiedFlag provides a mock function with given fields: ctx, organizationID
func (_m *OTPVerificationStorer) ChangeIsVerifiedFlag(ctx context.Context, organizationID int64) error {
	ret := _m.Called(ctx, organizationID)

	if len(ret) == 0 {
		panic("no return value specified for ChangeIsVerifiedFlag")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, organizationID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateOTPInfo provides a mock function with given fields: ctx, otpinfo
func (_m *OTPVerificationStorer) CreateOTPInfo(ctx context.Context, otpinfo db.OTP) error {
	ret := _m.Called(ctx, otpinfo)

	if len(ret) == 0 {
		panic("no return value specified for CreateOTPInfo")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, db.OTP) error); ok {
		r0 = rf(ctx, otpinfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteOTPData provides a mock function with given fields: ctx, orgId
func (_m *OTPVerificationStorer) DeleteOTPData(ctx context.Context, orgId int64) error {
	ret := _m.Called(ctx, orgId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOTPData")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, orgId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCountOfOrgId provides a mock function with given fields: ctx, orgId
func (_m *OTPVerificationStorer) GetCountOfOrgId(ctx context.Context, orgId int64) (int, error) {
	ret := _m.Called(ctx, orgId)

	if len(ret) == 0 {
		panic("no return value specified for GetCountOfOrgId")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (int, error)); ok {
		return rf(ctx, orgId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) int); ok {
		r0 = rf(ctx, orgId)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, orgId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOTPVerificationStatus provides a mock function with given fields: ctx, otpReq
func (_m *OTPVerificationStorer) GetOTPVerificationStatus(ctx context.Context, otpReq dto.OTP) (db.OTP, error) {
	ret := _m.Called(ctx, otpReq)

	if len(ret) == 0 {
		panic("no return value specified for GetOTPVerificationStatus")
	}

	var r0 db.OTP
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.OTP) (db.OTP, error)); ok {
		return rf(ctx, otpReq)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dto.OTP) db.OTP); ok {
		r0 = rf(ctx, otpReq)
	} else {
		r0 = ret.Get(0).(db.OTP)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dto.OTP) error); ok {
		r1 = rf(ctx, otpReq)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOTPVerificationStorer creates a new instance of OTPVerificationStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOTPVerificationStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *OTPVerificationStorer {
	mock := &OTPVerificationStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
