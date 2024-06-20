package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	ae "joshsoftware/peerly/apperrors"
	"joshsoftware/peerly/pkg/dto"

	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
)

const (
	createOTPQuery = `INSERT INTO otp (
		org_id,
		otp)
		VALUES ($1, $2, $3) RETURNING otp`
	getOTPVerificationQuery = `SELECT otp,created_at,org_id
	 FROM otp WHERE otp=$1 AND org_id=$2`

	 getCountOfOrgIdForOTPQuery = `SELECT COUNT(*) FROM otp WHERE org_id=$1`
	 ChangeIsVerifiedFlagQuery = `UPDATE organizations SET is_email_verified = true WHERE id = $1`
)

type OTPVerificationStorer interface {
	GetOTPVerificationStatus(ctx context.Context,otpReq dto.OTP)(otpInfo OTP,err error)
	GetCountOfOrgId(ctx context.Context,orgId int64)(count int,err error)
	ChangeIsVerifiedFlag(ctx context.Context,organizationID int64)(error)
	CreateOTPInfo(ctx context.Context,otpinfo OTP)(error)
	DeleteOTPData(ctx context.Context,orgId int64)(error)
}

type OTPVerificationStore struct {
	pgStore
}

func NewOTPVerificationRepo(db *sqlx.DB) OTPVerificationStorer {
	return &OTPVerificationStore{
		pgStore: pgStore{db},
	}
}

type OTP struct {
	CreatedAt   time.Time `db:"created_at"`
	OrgId      int64 `db:"org_id"`
	OTPCode     string `db:"otp"`
}

func (otp *OTPVerificationStore) GetOTPVerificationStatus(ctx context.Context, otpReq dto.OTP) (otpInfo OTP, err error) {
    fmt.Println("-------------------------------------------->repo")
    err = otp.db.Get(&otpInfo, getOTPVerificationQuery, otpReq.OTPCode,otpReq.OrgId)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            logger.WithField("organizationID", otpReq.OTPCode).Warn("Otp not found")
            return OTP{}, ae.InvalidOTP
        }
        logger.WithField("err", err.Error()).Error("Error fetching organization")
        return OTP{}, err
    }
	fmt.Println("otp Info: ",otpInfo)

    return
}


func (otp *OTPVerificationStore) GetCountOfOrgId(ctx context.Context,orgId int64)(count int,err error){
	
	err = otp.db.Get(&count, getCountOfOrgIdForOTPQuery, orgId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithField("organizationID", orgId).Warn("orgid in Otp not found")
			return 0, ae.InvalidReferenceId
		}
		logger.WithField("err", err.Error()).Error("Error fetching organization")
		return 0, err
	}
	return 
}

func (otp *OTPVerificationStore) ChangeIsVerifiedFlag(ctx context.Context,organizationID int64)(error){
	sqlRes, err := otp.db.Exec(ChangeIsVerifiedFlagQuery, organizationID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error dupdating is verified flag")
		return ae.InernalServer
	}

	rowsAffected, err := sqlRes.RowsAffected()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching rows affected count")
		return ae.InernalServer
	}

	if rowsAffected == 0 {
		// err = fmt.Errorf("organization with ID %d not found", organizationID)
		logger.WithField("organizationID", organizationID).Warn(err.Error())
		return ae.OrganizationNotFound
	}

	return nil
}

func (otp *OTPVerificationStore) CreateOTPInfo(ctx context.Context, otpInfo OTP) error {
    // Check if the OTP already exists
    var existingOTP string
    checkOTPQuery := "SELECT otp FROM otp WHERE otp = $1"
    err := otp.db.QueryRowContext(ctx, checkOTPQuery, otpInfo.OTPCode).Scan(&existingOTP)
    if err == nil {
        logger.WithField("otp", otpInfo.OTPCode).Warn("OTP already exists")
        return ae.ErrOTPAlreadyExists
    } else if err != sql.ErrNoRows {
        logger.WithField("err", err.Error()).Error("Error checking existing OTP")
        return err
    }

    // Set the current time for created_at
    otpInfo.CreatedAt = time.Now().UTC()

    // Use a named query to insert the OTP info
    createOTPQuery := `
        INSERT INTO otp (org_id, otp) 
        VALUES ($1, $2)
        RETURNING otp`

    // Insert the OTP info into the database and retrieve the otp code
    var insertedOTPCode string
    err = otp.db.QueryRowContext(ctx, createOTPQuery, otpInfo.OrgId, otpInfo.OTPCode).Scan(&insertedOTPCode)
    if err != nil {
        logger.WithField("err", err.Error()).Error("Error creating OTP")
        return err
    }

    // Retrieve the newly created OTP info
    err = otp.db.GetContext(ctx, &otpInfo, getOTPVerificationQuery, insertedOTPCode,otpInfo.OrgId)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            logger.WithField("Just created an OTP, but can't find it!", otpInfo).Warn(err.Error())
            return ae.ErrRecordNotFound
        }
        logger.WithField("err", err.Error()).Error("Error fetching created OTP")
        return err
    }

    return nil
}

func (otp *OTPVerificationStore) DeleteOTPData(ctx context.Context, orgId int64) error {
    
    deleteOTPQuery := "DELETE FROM otp WHERE org_id = $1"
    
    result, err := otp.db.ExecContext(ctx, deleteOTPQuery, orgId)
    if err != nil {
        logger.WithField("err", err.Error()).Error("Error deleting OTP data")
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        logger.WithField("err", err.Error()).Error("Error fetching rows affected")
        return err
    }
    logger.WithField("org_id", orgId).Infof("Deleted %d rows from OTP table", rowsAffected)

    return nil
}
