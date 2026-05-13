package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Cart struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	VendorID  string     `json:"vendor_id"`
	Status    string     `json:"status"`
	Items     []CartItem `json:"items"`
	Subtotal  float64    `json:"subtotal"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

type CartItem struct {
	ID        string   `json:"id"`
	CartID    string   `json:"cart_id"`
	ProductID string   `json:"product_id"`
	Product   Products `json:"product,omitempty"`
	Quantity  int      `json:"quantity"`
	UnitPrice float64  `json:"unit_price"`
	Subtotal  float64  `json:"subtotal"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type CartStore struct {
	db *sql.DB
}

func (s *CartStore) AddItem(ctx context.Context, userID, vendorID, productID string, quantity int) (*CartItem, error) {
	query := `
		WITH active_cart AS (
			INSERT INTO carts (user_id, vendor_id)
			VALUES ($1, $2)
			ON CONFLICT (user_id, vendor_id)
			WHERE status = 'active'
			DO UPDATE SET status = 'active', updated_at = now()
			RETURNING id
		),
		upserted_item AS (
			INSERT INTO cart_items (cart_id, product_id, quantity)
			SELECT ac.id, $3, $4
			FROM active_cart ac
			JOIN products p ON p.id = $3 AND p.vendor_id = $2 AND p.status = 'published'
			ON CONFLICT (cart_id, product_id)
			DO UPDATE SET
				quantity = cart_items.quantity + EXCLUDED.quantity,
				updated_at = now()
			RETURNING id, cart_id, product_id, quantity, created_at, updated_at
		)
		SELECT
			ui.id,
			ui.cart_id,
			ui.product_id,
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
			p.updated_at,
			ui.quantity,
			p.price AS unit_price,
			ui.quantity * p.price AS subtotal,
			ui.created_at,
			ui.updated_at
		FROM upserted_item ui
		JOIN products p ON p.id = ui.product_id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	item := &CartItem{}
	err := s.db.QueryRowContext(ctx, query, userID, vendorID, productID, quantity).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Product.ID,
		&item.Product.Name,
		&item.Product.Slug,
		&item.Product.Description,
		&item.Product.Category_ID,
		&item.Product.Price,
		&item.Product.Compare_Price,
		&item.Product.Stock_Quantity,
		&item.Product.SKU,
		&item.Product.Status,
		&item.Product.Weight,
		&item.Product.Created_At,
		&item.Product.Update_At,
		&item.Quantity,
		&item.UnitPrice,
		&item.Subtotal,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, mapCartError(err)
	}

	return item, nil
}

func (s *CartStore) GetCart(ctx context.Context, userID, vendorID string) (*Cart, error) {
	query := `
		SELECT id, user_id, vendor_id, status, created_at, updated_at
		FROM carts
		WHERE user_id = $1 AND vendor_id = $2 AND status = 'active'
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	cart := &Cart{
		UserID:   userID,
		VendorID: vendorID,
		Status:   "active",
		Items:    []CartItem{},
	}

	err := s.db.QueryRowContext(ctx, query, userID, vendorID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.VendorID,
		&cart.Status,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cart, nil
		}
		return nil, err
	}

	items, subtotal, err := s.getCartItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	cart.Items = items
	cart.Subtotal = subtotal

	return cart, nil
}

func (s *CartStore) UpdateItem(ctx context.Context, userID, vendorID, productID string, quantity int) (*CartItem, error) {
	query := `
		UPDATE cart_items ci
		SET quantity = $4, updated_at = now()
		FROM carts c, products p
		WHERE ci.cart_id = c.id
			AND ci.product_id = p.id
			AND c.user_id = $1
			AND c.vendor_id = $2
			AND c.status = 'active'
			AND ci.product_id = $3
		RETURNING
			ci.id,
			ci.cart_id,
			ci.product_id,
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
			p.updated_at,
			ci.quantity,
			p.price AS unit_price,
			ci.quantity * p.price AS subtotal,
			ci.created_at,
			ci.updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	item := &CartItem{}
	err := s.db.QueryRowContext(ctx, query, userID, vendorID, productID, quantity).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Product.ID,
		&item.Product.Name,
		&item.Product.Slug,
		&item.Product.Description,
		&item.Product.Category_ID,
		&item.Product.Price,
		&item.Product.Compare_Price,
		&item.Product.Stock_Quantity,
		&item.Product.SKU,
		&item.Product.Status,
		&item.Product.Weight,
		&item.Product.Created_At,
		&item.Product.Update_At,
		&item.Quantity,
		&item.UnitPrice,
		&item.Subtotal,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, mapCartError(err)
	}

	return item, nil
}

func (s *CartStore) RemoveItem(ctx context.Context, userID, vendorID, productID string) error {
	query := `
		DELETE FROM cart_items ci
		USING carts c
		WHERE ci.cart_id = c.id
			AND c.user_id = $1
			AND c.vendor_id = $2
			AND c.status = 'active'
			AND ci.product_id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, userID, vendorID, productID)
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

func (s *CartStore) ClearCart(ctx context.Context, userID, vendorID string) error {
	query := `
		DELETE FROM cart_items ci
		USING carts c
		WHERE ci.cart_id = c.id
			AND c.user_id = $1
			AND c.vendor_id = $2
			AND c.status = 'active'
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, vendorID)
	return err
}

func (s *CartStore) getCartItems(ctx context.Context, cartID string) ([]CartItem, float64, error) {
	query := `
		SELECT
			ci.id,
			ci.cart_id,
			ci.product_id,
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
			p.updated_at,
			ci.quantity,
			p.price AS unit_price,
			ci.quantity * p.price AS subtotal,
			ci.created_at,
			ci.updated_at
		FROM cart_items ci
		JOIN products p ON p.id = ci.product_id
		WHERE ci.cart_id = $1
		ORDER BY ci.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := []CartItem{}
	var total float64

	for rows.Next() {
		var item CartItem
		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Product.ID,
			&item.Product.Name,
			&item.Product.Slug,
			&item.Product.Description,
			&item.Product.Category_ID,
			&item.Product.Price,
			&item.Product.Compare_Price,
			&item.Product.Stock_Quantity,
			&item.Product.SKU,
			&item.Product.Status,
			&item.Product.Weight,
			&item.Product.Created_At,
			&item.Product.Update_At,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		total += item.Subtotal
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func mapCartError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23503" {
		return ErrNotFound
	}

	return err
}
