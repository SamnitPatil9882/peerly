package repository

import "time"
type OraganizationStorer interface {
	ListOfOrganization() ([]Organization,error)
}

type Organization struct {
	Id int64
	Name string
	ContactEmail string
	DomainName string
	SubscriptionStatus int
	SubscriptionValidUpto time.Time
	Hi5Limit int
	Hi5QuotaRenewalFrequency string
	Timezone string
	CreatedAt time.Time
}