package orgnization

import (
	"joshsoftware/peerly/db"
	"joshsoftware/peerly/pkg/dto"
)

func OrganizationDBToOrganization(orgDB db.Organization) dto.Organization {
	return dto.Organization{
		ID:                       orgDB.ID,
		Name:                     orgDB.Name,
		ContactEmail:             orgDB.ContactEmail,
		DomainName:               orgDB.DomainName,
		SubscriptionStatus:       orgDB.SubscriptionStatus,
		SubscriptionValidUpto:    orgDB.SubscriptionValidUpto,
		Hi5Limit:                 orgDB.Hi5Limit,
		Hi5QuotaRenewalFrequency: orgDB.Hi5QuotaRenewalFrequency,
		Timezone:                 orgDB.Timezone,
		CreatedAt:                orgDB.CreatedAt,
		CreatedBy:                orgDB.CreatedBy,
		UpdatedAt:                orgDB.UpdatedAt,
		// SoftDelete:               orgDB.SoftDelete,
		// SoftDeleteBy:             orgDB.SoftDeleteBy,
	}
}
