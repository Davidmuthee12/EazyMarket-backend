package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrDuplicateStoreName     = errors.New("a vendor with a similar storename already exists")
	ErrDuplicateSubdomain     = errors.New("a vendor with a similar subdomain already exists")
	ErrDuplicateBusinessEmail = errors.New("a vendore with similar business email already exists")
	ErrVendorProfileExists    = errors.New("vendor profile already exists")
)

type Vendor struct {
	ID             int64  `json:"-"`
	UserID         string `json:"user_id"`
	Storename      string `json:"storename"`
	Subdomain      string `json:"subdomain"`
	Description    string `json:"description"`
	Logo_URL       string `json:"logo_url"`
	Banner_URL     string `json:"banner_url"`
	Business_Email string `json:"business_email"`
	Business_Phone string `json:"business_phone"`
	Status         string `json:"status"`
	Address        string `json:"address"`
	CreatedAt      string `json:"created_at"`
}

type VenderStore struct {
	db *sql.DB
}

func (s *VenderStore) CreateVendorProfile(ctx context.Context, Vendor *Vendor, userUUID string) error {
	query := `
		INSERT INTO vendor_profiles(
			user_id,
			store_name,
			subdomain,
			description,
			logo_url,
			banner_url,
			business_email,
			business_phone,
			address
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		userUUID,
		Vendor.Storename,
		Vendor.Subdomain,
		Vendor.Description,
		Vendor.Logo_URL,
		Vendor.Banner_URL,
		Vendor.Business_Email,
		Vendor.Business_Phone,
		Vendor.Address,
	).Scan(
		&Vendor.ID,
		&Vendor.CreatedAt,
	)
	Vendor.UserID = userUUID

	if err != nil {
		return mapVendorProfileCreateError(err)
	}

	return nil

}

func mapVendorProfileCreateError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		switch pqErr.Constraint {
		case "vendor_profiles_user_id_key":
			return ErrVendorProfileExists
		case "vendors_email_key":
			return ErrDuplicateBusinessEmail
		case "vendor_profiles_subdomain_key":
			return ErrDuplicateSubdomain
		case "store_name_key":
			return ErrDuplicateStoreName
		}
	}

	return err
}

func (s *VenderStore) GetVendorByUUID(ctx context.Context, userID string) (*Vendor, error) {
	query := `
		SELECT
			user_id,
			store_name,
			subdomain,
			description,
			logo_url,
			banner_url,
			business_email,
			business_phone,
			status,
			address,
			created_at
		FROM vendor_profiles
		WHERE user_id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	vendor := &Vendor{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&vendor.UserID,
		&vendor.Storename,
		&vendor.Subdomain,
		&vendor.Description,
		&vendor.Logo_URL,
		&vendor.Banner_URL,
		&vendor.Business_Email,
		&vendor.Business_Phone,
		&vendor.Status,
		&vendor.Address,
		&vendor.CreatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return vendor, nil
}

func (s *VenderStore) GetVendorBySubdomain(ctx context.Context, subdomain string) (*Vendor, error) {
	query := `
		SELECT
			user_id,
			store_name,
			subdomain,
			description,
			logo_url,
			banner_url,
			business_email,
			business_phone,
			status,
			address,
			created_at
		FROM vendor_profiles
		WHERE subdomain = $1 AND status = 'approved'
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	vendor := &Vendor{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		subdomain,
	).Scan(
		&vendor.UserID,
		&vendor.Storename,
		&vendor.Subdomain,
		&vendor.Description,
		&vendor.Logo_URL,
		&vendor.Banner_URL,
		&vendor.Business_Email,
		&vendor.Business_Phone,
		&vendor.Status,
		&vendor.Address,
		&vendor.CreatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return vendor, nil
}

func (s *VenderStore) SetStatus(ctx context.Context, userUUID, status string) error {
	query := `
		UPDATE vendor_profiles
		SET status = $2, updated_at = NOW()
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, userUUID, status)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
