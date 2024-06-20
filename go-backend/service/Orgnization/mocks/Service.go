// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	dto "joshsoftware/peerly/pkg/dto"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// CreateOrganization provides a mock function with given fields: ctx, organization
func (_m *Service) CreateOrganization(ctx context.Context, organization dto.Organization) (dto.Organization, error) {
	ret := _m.Called(ctx, organization)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrganization")
	}

	var r0 dto.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.Organization) (dto.Organization, error)); ok {
		return rf(ctx, organization)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dto.Organization) dto.Organization); ok {
		r0 = rf(ctx, organization)
	} else {
		r0 = ret.Get(0).(dto.Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dto.Organization) error); ok {
		r1 = rf(ctx, organization)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOrganization provides a mock function with given fields: ctx, organizationID, userId
func (_m *Service) DeleteOrganization(ctx context.Context, organizationID int, userId int64) error {
	ret := _m.Called(ctx, organizationID, userId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOrganization")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int64) error); ok {
		r0 = rf(ctx, organizationID, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetOrganization provides a mock function with given fields: ctx, id
func (_m *Service) GetOrganization(ctx context.Context, id int) (dto.Organization, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganization")
	}

	var r0 dto.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (dto.Organization, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) dto.Organization); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(dto.Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrganizationByDomainName provides a mock function with given fields: ctx, domainName
func (_m *Service) GetOrganizationByDomainName(ctx context.Context, domainName string) (dto.Organization, error) {
	ret := _m.Called(ctx, domainName)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationByDomainName")
	}

	var r0 dto.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (dto.Organization, error)); ok {
		return rf(ctx, domainName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) dto.Organization); ok {
		r0 = rf(ctx, domainName)
	} else {
		r0 = ret.Get(0).(dto.Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, domainName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsValidContactEmail provides a mock function with given fields: ctx, otpInfo
func (_m *Service) IsValidContactEmail(ctx context.Context, otpInfo dto.OTP) error {
	ret := _m.Called(ctx, otpInfo)

	if len(ret) == 0 {
		panic("no return value specified for IsValidContactEmail")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.OTP) error); ok {
		r0 = rf(ctx, otpInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListOrganizations provides a mock function with given fields: ctx
func (_m *Service) ListOrganizations(ctx context.Context) ([]dto.Organization, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListOrganizations")
	}

	var r0 []dto.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]dto.Organization, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []dto.Organization); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.Organization)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResendOTPForContactEmail provides a mock function with given fields: ctx, orgId
func (_m *Service) ResendOTPForContactEmail(ctx context.Context, orgId int64) error {
	ret := _m.Called(ctx, orgId)

	if len(ret) == 0 {
		panic("no return value specified for ResendOTPForContactEmail")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, orgId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateOrganization provides a mock function with given fields: ctx, organization
func (_m *Service) UpdateOrganization(ctx context.Context, organization dto.Organization) (dto.Organization, error) {
	ret := _m.Called(ctx, organization)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrganization")
	}

	var r0 dto.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.Organization) (dto.Organization, error)); ok {
		return rf(ctx, organization)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dto.Organization) dto.Organization); ok {
		r0 = rf(ctx, organization)
	} else {
		r0 = ret.Get(0).(dto.Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dto.Organization) error); ok {
		r1 = rf(ctx, organization)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}