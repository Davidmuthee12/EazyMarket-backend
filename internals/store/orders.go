package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

const platformCommissionRate = 0.10

type Order struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	Status          string          `json:"status"`
	TotalAmount     float64         `json:"total_amount"`
	ShippingAddress json.RawMessage `json:"shipping_address,omitempty" swaggertype:"object"`
	Notes           string          `json:"notes,omitempty"`
	Items           []OrderItem     `json:"items"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

type OrderItem struct {
	ID            string   `json:"id"`
	OrderID       string   `json:"order_id"`
	ProductID     string   `json:"product_id"`
	VendorID      string   `json:"vendor_id"`
	Product       Products `json:"product,omitempty"`
	Quantity      int      `json:"quantity"`
	UnitPrice     float64  `json:"unit_price"`
	Subtotal      float64  `json:"subtotal"`
	CommissionAmt float64  `json:"commission_amt"`
	VendorPayout  float64  `json:"vendor_payout"`
}

type OrderStore struct {
	db *sql.DB
}

func (s *OrderStore) CreateFromCart(ctx context.Context, userID string, shippingAddress []byte, notes string) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	order := &Order{}

	err := withTx(s.db, ctx, func(tx *sql.Tx) error {
		var cartID string
		err := tx.QueryRowContext(ctx, `
			SELECT id
			FROM carts
			WHERE user_id = $1 AND status = 'active'
		`, userID).Scan(&cartID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrEmptyCart
			}
			return err
		}

		var totalAmount float64
		err = tx.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(ci.quantity * p.price), 0)
			FROM cart_items ci
			JOIN products p ON p.id = ci.product_id
			WHERE ci.cart_id = $1
		`, cartID).Scan(&totalAmount)
		if err != nil {
			return err
		}
		if totalAmount == 0 {
			return ErrEmptyCart
		}

		err = tx.QueryRowContext(ctx, `
			INSERT INTO orders (user_id, total_amount, shipping_address, notes)
			VALUES ($1, $2, NULLIF($3, '')::jsonb, NULLIF($4, ''))
			RETURNING id, user_id, status, total_amount, COALESCE(shipping_address, '{}'::jsonb), COALESCE(notes, ''), created_at, updated_at
		`, userID, totalAmount, string(shippingAddress), notes).Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmount,
			&order.ShippingAddress,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO order_items (
				order_id,
				product_id,
				vendor_id,
				quantity,
				unit_price,
				subtotal,
				commission_amt,
				vendor_payout
			)
			SELECT
				$1,
				p.id,
				p.vendor_id,
				ci.quantity,
				p.price,
				ci.quantity * p.price,
				(ci.quantity * p.price) * $2,
				(ci.quantity * p.price) - ((ci.quantity * p.price) * $2)
			FROM cart_items ci
			JOIN products p ON p.id = ci.product_id
			WHERE ci.cart_id = $3
		`, order.ID, platformCommissionRate, cartID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			DELETE FROM cart_items
			WHERE cart_id = $1
		`, cartID)
		if err != nil {
			return err
		}

		items, err := s.getOrderItems(ctx, tx, order.ID, "")
		if err != nil {
			return err
		}

		order.Items = items
		return nil
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderStore) GetAll(ctx context.Context, userID string) ([]Order, error) {
	query := `
		SELECT id, user_id, status, total_amount, COALESCE(shipping_address, '{}'::jsonb), COALESCE(notes, ''), created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	orders, err := s.getOrders(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, ErrNotFound
	}

	return orders, nil
}

func (s *OrderStore) GetByID(ctx context.Context, userID, orderID string) (*Order, error) {
	query := `
		SELECT id, user_id, status, total_amount, COALESCE(shipping_address, '{}'::jsonb), COALESCE(notes, ''), created_at, updated_at
		FROM orders
		WHERE id = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	order, err := s.getOrder(ctx, query, orderID, userID)
	if err != nil {
		return nil, err
	}

	items, err := s.getOrderItems(ctx, s.db, order.ID, "")
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (s *OrderStore) Cancel(ctx context.Context, userID, orderID string) (*Order, error) {
	query := `
		UPDATE orders
		SET status = 'cancelled', updated_at = now()
		WHERE id = $1
			AND user_id = $2
			AND status IN ('pending', 'confirmed')
		RETURNING id, user_id, status, total_amount, COALESCE(shipping_address, '{}'::jsonb), COALESCE(notes, ''), created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	order, err := s.getOrder(ctx, query, orderID, userID)
	if err != nil {
		return nil, err
	}

	items, err := s.getOrderItems(ctx, s.db, order.ID, "")
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (s *OrderStore) GetVendorOrders(ctx context.Context, vendorID string) ([]Order, error) {
	query := `
		SELECT DISTINCT o.id, o.user_id, o.status, o.total_amount, COALESCE(o.shipping_address, '{}'::jsonb), COALESCE(o.notes, ''), o.created_at, o.updated_at
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id
		WHERE oi.vendor_id = $1
		ORDER BY o.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	orders, err := s.getOrders(ctx, query, vendorID)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, ErrNotFound
	}

	for i := range orders {
		items, err := s.getOrderItems(ctx, s.db, orders[i].ID, vendorID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	return orders, nil
}

func (s *OrderStore) GetVendorOrderByID(ctx context.Context, vendorID, orderID string) (*Order, error) {
	query := `
		SELECT DISTINCT o.id, o.user_id, o.status, o.total_amount, COALESCE(o.shipping_address, '{}'::jsonb), COALESCE(o.notes, ''), o.created_at, o.updated_at
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id
		WHERE o.id = $1 AND oi.vendor_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	order, err := s.getOrder(ctx, query, orderID, vendorID)
	if err != nil {
		return nil, err
	}

	items, err := s.getOrderItems(ctx, s.db, order.ID, vendorID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (s *OrderStore) UpdateVendorOrderStatus(ctx context.Context, vendorID, orderID, status string) (*Order, error) {
	query := `
		UPDATE orders o
		SET status = $3, updated_at = now()
		WHERE o.id = $1
			AND EXISTS (
				SELECT 1
				FROM order_items oi
				WHERE oi.order_id = o.id AND oi.vendor_id = $2
			)
		RETURNING o.id, o.user_id, o.status, o.total_amount, COALESCE(o.shipping_address, '{}'::jsonb), COALESCE(o.notes, ''), o.created_at, o.updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	order, err := s.getOrder(ctx, query, orderID, vendorID, status)
	if err != nil {
		return nil, err
	}

	items, err := s.getOrderItems(ctx, s.db, order.ID, vendorID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (s *OrderStore) getOrders(ctx context.Context, query string, args ...any) ([]Order, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmount,
			&order.ShippingAddress,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderStore) getOrder(ctx context.Context, query string, args ...any) (*Order, error) {
	order := &Order{}
	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.ShippingAddress,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return order, nil
}

type queryer interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func (s *OrderStore) getOrderItems(ctx context.Context, q queryer, orderID, vendorID string) ([]OrderItem, error) {
	query := `
		SELECT
			oi.id,
			oi.order_id,
			oi.product_id,
			oi.vendor_id,
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
			oi.quantity,
			oi.unit_price,
			oi.subtotal,
			oi.commission_amt,
			oi.vendor_payout
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = $1
	`

	args := []any{orderID}
	if vendorID != "" {
		query += ` AND oi.vendor_id = $2`
		args = append(args, vendorID)
	}

	query += `
		ORDER BY oi.id
	`

	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []OrderItem{}
	for rows.Next() {
		var item OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.VendorID,
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
			&item.CommissionAmt,
			&item.VendorPayout,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
