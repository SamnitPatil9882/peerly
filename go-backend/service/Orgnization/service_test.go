package orgnization

import (
	"context"
	// "joshsoftware/peerly/config"
	"joshsoftware/peerly/db"
	"joshsoftware/peerly/db/mocks"
	// "joshsoftware/peerly/service/Email"
	"joshsoftware/peerly/pkg/dto"
	"os"
	"testing"
	"time"

	ae "joshsoftware/peerly/apperrors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// func envSetter(envs map[string]string) (closer func()) {
//     originalEnvs := map[string]string{}

//     for name, value := range envs {
//         if originalValue, ok := os.LookupEnv(name); ok {
//             originalEnvs[name] = originalValue
//         }
//         _ = os.Setenv(name, value)
//     }

//     return func() {
//         for name := range envs {
//             origValue, has := originalEnvs[name]
//             if has {
//                 _ = os.Setenv(name, origValue)
//             } else {
//                 _ = os.Unsetenv(name)
//             }
//         }
//     }
// }

func TestListOrganizations(t *testing.T) {
	tests := []struct {
		name              string
		wantOrganizations []dto.Organization
		wantErr           error
		setup             func(OranizationRepo *mocks.OrganizationStorer)
	}{
		{
			name: "Success",
			wantOrganizations: []dto.Organization{
				{
					ID:                       1,
					Name:                     "TestOrg",
					ContactEmail:             "test@example.com",
					DomainName:               "example.com",
					SubscriptionStatus:       1,
					SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
					Hi5Limit:                 100,
					Hi5QuotaRenewalFrequency: "monthly",
					Timezone:                 "UTC",
					CreatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618234, time.Local),
					CreatedBy:                1,
					UpdatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618272, time.Local),
				},
			},
			wantErr: nil,
			setup: func(OranizationRepo *mocks.OrganizationStorer) {
				OranizationRepo.On("ListOrganizations", mock.Anything).Return([]db.Organization{
					{
						ID:                       1,
						Name:                     "TestOrg",
						ContactEmail:             "test@example.com",
						DomainName:               "example.com",
						SubscriptionStatus:       1,
						SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
						Hi5Limit:                 100,
						Hi5QuotaRenewalFrequency: "monthly",
						Timezone:                 "UTC",
						CreatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618234, time.Local),
						CreatedBy:                1,
						UpdatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618272, time.Local),
						SoftDelete:               false,
						SoftDeleteBy:             0,
					},
				}, nil).Once()
			},
		},
		{
			name: "Error from Repository",
			wantOrganizations: []dto.Organization{},
			wantErr:           ae.InernalServer,
			setup: func(OranizationRepo *mocks.OrganizationStorer) {
				OranizationRepo.On("ListOrganizations", mock.Anything).Return([]db.Organization{}, ae.InernalServer).Once()
			},
		},
		{
			name: "Empty Result",
			wantOrganizations: []dto.Organization{},
			wantErr:           nil,
			setup: func(OranizationRepo *mocks.OrganizationStorer) {
				OranizationRepo.On("ListOrganizations", mock.Anything).Return([]db.Organization{}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock repository
			mockOrganizationRepo := new(mocks.OrganizationStorer)
			mockOTPVerificationRepo := new(mocks.OTPVerificationStorer)
			if tt.setup != nil {
				tt.setup(mockOrganizationRepo)
			}

			// Create service with mock repository
			orgSvc := NewService(mockOrganizationRepo,mockOTPVerificationRepo)

			// Call service method
			gotOrganizations, err := orgSvc.ListOrganizations(context.Background())

			// Assert expectations
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantOrganizations, gotOrganizations)

			// Assert that all expected calls were made
			mockOrganizationRepo.AssertExpectations(t)
		})
	}
}


