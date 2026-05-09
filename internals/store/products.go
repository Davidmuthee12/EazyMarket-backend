package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrDuplicateProductName = errors.New("a product with the similar name exists")
	ErrDuplicateProductSlug = errors.New("a product with the similar slug name exists")
)

type Products struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Description    string `json:"description"`
	Category_ID    string `json:"category_id"`
	Price          int64  `json:"price"`
	Compare_Price  int64  `json:"compare_price"`
	Stock_Quantity int64  `json:"stock_quantity"`
	SKU            string `json:"sku"`
	Status         string `json:"status"`
	Weight         int64  `json:"weight"`
	Created_At     string `json:"created_at"`
	Update_At      string `json:"updated_at"`
}

type ProductStore struct {
	db *sql.DB
}

func (s *ProductStore) CreateProduct(ctx context.Context, product *Products, vendorID string) error {
	query := `
		INSERT INTO products(
			vendor_id,
			name,
			slug,
			description,
			category_id,
			price,
			compare_price,
			stock_quantity,
			sku,
			weight
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		vendorID,
		product.Name,
		product.Slug,
		product.Description,
		product.Category_ID,
		product.Price,
		product.Compare_Price,
		product.Stock_Quantity,
		product.SKU,
		product.Weight,
	).Scan(
		&product.ID,
		&product.Created_At,
		&product.Update_At,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "products_name_key":
				return ErrDuplicateProductName
			case "products_slug_key":
				return ErrDuplicateProductSlug
			}
		}
		return err
	}

	return nil
}
