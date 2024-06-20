package api

import (
	"fmt"
	"joshsoftware/peerly/config"
	"joshsoftware/peerly/middleware"
	"joshsoftware/peerly/service"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	versionHeader = "Accept"
	authHeader    = "X-Auth-Token"
)

// InitRouter -  The routing mechanism. Mux helps us define handler functions and the access methods
func InitRouter(deps service.Dependencies) (router *mux.Router) {

	// router = mux.NewRouter()
	router = service.InitRouter(deps)
	
	// No version requirement for /ping
	// router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)
	// Version 1 API management
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())

	router.Handle("/organizations", middleware.JwtAuthMiddleware(listOrganizationHandler(deps.OrganizationService), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/organizations/{id:[0-9]+}", middleware.JwtAuthMiddleware(getOrganizationHandler(deps.OrganizationService), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/organizations/{domainName}", middleware.JwtAuthMiddleware(getOrganizationByDomainNameHandler(deps.OrganizationService), deps)).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/organizations", middleware.JwtAuthMiddleware(createOrganizationHandler(deps.OrganizationService), deps)).Methods(http.MethodPost).Headers(versionHeader, v1)

	router.Handle("/organizations/{id:[0-9]+}", middleware.JwtAuthMiddleware(deleteOrganizationHandler(deps.OrganizationService), deps)).Methods(http.MethodDelete).Headers(versionHeader, v1)

	router.Handle("/organizations/{id:[0-9]+}", middleware.JwtAuthMiddleware(updateOrganizationHandler(deps.OrganizationService), deps)).Methods(http.MethodPut).Headers(versionHeader, v1)
	
	router.Handle("/organizations/otp/verify",middleware.JwtAuthMiddleware(OTPVerificationHandler(deps.OrganizationService),deps)).Methods(http.MethodPost).Headers(versionHeader, v1)

	router.Handle("/organizations/otp/{id:[0-9]+}",middleware.JwtAuthMiddleware(ResendOTPhandler(deps.OrganizationService),deps)).Methods(http.MethodPost).Headers(versionHeader, v1)

	return
}