func TestGetOrganization(t *testing.T) {
	tests := []struct {
		name             string
		id               int
		wantOrganization dto.Organization
		wantErr          error
		setup            func(OranizationRepo *mocks.OrganizationStorer)
	}{
		{
			name:             "Success",
			id:               1,
			wantOrganization: dto.Organization{
				ID:                       1,
				Name:                     "TestOrg",
				ContactEmail:             "test@example.com",
				DomainName:               "example.com",
				SubscriptionStatus:       1,
				SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
				Hi5Limit:                 100,
				Hi5QuotaRenewalFrequency: "monthly",
				Timezone:                 "UTC",
				CreatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618234, time.Local),
				CreatedBy:                1,
				UpdatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618272, time.Local),
			},
			wantErr: nil,
			setup: func(OranizationRepo *mocks.OrganizationStorer) {
				OranizationRepo.On("GetOrganization", mock.Anything, 1).Return(db.Organization{
					ID:                       1,
					Name:                     "TestOrg",
					ContactEmail:             "test@example.com",
					DomainName:               "example.com",
					SubscriptionStatus:       1,
					SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
					Hi5Limit:                 100,
					Hi5QuotaRenewalFrequency: "monthly",
					Timezone:                 "UTC",
					CreatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618234, time.Local),
					CreatedBy:                1,
					UpdatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618272, time.Local),
					SoftDelete:               false,
					SoftDeleteBy:             0,
				}, nil).Once()
			},
		},
		{
			name:             "Error from Repository",
			id:               2,
			wantOrganization: dto.Organization{},
			wantErr:          ae.InernalServer,
			setup: func(OranizationRepo *mocks.OrganizationStorer) {
				OranizationRepo.On("GetOrganization", mock.Anything, 2).Return(db.Organization{}, ae.InernalServer).Once()
			},
		},
		{
			name:             "Organization Not Found",
			id:               3,
			wantOrganization: dto.Organization{},
			wantErr:          ae.OrganizationNotFound,
			setup: func(OranizationRepo *mocks.OrganizationStorer) {
				OranizationRepo.On("GetOrganization", mock.Anything, 3).Return(db.Organization{}, ae.OrganizationNotFound).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock repository
			mockOrganizationRepo := new(mocks.OrganizationStorer)
			mockOTPVerificationRepo := new(mocks.OTPVerificationStorer)
			if tt.setup != nil {
				tt.setup(mockOrganizationRepo)
			}

			// Create service with mock repository
			orgSvc := NewService(mockOrganizationRepo, mockOTPVerificationRepo)

			// Call service method
			gotOrganization, err := orgSvc.GetOrganization(context.Background(), tt.id)

			// Assert expectations
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantOrganization, gotOrganization)

			// Assert that all expected calls were made
			mockOrganizationRepo.AssertExpectations(t)
		})
	}
}

func TestGetOrganizationByDomainName(t *testing.T) {
	tests := []struct {
		name             string
		domainName       string
		wantOrganization dto.Organization
		wantErr          error
		setup            func(OranizationRepo *mocks.OrganizationStorer)
	}{
		{
			name:       "Success",
			domainName: "example.com",
			wantOrganization: dto.Organization{
				ID:                       1,
				Name:                     "TestOrg",
				ContactEmail:             "test@example.com",
				DomainName:               "example.com",
				SubscriptionStatus:       1,
				SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
				Hi5Limit:                 100,
				Hi5QuotaRenewalFrequency: "monthly",
				Timezone:                 "UTC",
				CreatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618234, time.Local),
				CreatedBy:                1,
				UpdatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618272, time.Local),
			},
			wantErr: nil,
			setup: func(OranizationRepo *mocks.OrganizationStorer) {
				OranizationRepo.On("GetOrganizationByDomainName", mock.Anything, "example.com").Return(db.Organization{
					ID:                       1,
					Name:                     "TestOrg",
					ContactEmail:             "test@example.com",
					DomainName:               "example.com",
					SubscriptionStatus:       1,
					SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
					Hi5Limit:                 100,
					Hi5QuotaRenewalFrequency: "monthly",
					Timezone:                 "UTC",
					CreatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618234, time.Local),
					CreatedBy:                1,
					UpdatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618272, time.Local),
					SoftDelete:               false,
					SoftDeleteBy:             0,
				}, nil).Once()
			},
		},
		{
			name:             "Error from Repository",
			domainName:       "nonexistent.com",
			wantOrganization: dto.Organization{},
			wantErr:          ae.InernalServer,
			setup: func(OranizationRepo *mocks.OrganizationStorer) {
				OranizationRepo.On("GetOrganizationByDomainName", mock.Anything, "nonexistent.com").Return(db.Organization{}, ae.InernalServer).Once()
			},
		},
		{
			name:             "Organization Not Found",
			domainName:       "notfound.com",
			wantOrganization: dto.Organization{},
			wantErr:          ae.OrganizationNotFound,
			setup: func(OranizationRepo *mocks.OrganizationStorer) {
				OranizationRepo.On("GetOrganizationByDomainName", mock.Anything, "notfound.com").Return(db.Organization{}, ae.OrganizationNotFound).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock repository
			mockOrganizationRepo := new(mocks.OrganizationStorer)
			mockOTPVerificationRepo := new(mocks.OTPVerificationStorer)
			if tt.setup != nil {
				tt.setup(mockOrganizationRepo)
			}

			// Create service with mock repository
			orgSvc := NewService(mockOrganizationRepo, mockOTPVerificationRepo)

			// Call service method
			gotOrganization, err := orgSvc.GetOrganizationByDomainName(context.Background(), tt.domainName)

			// Assert expectations
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantOrganization, gotOrganization)

			// Assert that all expected calls were made
			mockOrganizationRepo.AssertExpectations(t)
		})
	}
}

