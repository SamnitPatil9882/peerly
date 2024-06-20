package api

import (
	"context"
	"fmt"
	"joshsoftware/peerly/apperrors"
	"joshsoftware/peerly/pkg/dto"
	"joshsoftware/peerly/service/Orgnization/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListOrganizationHandler(t *testing.T) {
	orgSvc := mocks.NewService(t)
	orgSvcHandler := listOrganizationHandler(orgSvc)

	tests := []struct {
		name               string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Success",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListOrganizations", mock.Anything).Return([]dto.Organization{
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
						CreatedAt:                time.Date(2024, 6, 17, 11, 9, 19, 716618234, time.UTC),
						CreatedBy:                1,
						UpdatedAt:                time.Date(2024, 6, 17, 11, 9, 19, 716618272, time.UTC),
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":[{"id":1,"name":"TestOrg","email":"test@example.com","domain_name":"example.com","subscription_status":1,"subscription_valid_upto":"2024-06-17T16:00:00Z","hi5_limit":100,"hi5_quota_renewal_frequency":"monthly","timezone":"UTC","created_at":"2024-06-17T11:09:19.716618234Z","created_by":1,"updated_at":"2024-06-17T11:09:19.716618272Z"}]}`,
		},
		// {
		// 	name: "Unauthorized Access",
		// 	setup: func(mockSvc *mocks.Service) {
		// 		// mockSvc.On("ListOrganizations", mock.Anything).Return(nil, fmt.Errorf("service error")).Once()
		// 	},
		// 	expectedStatusCode: http.StatusUnauthorized,
		// 	expectedResponse:   `{"error":"Unauthorized access"}`,
		// },
		{
			name: "Service Error",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListOrganizations", mock.Anything).Return(nil, fmt.Errorf("service error")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"service error", "status":500}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(orgSvc)

			req, err := http.NewRequest("GET", "/organizations", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(orgSvcHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Result().StatusCode)

			// Additional checks for the response body
			assert.JSONEq(t, test.expectedResponse, rr.Body.String())
		})
	}
}

