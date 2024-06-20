package service

import (
	// "joshsoftware/peerly/aws"
	"joshsoftware/peerly/db"
	orgnization "joshsoftware/peerly/service/Orgnization"

	"github.com/jmoiron/sqlx"
)

// Dependencies - Stuff we need for the service package
type Dependencies struct {
	Store    db.Storer
	// AWSStore aws.AWSStorer
	// define other service dependencies
	OrganizationService orgnization.Service
}

func NewServices(dbInstance *sqlx.DB,store db.Storer) Dependencies {
	orgRepo := db.NewOrganizationRepo(dbInstance)
	otpRepo := db.NewOTPVerificationRepo(dbInstance)
	orgService := orgnization.NewService(orgRepo,otpRepo)
	return Dependencies{
		Store:    store,
		// AWSStore: awsstore,
		OrganizationService: orgService,
	}
}
