package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrDuplicateStoreName     = errors.New("a vendor with a similar storename already exists")
	ErrDuplicateBusinessEmail = errors.New("a vendore with similar business email already exists")
)

type Vendor struct {
	ID             int64  `json:"-"`
	Storename      string `json:"storename"`
	Subdomain      string `json:"subdomain"`
	Description    string `json:"description"`
	Logo_URL       string `json:"logo_url"`
	Banner_URL     string `json:"banner_url"`
	Business_Email string `json:"business_email"`
	Business_Phone string `json:"business_phone"`
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

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "vendors_email_key":
				return ErrDuplicateBusinessEmail
			case "store_name_key":
				return ErrDuplicateStoreName
			}
		}

		return err
	}

	return nil

}