func TestCreateOrganization(t *testing.T) {
	// Load configurations from application.yml
	// config.LoadConfig()

	tests := []struct {
		name             string
		organization     dto.Organization
		mockSetup        func(organizationRepo *mocks.OrganizationStorer, otpVerificationRepo *mocks.OTPVerificationStorer)
		wantOrganization dto.Organization
		wantErr          error
	}{
		{
			name: "Success",
			organization: dto.Organization{
				Name:                     "TestOrg",
				ContactEmail:             "test@example.com",
				DomainName:               "example.com",
				SubscriptionStatus:       1,
				SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
				Hi5Limit:                 100,
				Hi5QuotaRenewalFrequency: "monthly",
				Timezone:                 "UTC",
				CreatedAt:                time.Now(),
				CreatedBy:                1,
				UpdatedAt:                time.Now(),
			},
			mockSetup: func(organizationRepo *mocks.OrganizationStorer, otpVerificationRepo *mocks.OTPVerificationStorer) {
				organizationRepo.On("IsEmailPresent", mock.Anything, "test@example.com").Return(false).Once()
				organizationRepo.On("IsDomainPresent", mock.Anything, "example.com").Return(false).Once()
				organizationRepo.On("CreateOrganization", mock.Anything, mock.AnythingOfType("dto.Organization")).Return(db.Organization{
					ID:                       1,
					Name:                     "TestOrg",
					ContactEmail:             "test@example.com",
					DomainName:               "example.com",
					SubscriptionStatus:       1,
					SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
					Hi5Limit:                 100,
					Hi5QuotaRenewalFrequency: "monthly",
					Timezone:                 "UTC",
					CreatedAt:                time.Now(),
					CreatedBy:                1,
					UpdatedAt:                time.Now(),
					SoftDelete:               false,
					SoftDeleteBy:             0,
				}, nil).Once()
				otpVerificationRepo.On("CreateOTPInfo", mock.Anything, mock.AnythingOfType("db.OTP")).Return(nil).Once()
			},
			wantOrganization: dto.Organization{
				ID:                       1,
				Name:                     "TestOrg",
				ContactEmail:             "test@example.com",
				DomainName:               "example.com",
				SubscriptionStatus:       1,
				SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
				Hi5Limit:                 100,
				Hi5QuotaRenewalFrequency: "monthly",
				Timezone:                 "UTC",
				CreatedAt:                time.Now(),
				CreatedBy:                1,
				UpdatedAt:                time.Now(),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Initialize mock repositories
			mockOrganizationRepo := new(mocks.OrganizationStorer)
			mockOTPVerificationRepo := new(mocks.OTPVerificationStorer)
			if tt.mockSetup != nil {
				tt.mockSetup(mockOrganizationRepo, mockOTPVerificationRepo)
			}

			// Create service with mock repositories
			orgSvc := NewService(mockOrganizationRepo, mockOTPVerificationRepo)

			// Call service method
			gotOrganization, err := orgSvc.CreateOrganization(context.Background(), tt.organization)

			// Assert expectations
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantOrganization, gotOrganization)

			// Assert that all expected calls were made
			mockOrganizationRepo.AssertExpectations(t)
			mockOTPVerificationRepo.AssertExpectations(t)
		})
	}
}