func TestGetOrganizationHandler(t *testing.T) {
	orgSvc := mocks.NewService(t)
	orgSvcHandler := getOrganizationHandler(orgSvc)

	tests := []struct {
		name               string
		role               int
		orgID              int
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:  "Success",
			role:  1,
			orgID: 1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetOrganization", mock.Anything, 1).Return(dto.Organization{
					ID:                       1,
					Name:                     "TestOrg",
					ContactEmail:             "test@example.com",
					DomainName:               "example.com",
					SubscriptionStatus:       1,
					SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
					Hi5Limit:                 100,
					Hi5QuotaRenewalFrequency: "monthly",
					Timezone:                 "UTC",
					CreatedAt:                time.Date(2024, 6, 17, 11, 9, 19, 716618234, time.UTC),
					CreatedBy:                1,
					UpdatedAt:                time.Date(2024, 6, 17, 11, 9, 19, 716618272, time.UTC),
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"id":1,"name":"TestOrg","email":"test@example.com","domain_name":"example.com","subscription_status":1,"subscription_valid_upto":"2024-06-17T16:00:00Z","hi5_limit":100,"hi5_quota_renewal_frequency":"monthly","timezone":"UTC","created_at":"2024-06-17T11:09:19.716618234Z","created_by":1,"updated_at":"2024-06-17T11:09:19.716618272Z"}}`,
		},
		// {
		//     name:               "Unauthorized Access",
		//     role:               2,
		//     orgID:              1,
		//     setup:              func(mockSvc *mocks.Service) {},
		//     expectedStatusCode: http.StatusUnauthorized,
		//     expectedResponse:   `{"error":"Unauthorized access"}`,
		// },
		// {
		//     name:               "Invalid ID",
		//     role:               1,
		//     orgID:              100000000000, // This should trigger invalid ID error
		//     setup:              func(mockSvc *mocks.Service) {},
		//     expectedStatusCode: http.StatusBadRequest,
		//     expectedResponse:   `{"error":"Invalid ID"}`,
		// },
		// {
		//     name:  "Organization Not Found",
		//     role:  1,
		//     orgID: 2,
		//     setup: func(mockSvc *mocks.Service) {
		//         mockSvc.On("GetOrganization", mock.Anything, 2).Return(nil, apperrors.OrganizationNotFound).Once()
		//     },
		//     expectedStatusCode: http.StatusInternalServerError,
		//     expectedResponse:   `{"error":"organization not found"}`,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(orgSvc)

			req, err := http.NewRequest("GET", fmt.Sprintf("/organizations/%d", test.orgID), nil)
			require.NoError(t, err)

			// Adding url params
			req = mux.SetURLVars(req, map[string]string{
				"id": fmt.Sprint(test.orgID),
			})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(orgSvcHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Result().StatusCode)

			// Additional checks for the response body
			assert.JSONEq(t, test.expectedResponse, rr.Body.String())

			// Ensure the mock expectations were met
			orgSvc.AssertExpectations(t)
		})
	}
}

func TestGetOrganizationByDomainNameHandler(t *testing.T) {
	orgSvc := mocks.NewService(t)
	orgSvcHandler := getOrganizationByDomainNameHandler(orgSvc)

	tests := []struct {
		name               string
		role               int
		domainName         string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:       "Success",
			role:       1,
			domainName: "example.com",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetOrganizationByDomainName", mock.Anything, "example.com").Return(dto.Organization{
					ID:                       1,
					Name:                     "TestOrg",
					ContactEmail:             "test@example.com",
					DomainName:               "example.com",
					SubscriptionStatus:       1,
					SubscriptionValidUpto:    time.Date(2024, 6, 17, 16, 0, 0, 0, time.UTC),
					Hi5Limit:                 100,
					Hi5QuotaRenewalFrequency: "monthly",
					Timezone:                 "UTC",
					CreatedAt:                time.Date(2024, 6, 17, 11, 9, 19, 716618234, time.UTC),
					CreatedBy:                1,
					UpdatedAt:                time.Date(2024, 6, 17, 11, 9, 19, 716618272, time.UTC),
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{"data":{"id":1,"name":"TestOrg","email":"test@example.com","domain_name":"example.com",
				"subscription_status":1,"subscription_valid_upto":"2024-06-17T16:00:00Z","hi5_limit":100,
				"hi5_quota_renewal_frequency":"monthly","timezone":"UTC","created_at":"2024-06-17T11:09:19.716618234Z",
				"created_by":1,"updated_at":"2024-06-17T11:09:19.716618272Z"}}`,
		},
		{
			name:       "Organization Not Found",
			role:       1,
			domainName: "nonexistent.com",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetOrganizationByDomainName", mock.Anything, "nonexistent.com").Return(dto.Organization{}, apperrors.OrganizationNotFound).Once()
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `{"message":"organization  not found", "status":404}`,
		},
		// {
		// 	name:               "Unauthorized Access",
		// 	role:               2,
		// 	domainName:         "example.com",
		// 	setup:              func(mockSvc *mocks.Service) {

		// 	},
		// 	expectedStatusCode: http.StatusUnauthorized,
		// 	expectedResponse:   `{"message":"Unauthorised access", "status":401}`,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(orgSvc)

			req, err := http.NewRequest("GET", fmt.Sprintf("/organizations/domain/%s", test.domainName), nil)
			require.NoError(t, err)

			// Adding url params
			req = mux.SetURLVars(req, map[string]string{
				"domainName": fmt.Sprint(test.domainName),
			})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(orgSvcHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Result().StatusCode)

			// Additional checks for the response body
			assert.JSONEq(t, test.expectedResponse, rr.Body.String())

			// Ensure the mock expectations were met
			orgSvc.AssertExpectations(t)
		})
	}
}

