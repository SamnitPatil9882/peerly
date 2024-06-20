package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	ae "joshsoftware/peerly/apperrors"
	"joshsoftware/peerly/pkg/dto"
	"joshsoftware/peerly/util/log"
	"strconv"
	"strings"
	"time"
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
)

const (
	createOrganizationQuery = `INSERT INTO organizations (
		name,
		contact_email,
		domain_name,
		subscription_status,
		subscription_valid_upto,
		hi5_limit,
		hi5_quota_renewal_frequency,
		timezone,
		created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	updateOrganizationQuery = `UPDATE organizations SET (
		name,
		contact_email,
		domain_name,
		subscription_status,
		subscription_valid_upto,
		hi5_limit,
		hi5_quota_renewal_frequency,
		timezone) =
		($1, $2, $3, $4, $5, $6, $7, $8) where id = $9`

	deleteOrganizationQuery = `UPDATE organizations SET soft_delete = true, soft_delete_by = $1 WHERE id = $2`

	getOrganizationQuery = `SELECT id,
		name,
		contact_email,
		domain_name,
		subscription_status,
		subscription_valid_upto,
		hi5_limit,
		hi5_quota_renewal_frequency,
		timezone,
		created_at,
		created_by,
		updated_at FROM organizations WHERE id=$1 AND soft_delete = FALSE`

	listOrganizationsQuery = `SELECT id,
		name,
		contact_email,
		domain_name,
		subscription_status,
		subscription_valid_upto,
		hi5_limit,
		hi5_quota_renewal_frequency,
		timezone,
		created_at,
		created_by,
		updated_at FROM organizations WHERE soft_delete = FALSE ORDER BY name ASC`

	getOrganizationByDomainNameQuery = `SELECT id,
		name,
		contact_email,
		domain_name,
		subscription_status,
		subscription_valid_upto,
		hi5_limit,
		hi5_quota_renewal_frequency,
		timezone,
		created_at,
		created_by,
		updated_at FROM organizations WHERE domain_name=$1 AND soft_delete = FALSE LIMIT 1`
	getOrganizationByIDQuery = `SELECT id,
		name,
		contact_email,
		domain_name,
		subscription_status,
		subscription_valid_upto,
		hi5_limit,
		hi5_quota_renewal_frequency,
		timezone,
		created_at,
		created_by,
		updated_at FROM organizations WHERE id=$1 LIMIT 1`
	getCountOfContactEmailQuery = `SELECT COUNT(*) FROM organizations WHERE contact_email = $1 AND soft_delete = FALSE`
	getCountOfDomainNameQuery   = `SELECT COUNT(*) FROM organizations WHERE domain_name = $1 AND soft_delete = FALSE`
	getCountOfIdQuery           = `SELECT COUNT(*) FROM organizations WHERE id = $1 AND soft_delete = FALSE`
)

type OrganizationStorer interface {
	ListOrganizations(ctx context.Context) (organizations []Organization, err error)
	GetOrganization(ctx context.Context, organizationID int) (organization Organization, err error)
	GetOrganizationByDomainName(ctx context.Context, domainName string) (organization Organization, err error)
	DeleteOrganization(ctx context.Context, organizationID int, userId int64) (err error)
	UpdateOrganization(ctx context.Context, reqOrganization dto.Organization) (updatedOrganization Organization, err error)
	CreateOrganization(ctx context.Context, org dto.Organization) (createdOrganization Organization, err error)

	IsEmailPresent(ctx context.Context, email string) bool
	IsDomainPresent(ctx context.Context, domainName string) bool
	IsOrganizationIdPresent(ctx context.Context, organizationId int64) bool
}

type OrganizationStore struct {
	pgStore
}

func NewOrganizationRepo(db *sqlx.DB) OrganizationStorer {
	return &OrganizationStore{
		pgStore: pgStore{db}, // Use *sqlx.DB instead of *sql.DB
	}
}

// Organization - a struct representing an organization object in the database
type Organization struct {
	ID                       int64     `db:"id"`
	Name                     string    `db:"name"`
	ContactEmail             string    `db:"contact_email"`
	DomainName               string    `db:"domain_name"`
	SubscriptionStatus       int       `db:"subscription_status"`
	SubscriptionValidUpto    time.Time `db:"subscription_valid_upto"`
	Hi5Limit                 int       `db:"hi5_limit"`
	Hi5QuotaRenewalFrequency string    `db:"hi5_quota_renewal_frequency"`
	Timezone                 string    `db:"timezone"`
	CreatedAt                time.Time `db:"created_at"`
	CreatedBy                int64     `db:"created_by"`
	UpdatedAt                time.Time `db:"updated_at"`
	SoftDelete               bool      `db:"soft_delete"`
	SoftDeleteBy             int64     `db:"soft_delete_by"`
}

func (orgStr *OrganizationStore) ListOrganizations(ctx context.Context) (organizations []Organization, err error) {
	err = orgStr.db.Select(&organizations, listOrganizationsQuery)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing organizations")
		return organizations, ae.InernalServer
	}
	fmt.Println("orgs: ------------------->", organizations)
	return
}

func (s *OrganizationStore) CreateOrganization(ctx context.Context, org dto.Organization) (createdOrganization Organization, err error) {
	// Set org.CreatedAt so we get a valid created_at value from the database going forward
	org.CreatedAt = time.Now().UTC()

	lastInsertID := 0
	err = s.db.QueryRow(
		createOrganizationQuery,
		org.Name,
		org.ContactEmail,
		org.DomainName,
		org.SubscriptionStatus,
		org.SubscriptionValidUpto,
		org.Hi5Limit,
		org.Hi5QuotaRenewalFrequency,
		org.Timezone,
		org.CreatedBy,
	).Scan(&lastInsertID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error creating organization")
		return
	}

	err = s.db.Get(&createdOrganization, getOrganizationQuery, lastInsertID)
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: Log that we can't find the organization even though it's just been created
			log.Error(ae.ErrRecordNotFound, "Just created an Organization, but can't find it!", err)
		}
	}
	return
}

func (s *OrganizationStore) UpdateOrganization(ctx context.Context, reqOrganization dto.Organization) (updatedOrganization Organization, err error) {
	err = s.db.Get(&updatedOrganization, getOrganizationQuery, reqOrganization.ID)
	if err != nil {
		log.Error(ae.ErrRecordNotFound, "Cannot find organization id "+string(reqOrganization.ID), err)
		return Organization{}, ae.OrganizationNotFound
	}

	var dbOrganization Organization
	err = s.db.Get(&dbOrganization, getOrganizationQuery, reqOrganization.ID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching organization")
		return
	}

	updateFields := []string{}
	args := []interface{}{}
	argID := 1

	if reqOrganization.Name != "" {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", argID))
		args = append(args, reqOrganization.Name)
		argID++
	}
	if reqOrganization.ContactEmail != "" {
		updateFields = append(updateFields, fmt.Sprintf("contact_email = $%d", argID))
		updateFields = append(updateFields,"is_email_verified = false")
		args = append(args, reqOrganization.ContactEmail)
		argID++
	}
	if reqOrganization.DomainName != "" {
		updateFields = append(updateFields, fmt.Sprintf("domain_name = $%d", argID))
		args = append(args, reqOrganization.DomainName)
		argID++
	}

	if !reqOrganization.SubscriptionValidUpto.IsZero() {
		updateFields = append(updateFields, fmt.Sprintf("subscription_valid_upto = $%d", argID))
		args = append(args, reqOrganization.SubscriptionValidUpto)
		argID++
		updateFields = append(updateFields, fmt.Sprintf("subscription_status = $%d", argID))
		args = append(args, 1)
		argID++
	}
	if reqOrganization.Hi5Limit != 0 {
		updateFields = append(updateFields, fmt.Sprintf("hi5_limit = $%d", argID))
		args = append(args, reqOrganization.Hi5Limit)
		argID++
	}
	if reqOrganization.Hi5QuotaRenewalFrequency != "" {
		updateFields = append(updateFields, fmt.Sprintf("hi5_quota_renewal_frequency = $%d", argID))
		args = append(args, reqOrganization.Hi5QuotaRenewalFrequency)
		argID++
	}
	if reqOrganization.Timezone != "" {
		updateFields = append(updateFields, fmt.Sprintf("timezone = $%d", argID))
		args = append(args, reqOrganization.Timezone)
		argID++
	}

	if len(updateFields) > 0 {

		updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argID))
		args = append(args, time.Now())
		argID++
		// Append the organization ID for the WHERE clause

		args = append(args, reqOrganization.ID)
		updateQuery := fmt.Sprintf("UPDATE organizations SET %s WHERE id = $%d", strings.Join(updateFields, ", "), argID)
		fmt.Println("update query: ------------->\n", updateQuery)
		fmt.Println("update args: ------------->\n", args)
		stmt, err := s.db.Prepare(updateQuery)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error preparing update statement")
			return Organization{}, err
		}
		defer stmt.Close()
		_, err = stmt.Exec(args...)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error executing update statement")
			return Organization{}, err
		}
	}

	err = s.db.Get(&updatedOrganization, getOrganizationQuery, reqOrganization.ID)
	if err != nil {
		log.Error(ae.ErrRecordNotFound, "Cannot find organization id "+string(reqOrganization.ID), err)
		return
	}

	return
}

func (s *OrganizationStore) DeleteOrganization(ctx context.Context, organizationID int, userId int64) (err error) {
	sqlRes, err := s.db.Exec(deleteOrganizationQuery, userId, organizationID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting organization")
		return ae.InernalServer
	}

	rowsAffected, err := sqlRes.RowsAffected()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching rows affected count")
		return ae.InernalServer
	}

	if rowsAffected == 0 {
		err = fmt.Errorf("organization with ID %d not found", organizationID)
		logger.WithField("organizationID", organizationID).Warn(err.Error())
		return ae.OrganizationNotFound
	}

	return nil
}

// GetOrganization - returns an organization from the database if it exists based on its ID primary key
func (s *OrganizationStore) GetOrganization(ctx context.Context, organizationID int) (organization Organization, err error) {
	err = s.db.Get(&organization, getOrganizationQuery, organizationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithField("organizationID", organizationID).Warn("Organization not found")
			return Organization{}, ae.OrganizationNotFound
		}
		logger.WithField("err", err.Error()).Error("Error fetching organization")
		return Organization{}, err
	}

	return
}

func (s *OrganizationStore) GetOrganizationByDomainName(ctx context.Context, domainName string) (organization Organization, err error) {
	fmt.Println("GetOrganizationByDomainName ------------------------>")
	err = s.db.Get(&organization, getOrganizationByDomainNameQuery, domainName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithField("organization domain name", domainName).Warn("Organization not found by domain name")
			return Organization{}, ae.OrganizationNotFound
		}
		logger.WithField("err", err.Error()).Error("Error fetching organization")
		return Organization{}, err
	}
	return
}

func (s *OrganizationStore) IsEmailPresent(ctx context.Context, email string) bool {

	var count int

	err := s.db.QueryRowContext(ctx, getCountOfContactEmailQuery, email).Scan(&count)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching contact email of organization by contact email id: " + email)
		return false
	}

	return count > 0
}

func (s *OrganizationStore) IsDomainPresent(ctx context.Context, domainName string) bool {
	fmt.Println("domain name repo------------------------------------------>", domainName)

	var count int

	err := s.db.QueryRowContext(ctx, getCountOfDomainNameQuery, domainName).Scan(&count)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching domain name of organization by contact email id: " + domainName)
		return false
	}

	return count > 0
}

func (s *OrganizationStore) IsOrganizationIdPresent(ctx context.Context, organizationId int64) bool {
	var count int

	err := s.db.QueryRowContext(ctx, getCountOfIdQuery, organizationId).Scan(&count)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching id of organization: " + strconv.FormatInt(organizationId, 10))
		return false
	}

	return count > 0
}

///helper functions Organization

func OrganizationToDB(org dto.Organization) Organization {
	return Organization{
		ID:                       org.ID,
		Name:                     org.Name,
		ContactEmail:             org.ContactEmail,
		DomainName:               org.DomainName,
		SubscriptionStatus:       org.SubscriptionStatus,
		SubscriptionValidUpto:    org.SubscriptionValidUpto,
		Hi5Limit:                 org.Hi5Limit,
		Hi5QuotaRenewalFrequency: org.Hi5QuotaRenewalFrequency,
		Timezone:                 org.Timezone,
		CreatedAt:                org.CreatedAt,
	}
}