func envSetter(envs map[string]string) (closer func()) {
	originalEnvs := map[string]string{}

	for name, value := range envs {
		if originalValue, ok := os.LookupEnv(name); ok {
			originalEnvs[name] = originalValue
		}
		_ = os.Setenv(name, value)
	}

	return func() {
		for name := range envs {
			origValue, has := originalEnvs[name]
			if has {
				_ = os.Setenv(name, origValue)
			} else {
				_ = os.Unsetenv(name)
			}
		}
	}
}


func TestUpdateOrganization(t *testing.T) {
	tests := []struct {
		name             string
		organization     dto.Organization
		wantOrganization dto.Organization
		wantErr          error
		setup            func(OrganizationRepo *mocks.OrganizationStorer)
	}{
		{
			name: "Success",
			organization: dto.Organization{
				ID:                       1,
				Name:                     "UpdatedOrg",
				ContactEmail:             "updated@example.com",
				DomainName:               "updated.com",
				SubscriptionStatus:       1,
				SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
				Hi5Limit:                 200,
				Hi5QuotaRenewalFrequency: "yearly",
				Timezone:                 "UTC",
				CreatedAt:                time.Date(2023, time.June, 17, 11, 9, 19, 716618234, time.Local),
				CreatedBy:                1,
				UpdatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618272, time.Local),
			},
			wantOrganization: dto.Organization{
				ID:                       1,
				Name:                     "UpdatedOrg",
				ContactEmail:             "updated@example.com",
				DomainName:               "updated.com",
				SubscriptionStatus:       1,
				SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
				Hi5Limit:                 200,
				Hi5QuotaRenewalFrequency: "yearly",
				Timezone:                 "UTC",
				CreatedAt:                time.Date(2023, time.June, 17, 11, 9, 19, 716618234, time.Local),
				CreatedBy:                1,
				UpdatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618272, time.Local),
			},
			wantErr: nil,
			setup: func(OrganizationRepo *mocks.OrganizationStorer) {
				OrganizationRepo.On("IsOrganizationIdPresent", mock.Anything, int64(1)).Return(true).Once()
				OrganizationRepo.On("IsEmailPresent", mock.Anything, "updated@example.com").Return(false).Once()
				OrganizationRepo.On("IsDomainPresent", mock.Anything, "updated.com").Return(false).Once()
				OrganizationRepo.On("UpdateOrganization", mock.Anything, mock.Anything).Return(db.Organization{
					ID:                       1,
					Name:                     "UpdatedOrg",
					ContactEmail:             "updated@example.com",
					DomainName:               "updated.com",
					SubscriptionStatus:       1,
					SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
					Hi5Limit:                 200,
					Hi5QuotaRenewalFrequency: "yearly",
					Timezone:                 "UTC",
					CreatedAt:                time.Date(2023, time.June, 17, 11, 9, 19, 716618234, time.Local),
					CreatedBy:                1,
					UpdatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618272, time.Local),
					SoftDelete:               false,
					SoftDeleteBy:             0,
				}, nil).Once()
			},
		},
		{
			name: "Organization Not Found",
			organization: dto.Organization{
				ID: 2,
			},
			wantOrganization: dto.Organization{},
			wantErr:          ae.OrganizationNotFound,
			setup: func(OrganizationRepo *mocks.OrganizationStorer) {
				OrganizationRepo.On("IsOrganizationIdPresent", mock.Anything, int64(2)).Return(false).Once()
			},
		},
		{
			name: "Email Already Present",
			organization: dto.Organization{
				ID:           3,
				ContactEmail: "duplicate@example.com",
			},
			wantOrganization: dto.Organization{},
			wantErr:          ae.InvalidContactEmail,
			setup: func(OrganizationRepo *mocks.OrganizationStorer) {
				OrganizationRepo.On("IsOrganizationIdPresent", mock.Anything, int64(3)).Return(true).Once()
				OrganizationRepo.On("IsEmailPresent", mock.Anything, "duplicate@example.com").Return(true).Once()
			},
		},
		{
			name: "Domain Name Already Present",
			organization: dto.Organization{
				ID:         4,
				DomainName: "duplicate.com",
			},
			wantOrganization: dto.Organization{},
			wantErr:          ae.InvalidDomainName,
			setup: func(OrganizationRepo *mocks.OrganizationStorer) {
				OrganizationRepo.On("IsOrganizationIdPresent", mock.Anything, int64(4)).Return(true).Once()
				OrganizationRepo.On("IsEmailPresent", mock.Anything, mock.Anything).Return(false).Once()
				OrganizationRepo.On("IsDomainPresent", mock.Anything, "duplicate.com").Return(true).Once()
			},
		},
		{
			name: "Repository Update Error",
			organization: dto.Organization{
				ID:                       5,
				Name:                     "ErrorOrg",
				ContactEmail:             "error@example.com",
				DomainName:               "error.com",
				SubscriptionStatus:       1,
				SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
				Hi5Limit:                 100,
				Hi5QuotaRenewalFrequency: "monthly",
				Timezone:                 "UTC",
				CreatedAt:                time.Date(2023, time.June, 17, 11, 9, 19, 716618234, time.Local),
				CreatedBy:                1,
				UpdatedAt:                time.Date(2024, time.June, 17, 11, 9, 19, 716618272, time.Local),
			},
			wantOrganization: dto.Organization{},
			wantErr:         ae.InernalServer,
			setup: func(OrganizationRepo *mocks.OrganizationStorer) {
				OrganizationRepo.On("IsOrganizationIdPresent", mock.Anything, int64(5)).Return(true).Once()
				OrganizationRepo.On("IsEmailPresent", mock.Anything, "error@example.com").Return(false).Once()
				OrganizationRepo.On("IsDomainPresent", mock.Anything, "error.com").Return(false).Once()
				OrganizationRepo.On("UpdateOrganization", mock.Anything, mock.Anything).Return(db.Organization{}, ae.InernalServer).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock repository
			mockOrganizationRepo := new(mocks.OrganizationStorer)
			mockOTPVerificationRepo := new(mocks.OTPVerificationStorer)
			if tt.setup != nil {
				tt.setup(mockOrganizationRepo)
			}

			// Create service with mock repository
			orgSvc := NewService(mockOrganizationRepo, mockOTPVerificationRepo)

			// Call service method
			gotOrganization, err := orgSvc.UpdateOrganization(context.Background(), tt.organization)

			// Assert expectations
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantOrganization, gotOrganization)

			// Assert that all expected calls were made
			mockOrganizationRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteOrganization(t *testing.T) {
	tests := []struct {
		name          string
		organizationID int
		userId        int64
		wantErr       error
		setup         func(OrganizationRepo *mocks.OrganizationStorer)
	}{
		{
			name:          "Success",
			organizationID: 1,
			userId:        123,
			wantErr:       nil,
			setup: func(OrganizationRepo *mocks.OrganizationStorer) {
				OrganizationRepo.On("IsOrganizationIdPresent", mock.Anything, int64(1)).Return(true).Once()
				OrganizationRepo.On("DeleteOrganization", mock.Anything, 1, int64(123)).Return(nil).Once()
			},
		},
		{
			name:          "Organization Not Found",
			organizationID: 2,
			userId:        123,
			wantErr:       ae.OrganizationNotFound,
			setup: func(OrganizationRepo *mocks.OrganizationStorer) {
				OrganizationRepo.On("IsOrganizationIdPresent", mock.Anything, int64(2)).Return(false).Once()
			},
		},
		{
			name:          "Repository Delete Error",
			organizationID: 3,
			userId:        123,
			wantErr:       ae.InernalServer,
			setup: func(OrganizationRepo *mocks.OrganizationStorer) {
				OrganizationRepo.On("IsOrganizationIdPresent", mock.Anything, int64(3)).Return(true).Once()
				OrganizationRepo.On("DeleteOrganization", mock.Anything, 3, int64(123)).Return(ae.InernalServer).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock repository
			mockOrganizationRepo := new(mocks.OrganizationStorer)
			mockOTPVerificationRepo := new(mocks.OTPVerificationStorer)
			if tt.setup != nil {
				tt.setup(mockOrganizationRepo)
			}

			// Create service with mock repository
			orgSvc := NewService(mockOrganizationRepo, mockOTPVerificationRepo)

			// Call service method
			err := orgSvc.DeleteOrganization(context.Background(), tt.organizationID, tt.userId)

			// Assert expectations
			assert.Equal(t, tt.wantErr, err)

			// Assert that all expected calls were made
			mockOrganizationRepo.AssertExpectations(t)
		})
	}
}

func TestIsValidContactEmail(t *testing.T) {
	tests := []struct {
		name      string
		otpInfo   dto.OTP
		setup     func(OTPVerificationRepo *mocks.OTPVerificationStorer)
		wantErr   error
	}{
		{
			name: "Success",
			otpInfo: dto.OTP{
				OTPCode: "123456",
				OrgId:   1,
			},
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer) {
				otp := db.OTP{
					OTPCode:   "123456",
					CreatedAt: time.Now().Add(-1 * time.Minute),
				}
				OTPVerificationRepo.On("GetOTPVerificationStatus", mock.Anything, mock.Anything).Return(otp, nil).Once()
				OTPVerificationRepo.On("ChangeIsVerifiedFlag", mock.Anything, int64(1)).Return(nil).Once()
				OTPVerificationRepo.On("DeleteOTPData", mock.Anything, int64(1)).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "Invalid Reference ID",
			otpInfo: dto.OTP{
				OTPCode: "123456",
				OrgId:   2,
			},
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer) {
				OTPVerificationRepo.On("GetOTPVerificationStatus", mock.Anything, mock.Anything).Return(db.OTP{}, ae.InvalidReferenceId).Once()
			},
			wantErr: ae.InvalidReferenceId,
		},
		{
			name: "Expired OTP",
			otpInfo: dto.OTP{
				OTPCode: "123456",
				OrgId:   3,
			},
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer) {
				otp := db.OTP{
					OTPCode:   "123456",
					CreatedAt: time.Now().Add(-3 * time.Minute),
				}
				OTPVerificationRepo.On("GetOTPVerificationStatus", mock.Anything, mock.Anything).Return(otp, nil).Once()
			},
			wantErr: ae.TimeExceeded,
		},
		{
			name: "Invalid OTP",
			otpInfo: dto.OTP{
				OTPCode: "654321",
				OrgId:   4,
			},
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer) {
				otp := db.OTP{
					OTPCode:   "123456",
					CreatedAt: time.Now(),
				}
				OTPVerificationRepo.On("GetOTPVerificationStatus", mock.Anything, mock.Anything).Return(otp, nil).Once()
			},
			wantErr: ae.InvalidOTP,
		},
		{
			name: "Change Is Verified Flag Error",
			otpInfo: dto.OTP{
				OTPCode: "123456",
				OrgId:   5,
			},
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer) {
				otp := db.OTP{
					OTPCode:   "123456",
					CreatedAt: time.Now(),
				}
				OTPVerificationRepo.On("GetOTPVerificationStatus", mock.Anything, mock.Anything).Return(otp, nil).Once()
				OTPVerificationRepo.On("ChangeIsVerifiedFlag", mock.Anything, int64(5)).Return(ae.InernalServer).Once()
			},
			wantErr: ae.InernalServer,
		},
		{
			name: "Delete OTP Data Error",
			otpInfo: dto.OTP{
				OTPCode: "123456",
				OrgId:   6,
			},
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer) {
				otp := db.OTP{
					OTPCode:   "123456",
					CreatedAt: time.Now(),
				}
				OTPVerificationRepo.On("GetOTPVerificationStatus", mock.Anything, mock.Anything).Return(otp, nil).Once()
				OTPVerificationRepo.On("ChangeIsVerifiedFlag", mock.Anything, int64(6)).Return(nil).Once()
				OTPVerificationRepo.On("DeleteOTPData", mock.Anything, int64(6)).Return(ae.InernalServer).Once()
			},
			wantErr: ae.InernalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock repository
			mockOTPVerificationRepo := new(mocks.OTPVerificationStorer)
			mockOrganizationRepo := new(mocks.OrganizationStorer)
			if tt.setup != nil {
				tt.setup(mockOTPVerificationRepo)
			}

			// Create service with mock repository
			orgSvc := NewService(mockOrganizationRepo, mockOTPVerificationRepo)

			// Call service method
			err := orgSvc.IsValidContactEmail(context.Background(), tt.otpInfo)

			// Assert expectations
			assert.Equal(t, tt.wantErr, err)

			// Assert that all expected calls were made
			mockOTPVerificationRepo.AssertExpectations(t)
		})
	}
}