func TestCreateOrganizationHandler(t *testing.T) {
	orgSvc := &mocks.Service{}
	orgSvcHandler := createOrganizationHandler(orgSvc)

	tests := []struct {
		name               string
		roleID             int
		requestBody        string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:   "Success",
			roleID: 1,
			requestBody: `{
				"name": "Samnit Patil",
				"email": "samnitpatil9882@gmail.com",
				"domain_name": "samnit.com",
				"subscription_valid_upto": "2024-06-30T23:59:59Z",
				"hi5_limit": 9999,
				"hi5_quota_renewal_frequency": "week",
				"timezone": "UTC"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateOrganization", mock.Anything, mock.AnythingOfType("dto.Organization")).Return(dto.Organization{
					ID:                       1,
					Name:                     "TestOrg",
					ContactEmail:             "test@example.com",
					DomainName:               "example.com",
					SubscriptionStatus:       1,
					SubscriptionValidUpto:    time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
					Hi5Limit:                 100,
					Hi5QuotaRenewalFrequency: "monthly",
					Timezone:                 "UTC",
					CreatedAt:                time.Date(2024, 6, 19, 12, 0, 0, 0, time.UTC),
					CreatedBy:                1,
					UpdatedAt:                time.Date(2024, 6, 19, 12, 0, 0, 0, time.UTC),
				}, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: `{
				"data": {
					"id": 1,
					"name": "TestOrg",
					"email": "test@example.com",
					"domain_name": "example.com",
					"subscription_status": 1,
					"subscription_valid_upto": "2024-12-31T23:59:59Z",
					"hi5_limit": 100,
					"hi5_quota_renewal_frequency": "monthly",
					"timezone": "UTC",
					"created_at": "2024-06-19T12:00:00Z",
					"created_by": 1,
					"updated_at": "2024-06-19T12:00:00Z"
				}
			}`,
		},
		{
			name:               "Invalid JSON",
			roleID:             1,
			requestBody:        `invalid-json`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"error in parsing request in json","status":400}`,
		},
		{
			name:        "Validation Errors",
			roleID:      1,
			requestBody: `{}`,
			setup: func(mockSvc *mocks.Service) {
				// mockSvc.On("CreateOrganization", mock.Anything, mock.AnythingOfType("dto.Organization")).Return(dto.Organization{}, nil).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: `{
    "error": {
        "error": {
            "code": "invalid_data",
            "message": "Please provide valid organization data",
            "fields": {
                "domain_name": "Please enter valid domain",
                "email": "Please enter a valid email",
                "hi5_limit": "Please enter hi5 limit greater than 0",
                "hi5_quota_renewal_frequency": "Please enter valid hi5 renewal frequency",
                "name": "Can't be blank",
                "subscription_valid_upto": "Please enter subscription valid upto date",
                "timezone": "Please enter valid timezone"
            }
        }
    }
}`,
		},
		{
			name:   "Create Organization Error",
			roleID: 1,
			requestBody: `{
				    "name": "Samnit Patil",
    				"email": "samnitpatil9882@gmail.com",
    				"domain_name": "samnit.com",
    				"subscription_status": 1,
    				"subscription_valid_upto": "2024-06-30T23:59:59Z",
    				"hi5_limit": 9999,
    				"hi5_quota_renewal_frequency": "week",
    				"timezone": "UTC"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateOrganization", mock.Anything, mock.AnythingOfType("dto.Organization")).Return(dto.Organization{}, apperrors.InernalServer).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"Failed to create database record", "status":500}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(orgSvc)
			req, err := http.NewRequest("POST", "/organizations", strings.NewReader(test.requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(context.WithValue(req.Context(), "roleid", test.roleID))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(orgSvcHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Result().StatusCode)
			assert.JSONEq(t, test.expectedResponse, rr.Body.String())

			// Ensure the mock expectations were met
			orgSvc.AssertExpectations(t)
		})
	}
}

func TestUpdateOrganizationHandler(t *testing.T) {
	// Define mock organization service
	mockSvc := &mocks.Service{}

	// Create handler function
	handler := updateOrganizationHandler(mockSvc)

	// Define test cases using a table-driven approach
	tests := []struct {
		name               string
		roleID             int
		orgID              string
		requestMethod      string
		requestPath        string
		requestBody        string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:          "Success",
			roleID:        1,
			orgID:         "1",
			requestMethod: "PUT",
			requestPath:   "/organizations/1",
			requestBody: `{
                "id": 1,
                "name": "Updated Organization",
                "email": "updated@example.com",
                "domain_name": "updated.com",
                "subscription_valid_upto": "2024-12-31T23:59:59Z",
                "hi5_limit": 100,
                "hi5_quota_renewal_frequency": "month",
                "timezone": "UTC"
            }`,
			setup: func(mockSvc *mocks.Service) {
				org := dto.Organization{
					ID:                       1,
					Name:                     "Updated Organization",
					ContactEmail:             "updated@example.com",
					DomainName:               "updated.com",
					SubscriptionValidUpto:    time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
					Hi5Limit:                 100,
					Hi5QuotaRenewalFrequency: "month",
					Timezone:                 "UTC",
					CreatedAt:                time.Date(2024, 6, 19, 12, 0, 0, 0, time.UTC),
					CreatedBy:                1,
					UpdatedAt:                time.Date(2024, 6, 19, 12, 0, 0, 0, time.UTC),
				}
				mockSvc.On("UpdateOrganization", mock.Anything, mock.MatchedBy(func(org dto.Organization) bool {
					return org.ID == 1 && org.Name == "Updated Organization"
				})).Return(org, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"data":{"id":1,"name":"Updated Organization","email":"updated@example.com","domain_name":"updated.com","subscription_status":0,"subscription_valid_upto":"2024-12-31T23:59:59Z","hi5_limit":100,"hi5_quota_renewal_frequency":"month","timezone":"UTC","created_at":"2024-06-19T12:00:00Z","created_by":1,"updated_at":"2024-06-19T12:00:00Z"}}`,
		},
		{
			name:          "Invalid ID Parameter",
			roleID:        1,
			orgID:         "ab",
			requestMethod: "PUT",
			requestPath:   "/organizations/ab",
			requestBody: `{
                "id": 1,
                "name": "Updated Organization",
                "email": "updated@example.com",
                "domain_name": "updated.com",
                "subscription_valid_upto": "2024-12-31T23:59:59Z",
                "hi5_limit": 100,
                "hi5_quota_renewal_frequency": "month",
                "timezone": "UTC"
            }`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ``,
		},
		{
			name:               "Invalid JSON",
			roleID:             1,
			orgID:              "1",
			requestMethod:      "PUT",
			requestPath:        "/organizations/1",
			requestBody:        `{invalid-json}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"error in parsing request in json","status":400}`,
		},
		{
			name:          "Validation Errors",
			roleID:        1,
			orgID:         "1",
			requestMethod: "PUT",
			requestPath:   "/organizations/1",
			requestBody: `{
    "email": "invalid-email", 
    "subscription_valid_upto": "2020-01-01T00:00:00Z", 
    "hi5_quota_renewal_frequency": "invalid_frequency" 
}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateOrganization", mock.Anything, mock.AnythingOfType("dto.Organization")).Return(dto.Organization{}, nil).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":{"error":{"code":"invalid_data","message":"Please provide valid organization data","fields":{"email":"Please enter a valid email","hi5_quota_renewal_frequency":"Please enter valid hi5 renewal frequency","subscription_valid_upto":"Please enter subscription valid upto date"}}}}`,
		},
		// {
		// 	name:          "Update Organization Error",
		// 	roleID:        1,
		// 	orgID:         "1",
		// 	requestMethod: "PUT",
		// 	requestPath:   "/organizations/1",
		// 	requestBody: `{
		//         "id": 1,
		//         "name": "Updated Organization",
		//         "email": "updated@example.com",
		//         "domain_name": "updated.com",
		//         "subscription_valid_upto": "2024-12-31T23:59:59Z",
		//         "hi5_limit": 100,
		//         "hi5_quota_renewal_frequency": "month",
		//         "timezone": "UTC"
		//     }`,
		// 	setup: func(mockSvc *mocks.Service) {
		// 		mockSvc.On("UpdateOrganization", mock.Anything, mock.AnythingOfType("dto.Organization")).Return(dto.Organization{}, apperrors.InvalidContactEmail).Once()
		// 	},
		// 	expectedStatusCode: http.StatusInternalServerError,
		// 	expectedResponse: `{
		//         "message": "Internal server error",
		//         "status": 500
		//     }`,
		// },
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Perform setup
			tt.setup(mockSvc)

			// Create request
			reqBody := strings.NewReader(tt.requestBody)
			req, err := http.NewRequest(tt.requestMethod, tt.requestPath, reqBody)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(context.WithValue(req.Context(), "roleid", tt.roleID))
			// Adding url params
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.orgID,
			})
			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve HTTP
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatusCode)
			}

			// Check response body
			if rr.Body.String() != tt.expectedResponse {
				t.Errorf("handler returned unexpected body: got %v want %v, %s",
					rr.Body.String(), tt.expectedResponse, cmp.Diff(rr.Body.String(), tt.expectedResponse))
			}
		})
	}
}
func TestDeleteOrganizationHandler(t *testing.T) {
	// Define mock organization service
	mockSvc := &mocks.Service{}

	// Create handler function
	handler := deleteOrganizationHandler(mockSvc)

	// Define test cases using a table-driven approach
	tests := []struct {
		name               string
		roleID             int
		orgID              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:   "Successful Deletion",
			roleID: 1,
			orgID:  "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteOrganization", mock.Anything, 1, mock.AnythingOfType("int64")).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   ``,
		},
		{
			name:               "Invalid ID Parameter",
			roleID:             1,
			orgID:              "ab",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"Invalid id","status":400}`,
		},
		{
			name:   "Service Error",
			roleID: 1,
			orgID:  "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteOrganization", mock.Anything, 1, mock.AnythingOfType("int64")).Return(apperrors.InernalServer).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"Internal server error","status":500}`,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Perform setup
			tt.setup(mockSvc)

			// Create request
			req, err := http.NewRequest("DELETE", "/organizations/"+tt.orgID, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			req = req.WithContext(context.WithValue(req.Context(), "roleid", tt.roleID))

			// Adding URL params
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.orgID,
			})

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve HTTP
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatusCode)
			}

		})
	}
}

