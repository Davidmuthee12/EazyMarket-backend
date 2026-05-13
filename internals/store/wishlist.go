package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Wishlist struct {
	UserID    string   `json:"user_id"`
	ProductID string   `json:"product_id"`
	Product   Products `json:"product,omitempty"`
	CreatedAt string   `json:"created_at"`
}

type WishlistStore struct {
	db *sql.DB
}

func (s *WishlistStore) AddToWishList(ctx context.Context, userUUID, vendorID, productID string) (*Wishlist, error) {
	query := `
		INSERT INTO wishlists (user_id, product_id)
		SELECT $1, p.id
		FROM products p
		WHERE p.id = $2 AND p.vendor_id = $3 AND p.status = 'published'
		RETURNING user_id, product_id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	wishlist := &Wishlist{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userUUID,
		productID,
		vendorID,
	).Scan(
		&wishlist.UserID,
		&wishlist.ProductID,
		&wishlist.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23503":
				return nil, ErrNotFound
			case "23505":
				return nil, ErrConflict
			}
		}

		return nil, err
	}

	return wishlist, nil
}

func (s *WishlistStore) GetUserWishlist(ctx context.Context, userUUID, vendorID string) ([]Wishlist, error) {
	query := `
		SELECT
			w.user_id,
			w.product_id,
			w.created_at,
			p.id,
			p.name,
			p.slug,
			p.description,
			COALESCE(p.category_id::text, ''),
			p.price,
			COALESCE(p.compare_price, 0),
			p.stock_quantity,
			COALESCE(p.sku, ''),
			p.status,
			COALESCE(p.weight, 0),
			p.created_at,
			p.updated_at
		FROM wishlists w
		JOIN products p ON p.id = w.product_id
		WHERE w.user_id = $1 AND p.vendor_id = $2
		ORDER BY w.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userUUID, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wishlists := []Wishlist{}

	for rows.Next() {
		var wishlist Wishlist

		err := rows.Scan(
			&wishlist.UserID,
			&wishlist.ProductID,
			&wishlist.CreatedAt,
			&wishlist.Product.ID,
			&wishlist.Product.Name,
			&wishlist.Product.Slug,
			&wishlist.Product.Description,
			&wishlist.Product.Category_ID,
			&wishlist.Product.Price,
			&wishlist.Product.Compare_Price,
			&wishlist.Product.Stock_Quantity,
			&wishlist.Product.SKU,
			&wishlist.Product.Status,
			&wishlist.Product.Weight,
			&wishlist.Product.Created_At,
			&wishlist.Product.Update_At,
		)

		if err != nil {
			return nil, err
		}

		wishlists = append(wishlists, wishlist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wishlists, nil
}

func (s *WishlistStore) GetWishlistByID(ctx context.Context, userUUID, vendorID, productID string) (*Wishlist, error) {
	query := `
		SELECT
			w.user_id,
			w.product_id,
			w.created_at,
			p.id,
			p.name,
			p.slug,
			p.description,
			COALESCE(p.category_id::text, ''),
			p.price,
			COALESCE(p.compare_price, 0),
			p.stock_quantity,
			COALESCE(p.sku, ''),
			p.status,
			COALESCE(p.weight, 0),
			p.created_at,
			p.updated_at
		FROM wishlists w
		JOIN products p ON p.id = w.product_id
		WHERE w.user_id = $1 AND p.vendor_id = $2 AND w.product_id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	wishlist := &Wishlist{}
	err := s.db.QueryRowContext(ctx, query, userUUID, vendorID, productID).Scan(
		&wishlist.UserID,
		&wishlist.ProductID,
		&wishlist.CreatedAt,
		&wishlist.Product.ID,
		&wishlist.Product.Name,
		&wishlist.Product.Slug,
		&wishlist.Product.Description,
		&wishlist.Product.Category_ID,
		&wishlist.Product.Price,
		&wishlist.Product.Compare_Price,
		&wishlist.Product.Stock_Quantity,
		&wishlist.Product.SKU,
		&wishlist.Product.Status,
		&wishlist.Product.Weight,
		&wishlist.Product.Created_At,
		&wishlist.Product.Update_At,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return wishlist, nil
}

func (s *WishlistStore) DeleteFromWishlist(ctx context.Context, userUUID, vendorID, productID string) error {
	query := `
		DELETE FROM wishlists w
		USING products p
		WHERE w.product_id = p.id
			AND w.user_id = $1
			AND p.vendor_id = $2
			AND w.product_id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, userUUID, vendorID, productID)
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