func TestResendOTPForContactEmail(t *testing.T) {
	tests := []struct {
		name      string
		orgId     int64
		setup     func(OTPVerificationRepo *mocks.OTPVerificationStorer, OrganizationRepo *mocks.OrganizationStorer)
		mockSendMail func() error
		wantErr   error
	}{
		{
			name:  "Success",
			orgId: 1,
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer, OrganizationRepo *mocks.OrganizationStorer) {
				OTPVerificationRepo.On("GetCountOfOrgId", mock.Anything, int64(1)).Return(1, nil).Once()
				OrganizationRepo.On("GetOrganization", mock.Anything, 1).Return(db.Organization{
					ID:           1,
					ContactEmail: "test@example.com",
				}, nil).Once()
				OTPVerificationRepo.On("CreateOTPInfo", mock.Anything, mock.Anything).Return(nil).Once()
			},
			mockSendMail: func() error {
				return nil
			},
			wantErr: nil,
		},
		{
			name:  "Organization Not Found",
			orgId: 2,
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer, OrganizationRepo *mocks.OrganizationStorer) {
				OTPVerificationRepo.On("GetCountOfOrgId", mock.Anything, int64(2)).Return(0, nil).Once()
			},
			mockSendMail: nil,
			wantErr: ae.OrganizationNotFound,
		},
		{
			name:  "Attempt Exceeded",
			orgId: 3,
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer, OrganizationRepo *mocks.OrganizationStorer) {
				OTPVerificationRepo.On("GetCountOfOrgId", mock.Anything, int64(3)).Return(3, nil).Once()
			},
			mockSendMail: nil,
			wantErr: ae.AttemptExceeded,
		},
		{
			name:  "GetCountOfOrgId Error",
			orgId: 4,
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer, OrganizationRepo *mocks.OrganizationStorer) {
				OTPVerificationRepo.On("GetCountOfOrgId", mock.Anything, int64(4)).Return(0, ae.InernalServer).Once()
			},
			mockSendMail: nil,
			wantErr: ae.InernalServer,
		},
		{
			name:  "GetOrganization Error",
			orgId: 5,
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer, OrganizationRepo *mocks.OrganizationStorer) {
				OTPVerificationRepo.On("GetCountOfOrgId", mock.Anything, int64(5)).Return(1, nil).Once()
				OrganizationRepo.On("GetOrganization", mock.Anything, 5).Return(db.Organization{}, ae.InernalServer).Once()
			},
			mockSendMail: nil,
			wantErr: ae.InernalServer,
		},
		{
			name:  "CreateOTPInfo Error",
			orgId: 6,
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer, OrganizationRepo *mocks.OrganizationStorer) {
				OTPVerificationRepo.On("GetCountOfOrgId", mock.Anything, int64(6)).Return(1, nil).Once()
				OrganizationRepo.On("GetOrganization", mock.Anything, 6).Return(db.Organization{
					ID:           6,
					ContactEmail: "test@example.com",
				}, nil).Once()
				OTPVerificationRepo.On("CreateOTPInfo", mock.Anything, mock.Anything).Return(ae.InernalServer).Once()
			},
			mockSendMail: nil,
			wantErr: ae.InernalServer,
		},
		{
			name:  "SendMail Error",
			orgId: 7,
			setup: func(OTPVerificationRepo *mocks.OTPVerificationStorer, OrganizationRepo *mocks.OrganizationStorer) {
				OTPVerificationRepo.On("GetCountOfOrgId", mock.Anything, int64(7)).Return(1, nil).Once()
				OrganizationRepo.On("GetOrganization", mock.Anything, 7).Return(db.Organization{
					ID:           7,
					ContactEmail: "test@example.com",
				}, nil).Once()
				OTPVerificationRepo.On("CreateOTPInfo", mock.Anything, mock.Anything).Return(nil).Once()
			},
			mockSendMail: func() error {
				return ae.InernalServer
			},
			wantErr: ae.InernalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock repository
			mockOTPVerificationRepo := new(mocks.OTPVerificationStorer)
			mockOrganizationRepo := new(mocks.OrganizationStorer)
			if tt.setup != nil {
				tt.setup(mockOTPVerificationRepo, mockOrganizationRepo)
			}

			// // Mock SendMail function
			// email.SendMail = tt.mockSendMail

			// Create service with mock repository
			orgSvc := NewService(mockOrganizationRepo, mockOTPVerificationRepo)

			// Call service method
			err := orgSvc.ResendOTPForContactEmail(context.Background(), tt.orgId)

			// Assert expectations
			assert.Equal(t, tt.wantErr, err)

			// Assert that all expected calls were made
			mockOTPVerificationRepo.AssertExpectations(t)
			mockOrganizationRepo.AssertExpectations(t)
		})
	}
}