func TestOTPVerificationHandler(t *testing.T) {
	// Define mock organization service
	mockSvc := &mocks.Service{}
	// emailService := &emailMock.MailService{}

	// Create handler function
	handler := OTPVerificationHandler(mockSvc)

	// Define test cases using a table-driven approach
	tests := []struct {
		name               string
		roleID             int
		requestBody        string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Successful Verification",
			requestBody: `{
        "org_id":12,
        "otpcode":"741118"
}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("IsValidContactEmail", mock.Anything, mock.AnythingOfType("dto.OTP")).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "Invalid JSON",
			roleID:      1,
			requestBody: `{invalid-json}`,
			setup:       func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "Validation Errors",
			roleID:      1,
			requestBody: `{"org_id":1,"otpcode":""}`,
			setup: func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "Service Error",
			roleID:      1,
			requestBody: `{"org_id":12,"otpcode":"123456"}`,
			setup: func(mockSvc *mocks.Service) {
				otpInfo := dto.OTP{OrgId: 12, OTPCode: "123456"}
				mockSvc.On("IsValidContactEmail", mock.Anything, otpInfo).Return(apperrors.InvalidOTP).Once()
			},
			expectedStatusCode: http.StatusGone,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Perform setup
			tt.setup(mockSvc)

			// Create request
			reqBody := strings.NewReader(tt.requestBody)
			req, err := http.NewRequest("POST", "/otp/verify", reqBody)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve HTTP
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatusCode)

				fmt.Println("body: ", rr.Body)
			}
		})
	}
}


func TestResendOTPhandler(t *testing.T) {
	// Define mock organization service
	mockSvc := &mocks.Service{}

	// Create handler function
	handler := ResendOTPhandler(mockSvc)

	// Define test cases using a table-driven approach
	tests := []struct {
		name               string
		orgID              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:  "Successful Resend OTP",
			orgID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ResendOTPForContactEmail", mock.Anything, int64(1)).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid ID Parameter",
			orgID:              "invalid",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Service Error",
			orgID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ResendOTPForContactEmail", mock.Anything, int64(1)).Return(apperrors.InernalServer).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Perform setup
			tt.setup(mockSvc)

			// Create request
			req, err := http.NewRequest("POST", "/organizations/"+tt.orgID+"/resend-otp", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			// Adding URL params
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.orgID,
			})

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve HTTP
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatusCode)
			}
		})
	}
}