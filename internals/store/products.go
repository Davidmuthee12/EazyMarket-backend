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
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Slug           string  `json:"slug"`
	Description    string  `json:"description"`
	Category_ID    string  `json:"category_id"`
	Price          float64 `json:"price"`
	Compare_Price  float64 `json:"compare_price"`
	Stock_Quantity int     `json:"stock_quantity"`
	SKU            string  `json:"sku"`
	Status         string  `json:"status"`
	Weight         float64 `json:"weight"`
	Created_At     string  `json:"created_at"`
	Update_At      string  `json:"updated_at"`
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
			status,
			weight
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, COALESCE(NULLIF($10, ''), 'draft'), $11)
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
		product.Status,
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
			case "products_slug_key", "products_vendor_slug_idx":
				return ErrDuplicateProductSlug
			}
		}
		return err
	}

	return nil
}

func (s *ProductStore) GetAllProduct(ctx context.Context, vendorID string) ([]Products, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			description,
			COALESCE(category_id::text, ''),
			price,
			COALESCE(compare_price, 0),
			stock_quantity,
			sku,
			status,
			COALESCE(weight, 0),
			created_at,
			updated_at
		FROM products
		WHERE vendor_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		vendorID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Products{}
	for rows.Next() {
		var product Products

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Category_ID,
			&product.Price,
			&product.Compare_Price,
			&product.Stock_Quantity,
			&product.SKU,
			&product.Status,
			&product.Weight,
			&product.Created_At,
			&product.Update_At,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, ErrNotFound
	}

	return products, nil
}

func (s *ProductStore) GetProductByUUID(ctx context.Context, productID string) (*Products, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			description,
			COALESCE(category_id::text, ''),
			price,
			COALESCE(compare_price, 0),
			stock_quantity,
			sku,
			status,
			COALESCE(weight, 0),
			created_at,
			updated_at
		FROM products
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	product := &Products{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		productID,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Category_ID,
		&product.Price,
		&product.Compare_Price,
		&product.Stock_Quantity,
		&product.SKU,
		&product.Status,
		&product.Weight,
		&product.Created_At,
		&product.Update_At,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return product, nil

}

func (s *ProductStore) GetPublishedProductsByVendor(ctx context.Context, vendorID string) ([]Products, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			description,
			COALESCE(category_id::text, ''),
			price,
			COALESCE(compare_price, 0),
			stock_quantity,
			sku,
			status,
			COALESCE(weight, 0),
			created_at,
			updated_at
		FROM products
		WHERE vendor_id = $1 AND status = 'published'
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Products{}
	for rows.Next() {
		var product Products

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Category_ID,
			&product.Price,
			&product.Compare_Price,
			&product.Stock_Quantity,
			&product.SKU,
			&product.Status,
			&product.Weight,
			&product.Created_At,
			&product.Update_At,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ProductStore) GetPublishedProductBySlug(ctx context.Context, vendorID, slug string) (*Products, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			description,
			COALESCE(category_id::text, ''),
			price,
			COALESCE(compare_price, 0),
			stock_quantity,
			sku,
			status,
			COALESCE(weight, 0),
			created_at,
			updated_at
		FROM products
		WHERE vendor_id = $1 AND slug = $2 AND status = 'published'
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	product := &Products{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		vendorID,
		slug,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Category_ID,
		&product.Price,
		&product.Compare_Price,
		&product.Stock_Quantity,
		&product.SKU,
		&product.Status,
		&product.Weight,
		&product.Created_At,
		&product.Update_At,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return product, nil
}

func (s *ProductStore) UpdateProduct(ctx context.Context, product *Products, vendorID string) error {
	query := `
		UPDATE products
		SET
			name = $3,
			slug = $4,
			description = $5,
			category_id = NULLIF($6, '')::uuid,
			price = $7,
			compare_price = $8,
			stock_quantity = $9,
			sku = $10,
			status = COALESCE(NULLIF($11, ''), status),
			weight = $12,
			updated_at = now()
		WHERE id = $1 AND vendor_id = $2
		RETURNING
			id,
			name,
			slug,
			description,
			COALESCE(category_id::text, ''),
			price,
			COALESCE(compare_price, 0),
			stock_quantity,
			sku,
			status,
			COALESCE(weight, 0),
			created_at,
			updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		product.ID,
		vendorID,
		product.Name,
		product.Slug,
		product.Description,
		product.Category_ID,
		product.Price,
		product.Compare_Price,
		product.Stock_Quantity,
		product.SKU,
		product.Status,
		product.Weight,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Category_ID,
		&product.Price,
		&product.Compare_Price,
		&product.Stock_Quantity,
		&product.SKU,
		&product.Status,
		&product.Weight,
		&product.Created_At,
		&product.Update_At,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				switch pqErr.Constraint {
				case "products_name_key":
					return ErrDuplicateProductName
				case "products_slug_key", "products_vendor_slug_idx":
					return ErrDuplicateProductSlug
				}
			}
			return err
		}
	}

	return nil
}

func (s *ProductStore) DeleteProduct(ctx context.Context, productID string, vendorID string) error {
	query := `
		DELETE FROM products
		WHERE id = $1 AND vendor_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		productID,
		vendorID,
	)

	if err != nil {
		return err
	}

	return nil
}
