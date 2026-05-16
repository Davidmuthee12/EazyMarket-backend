package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type AnalyticsSummary struct {
	OrdersCount        int     `json:"orders_count"`
	GrossRevenue       float64 `json:"gross_revenue"`
	VendorPayout       float64 `json:"vendor_payout"`
	PlatformCommission float64 `json:"platform_commission"`
	AverageOrderValue  float64 `json:"average_order_value"`
	ProductsSold       int     `json:"products_sold"`
	CustomersCount     int     `json:"customers_count"`
	PendingOrders      int     `json:"pending_orders"`
	ConfirmedOrders    int     `json:"confirmed_orders"`
	ShippedOrders      int     `json:"shipped_orders"`
	DeliveredOrders    int     `json:"delivered_orders"`
	CancelledOrders    int     `json:"cancelled_orders"`
	StorefrontViews    int     `json:"storefront_views"`
	ProductViews       int     `json:"product_views"`
	AddToCartEvents    int     `json:"add_to_cart_events"`
	WishlistAddEvents  int     `json:"wishlist_add_events"`
}

type AdminAnalyticsSummary struct {
	AnalyticsSummary
	TotalUsers                int `json:"total_users"`
	TotalCustomers            int `json:"total_customers"`
	TotalVendors              int `json:"total_vendors"`
	ApprovedVendors           int `json:"approved_vendors"`
	PendingVendorApplications int `json:"pending_vendor_applications"`
	SuspendedVendors          int `json:"suspended_vendors"`
}

type RevenuePoint struct {
	Period             string  `json:"period"`
	OrdersCount        int     `json:"orders_count"`
	GrossRevenue       float64 `json:"gross_revenue"`
	VendorPayout       float64 `json:"vendor_payout"`
	PlatformCommission float64 `json:"platform_commission"`
}

type TopProductAnalytics struct {
	ProductID    string  `json:"product_id"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	QuantitySold int     `json:"quantity_sold"`
	GrossRevenue float64 `json:"gross_revenue"`
	ProductViews int     `json:"product_views"`
	AddToCart    int     `json:"add_to_cart"`
	WishlistAdds int     `json:"wishlist_adds"`
}

type TopVendorAnalytics struct {
	VendorID     string  `json:"vendor_id"`
	StoreName    string  `json:"store_name"`
	Subdomain    string  `json:"subdomain"`
	OrdersCount  int     `json:"orders_count"`
	GrossRevenue float64 `json:"gross_revenue"`
	VendorPayout float64 `json:"vendor_payout"`
	ProductsSold int     `json:"products_sold"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type CustomerAnalytics struct {
	TotalCustomers  int `json:"total_customers"`
	RepeatCustomers int `json:"repeat_customers"`
}

type UserAnalytics struct {
	TotalUsers     int `json:"total_users"`
	ActiveUsers    int `json:"active_users"`
	SuspendedUsers int `json:"suspended_users"`
	Customers      int `json:"customers"`
	Vendors        int `json:"vendors"`
	Admins         int `json:"admins"`
}

type VendorApplicationAnalytics struct {
	Pending  int `json:"pending"`
	Approved int `json:"approved"`
	Rejected int `json:"rejected"`
}

type AnalyticsEvent struct {
	ID        string          `json:"id"`
	VendorID  string          `json:"vendor_id"`
	UserID    string          `json:"user_id,omitempty"`
	SessionID string          `json:"session_id,omitempty"`
	EventType string          `json:"event_type"`
	ProductID string          `json:"product_id,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty" swaggertype:"object"`
	CreatedAt string          `json:"created_at"`
}

type AnalyticsStore struct {
	db *sql.DB
}

func (s *AnalyticsStore) TrackEvent(ctx context.Context, event *AnalyticsEvent) error {
	query := `
		INSERT INTO analytics_events (vendor_id, user_id, session_id, event_type, product_id, metadata)
		SELECT $1, NULLIF($2, '')::uuid, NULLIF($3, ''), $4, NULLIF($5, '')::uuid, COALESCE(NULLIF($6, '')::jsonb, '{}'::jsonb)
		WHERE NULLIF($5, '') IS NULL
			OR EXISTS (
				SELECT 1
				FROM products p
				WHERE p.id = NULLIF($5, '')::uuid AND p.vendor_id = $1
			)
		RETURNING id, vendor_id, COALESCE(user_id::text, ''), COALESCE(session_id, ''), event_type, COALESCE(product_id::text, ''), metadata, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	metadata := string(event.Metadata)
	err := s.db.QueryRowContext(ctx, query, event.VendorID, event.UserID, event.SessionID, event.EventType, event.ProductID, metadata).Scan(
		&event.ID,
		&event.VendorID,
		&event.UserID,
		&event.SessionID,
		&event.EventType,
		&event.ProductID,
		&event.Metadata,
		&event.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (s *AnalyticsStore) GetVendorSummary(ctx context.Context, vendorID string, from, to *time.Time) (*AnalyticsSummary, error) {
	query := `
		WITH order_totals AS (
			SELECT
				o.id,
				o.user_id,
				o.status,
				SUM(oi.subtotal) AS gross_revenue,
				SUM(oi.vendor_payout) AS vendor_payout,
				SUM(oi.commission_amt) AS platform_commission,
				SUM(oi.quantity) AS products_sold
			FROM orders o
			JOIN order_items oi ON oi.order_id = o.id
			WHERE oi.vendor_id = $1
				AND ($2::timestamptz IS NULL OR o.created_at >= $2)
				AND ($3::timestamptz IS NULL OR o.created_at <= $3)
			GROUP BY o.id, o.user_id, o.status
		),
		order_agg AS (
			SELECT
				COUNT(id) AS orders_count,
				COALESCE(SUM(gross_revenue), 0) AS gross_revenue,
				COALESCE(SUM(vendor_payout), 0) AS vendor_payout,
				COALESCE(SUM(platform_commission), 0) AS platform_commission,
				COALESCE(AVG(gross_revenue), 0) AS average_order_value,
				COALESCE(SUM(products_sold), 0) AS products_sold,
				COUNT(DISTINCT user_id) AS customers_count,
				COUNT(*) FILTER (WHERE status = 'pending') AS pending_orders,
				COUNT(*) FILTER (WHERE status = 'confirmed') AS confirmed_orders,
				COUNT(*) FILTER (WHERE status = 'shipped') AS shipped_orders,
				COUNT(*) FILTER (WHERE status = 'delivered') AS delivered_orders,
				COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled_orders
			FROM order_totals
		),
		event_totals AS (
			SELECT
				COUNT(*) FILTER (WHERE event_type = 'storefront_view') AS storefront_views,
				COUNT(*) FILTER (WHERE event_type = 'product_view') AS product_views,
				COUNT(*) FILTER (WHERE event_type = 'add_to_cart') AS add_to_cart_events,
				COUNT(*) FILTER (WHERE event_type = 'wishlist_add') AS wishlist_add_events
			FROM analytics_events
			WHERE vendor_id = $1
				AND ($2::timestamptz IS NULL OR created_at >= $2)
				AND ($3::timestamptz IS NULL OR created_at <= $3)
		)
		SELECT
			oa.orders_count,
			oa.gross_revenue,
			oa.vendor_payout,
			oa.platform_commission,
			oa.average_order_value,
			oa.products_sold,
			oa.customers_count,
			oa.pending_orders,
			oa.confirmed_orders,
			oa.shipped_orders,
			oa.delivered_orders,
			oa.cancelled_orders,
			COALESCE(et.storefront_views, 0),
			COALESCE(et.product_views, 0),
			COALESCE(et.add_to_cart_events, 0),
			COALESCE(et.wishlist_add_events, 0)
		FROM order_agg oa
		CROSS JOIN event_totals et
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	summary := &AnalyticsSummary{}
	err := s.db.QueryRowContext(ctx, query, vendorID, from, to).Scan(
		&summary.OrdersCount,
		&summary.GrossRevenue,
		&summary.VendorPayout,
		&summary.PlatformCommission,
		&summary.AverageOrderValue,
		&summary.ProductsSold,
		&summary.CustomersCount,
		&summary.PendingOrders,
		&summary.ConfirmedOrders,
		&summary.ShippedOrders,
		&summary.DeliveredOrders,
		&summary.CancelledOrders,
		&summary.StorefrontViews,
		&summary.ProductViews,
		&summary.AddToCartEvents,
		&summary.WishlistAddEvents,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return summary, nil
}

func (s *AnalyticsStore) GetAdminSummary(ctx context.Context, from, to *time.Time) (*AdminAnalyticsSummary, error) {
	query := `
		WITH order_totals AS (
			SELECT
				o.id,
				o.user_id,
				o.status,
				SUM(oi.subtotal) AS gross_revenue,
				SUM(oi.vendor_payout) AS vendor_payout,
				SUM(oi.commission_amt) AS platform_commission,
				SUM(oi.quantity) AS products_sold
			FROM orders o
			JOIN order_items oi ON oi.order_id = o.id
			WHERE ($1::timestamptz IS NULL OR o.created_at >= $1)
				AND ($2::timestamptz IS NULL OR o.created_at <= $2)
			GROUP BY o.id, o.user_id, o.status
		),
		order_agg AS (
			SELECT
				COUNT(id) AS orders_count,
				COALESCE(SUM(gross_revenue), 0) AS gross_revenue,
				COALESCE(SUM(vendor_payout), 0) AS vendor_payout,
				COALESCE(SUM(platform_commission), 0) AS platform_commission,
				COALESCE(AVG(gross_revenue), 0) AS average_order_value,
				COALESCE(SUM(products_sold), 0) AS products_sold,
				COUNT(DISTINCT user_id) AS customers_count,
				COUNT(*) FILTER (WHERE status = 'pending') AS pending_orders,
				COUNT(*) FILTER (WHERE status = 'confirmed') AS confirmed_orders,
				COUNT(*) FILTER (WHERE status = 'shipped') AS shipped_orders,
				COUNT(*) FILTER (WHERE status = 'delivered') AS delivered_orders,
				COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled_orders
			FROM order_totals
		),
		event_totals AS (
			SELECT
				COUNT(*) FILTER (WHERE event_type = 'storefront_view') AS storefront_views,
				COUNT(*) FILTER (WHERE event_type = 'product_view') AS product_views,
				COUNT(*) FILTER (WHERE event_type = 'add_to_cart') AS add_to_cart_events,
				COUNT(*) FILTER (WHERE event_type = 'wishlist_add') AS wishlist_add_events
			FROM analytics_events
			WHERE ($1::timestamptz IS NULL OR created_at >= $1)
				AND ($2::timestamptz IS NULL OR created_at <= $2)
		),
		user_totals AS (
			SELECT
				COUNT(*) AS total_users,
				COUNT(*) FILTER (WHERE role = 'user') AS total_customers,
				COUNT(*) FILTER (WHERE role = 'vendor') AS total_vendors
			FROM users
		),
		vendor_totals AS (
			SELECT
				COUNT(*) FILTER (WHERE status = 'approved') AS approved_vendors,
				COUNT(*) FILTER (WHERE status = 'pending') AS pending_vendor_applications,
				COUNT(*) FILTER (WHERE status = 'suspended') AS suspended_vendors
			FROM vendor_profiles
		)
		SELECT
			oa.orders_count,
			oa.gross_revenue,
			oa.vendor_payout,
			oa.platform_commission,
			oa.average_order_value,
			oa.products_sold,
			oa.customers_count,
			oa.pending_orders,
			oa.confirmed_orders,
			oa.shipped_orders,
			oa.delivered_orders,
			oa.cancelled_orders,
			COALESCE(et.storefront_views, 0),
			COALESCE(et.product_views, 0),
			COALESCE(et.add_to_cart_events, 0),
			COALESCE(et.wishlist_add_events, 0),
			ut.total_users,
			ut.total_customers,
			ut.total_vendors,
			vt.approved_vendors,
			vt.pending_vendor_applications,
			vt.suspended_vendors
		FROM order_agg oa
		CROSS JOIN event_totals et
		CROSS JOIN user_totals ut
		CROSS JOIN vendor_totals vt
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	summary := &AdminAnalyticsSummary{}
	err := s.db.QueryRowContext(ctx, query, from, to).Scan(
		&summary.OrdersCount,
		&summary.GrossRevenue,
		&summary.VendorPayout,
		&summary.PlatformCommission,
		&summary.AverageOrderValue,
		&summary.ProductsSold,
		&summary.CustomersCount,
		&summary.PendingOrders,
		&summary.ConfirmedOrders,
		&summary.ShippedOrders,
		&summary.DeliveredOrders,
		&summary.CancelledOrders,
		&summary.StorefrontViews,
		&summary.ProductViews,
		&summary.AddToCartEvents,
		&summary.WishlistAddEvents,
		&summary.TotalUsers,
		&summary.TotalCustomers,
		&summary.TotalVendors,
		&summary.ApprovedVendors,
		&summary.PendingVendorApplications,
		&summary.SuspendedVendors,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return summary, nil
}

func (s *AnalyticsStore) GetRevenue(ctx context.Context, vendorID string, from, to *time.Time, interval string) ([]RevenuePoint, error) {
	vendorClause := ""
	args := []any{from, to}
	if vendorID != "" {
		vendorClause = "AND oi.vendor_id = $3"
		args = append(args, vendorID)
	}

	query := fmt.Sprintf(`
		SELECT
			date_trunc('%s', o.created_at)::date::text AS period,
			COUNT(DISTINCT o.id) AS orders_count,
			COALESCE(SUM(oi.subtotal), 0) AS gross_revenue,
			COALESCE(SUM(oi.vendor_payout), 0) AS vendor_payout,
			COALESCE(SUM(oi.commission_amt), 0) AS platform_commission
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id
		WHERE ($1::timestamptz IS NULL OR o.created_at >= $1)
			AND ($2::timestamptz IS NULL OR o.created_at <= $2)
			%s
		GROUP BY period
		ORDER BY period
	`, interval, vendorClause)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := []RevenuePoint{}
	for rows.Next() {
		var point RevenuePoint
		if err := rows.Scan(&point.Period, &point.OrdersCount, &point.GrossRevenue, &point.VendorPayout, &point.PlatformCommission); err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, rows.Err()
}

func (s *AnalyticsStore) GetTopProducts(ctx context.Context, vendorID string, from, to *time.Time, limit int) ([]TopProductAnalytics, error) {
	query := `
		WITH sales AS (
			SELECT
				p.id,
				p.name,
				p.slug,
				COALESCE(SUM(oi.quantity) FILTER (WHERE o.id IS NOT NULL), 0) AS quantity_sold,
				COALESCE(SUM(oi.subtotal) FILTER (WHERE o.id IS NOT NULL), 0) AS gross_revenue
			FROM products p
			LEFT JOIN order_items oi ON oi.product_id = p.id
			LEFT JOIN orders o ON o.id = oi.order_id
				AND ($2::timestamptz IS NULL OR o.created_at >= $2)
				AND ($3::timestamptz IS NULL OR o.created_at <= $3)
			WHERE (NULLIF($1, '')::uuid IS NULL OR p.vendor_id = NULLIF($1, '')::uuid)
			GROUP BY p.id, p.name, p.slug
		),
		events AS (
			SELECT
				product_id,
				COUNT(*) FILTER (WHERE event_type = 'product_view') AS product_views,
				COUNT(*) FILTER (WHERE event_type = 'add_to_cart') AS add_to_cart,
				COUNT(*) FILTER (WHERE event_type = 'wishlist_add') AS wishlist_adds
			FROM analytics_events
			WHERE product_id IS NOT NULL
				AND (NULLIF($1, '')::uuid IS NULL OR vendor_id = NULLIF($1, '')::uuid)
				AND ($2::timestamptz IS NULL OR created_at >= $2)
				AND ($3::timestamptz IS NULL OR created_at <= $3)
			GROUP BY product_id
		)
		SELECT
			s.id,
			s.name,
			s.slug,
			s.quantity_sold,
			s.gross_revenue,
			COALESCE(e.product_views, 0),
			COALESCE(e.add_to_cart, 0),
			COALESCE(e.wishlist_adds, 0)
		FROM sales s
		LEFT JOIN events e ON e.product_id = s.id
		ORDER BY s.quantity_sold DESC, s.gross_revenue DESC, COALESCE(e.product_views, 0) DESC
		LIMIT $4
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, vendorID, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []TopProductAnalytics{}
	for rows.Next() {
		var product TopProductAnalytics
		if err := rows.Scan(&product.ProductID, &product.Name, &product.Slug, &product.QuantitySold, &product.GrossRevenue, &product.ProductViews, &product.AddToCart, &product.WishlistAdds); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, rows.Err()
}

func (s *AnalyticsStore) GetOrderStatus(ctx context.Context, vendorID string, from, to *time.Time) ([]StatusCount, error) {
	vendorClause := ""
	args := []any{from, to}
	if vendorID != "" {
		vendorClause = "AND oi.vendor_id = $3"
		args = append(args, vendorID)
	}

	query := fmt.Sprintf(`
		SELECT o.status, COUNT(DISTINCT o.id)
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id
		WHERE ($1::timestamptz IS NULL OR o.created_at >= $1)
			AND ($2::timestamptz IS NULL OR o.created_at <= $2)
			%s
		GROUP BY o.status
		ORDER BY o.status
	`, vendorClause)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statuses := []StatusCount{}
	for rows.Next() {
		var status StatusCount
		if err := rows.Scan(&status.Status, &status.Count); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	return statuses, rows.Err()
}

func (s *AnalyticsStore) GetCustomerAnalytics(ctx context.Context, vendorID string, from, to *time.Time) (*CustomerAnalytics, error) {
	query := `
		WITH customer_orders AS (
			SELECT o.user_id, COUNT(DISTINCT o.id) AS orders_count
			FROM orders o
			JOIN order_items oi ON oi.order_id = o.id
			WHERE oi.vendor_id = $1
				AND ($2::timestamptz IS NULL OR o.created_at >= $2)
				AND ($3::timestamptz IS NULL OR o.created_at <= $3)
			GROUP BY o.user_id
		)
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE orders_count > 1)
		FROM customer_orders
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	customers := &CustomerAnalytics{}
	err := s.db.QueryRowContext(ctx, query, vendorID, from, to).Scan(&customers.TotalCustomers, &customers.RepeatCustomers)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (s *AnalyticsStore) GetTopVendors(ctx context.Context, from, to *time.Time, limit int) ([]TopVendorAnalytics, error) {
	query := `
		SELECT
			v.user_id,
			v.store_name,
			v.subdomain,
			COUNT(DISTINCT o.id) AS orders_count,
			COALESCE(SUM(oi.subtotal) FILTER (WHERE o.id IS NOT NULL), 0) AS gross_revenue,
			COALESCE(SUM(oi.vendor_payout) FILTER (WHERE o.id IS NOT NULL), 0) AS vendor_payout,
			COALESCE(SUM(oi.quantity) FILTER (WHERE o.id IS NOT NULL), 0) AS products_sold
		FROM vendor_profiles v
		LEFT JOIN order_items oi ON oi.vendor_id = v.user_id
		LEFT JOIN orders o ON o.id = oi.order_id
			AND ($1::timestamptz IS NULL OR o.created_at >= $1)
			AND ($2::timestamptz IS NULL OR o.created_at <= $2)
		GROUP BY v.user_id, v.store_name, v.subdomain
		ORDER BY gross_revenue DESC, orders_count DESC
		LIMIT $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vendors := []TopVendorAnalytics{}
	for rows.Next() {
		var vendor TopVendorAnalytics
		if err := rows.Scan(&vendor.VendorID, &vendor.StoreName, &vendor.Subdomain, &vendor.OrdersCount, &vendor.GrossRevenue, &vendor.VendorPayout, &vendor.ProductsSold); err != nil {
			return nil, err
		}
		vendors = append(vendors, vendor)
	}

	return vendors, rows.Err()
}

func (s *AnalyticsStore) GetUserAnalytics(ctx context.Context, from, to *time.Time) (*UserAnalytics, error) {
	query := `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'active'),
			COUNT(*) FILTER (WHERE status = 'suspended'),
			COUNT(*) FILTER (WHERE role = 'user'),
			COUNT(*) FILTER (WHERE role = 'vendor'),
			COUNT(*) FILTER (WHERE role = 'admin')
		FROM users
		WHERE ($1::timestamptz IS NULL OR created_at >= $1)
			AND ($2::timestamptz IS NULL OR created_at <= $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	users := &UserAnalytics{}
	err := s.db.QueryRowContext(ctx, query, from, to).Scan(&users.TotalUsers, &users.ActiveUsers, &users.SuspendedUsers, &users.Customers, &users.Vendors, &users.Admins)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *AnalyticsStore) GetVendorApplicationAnalytics(ctx context.Context, from, to *time.Time) (*VendorApplicationAnalytics, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending'),
			COUNT(*) FILTER (WHERE status = 'approved'),
			COUNT(*) FILTER (WHERE status = 'rejected')
		FROM role_upgrade_requests
		WHERE ($1::timestamptz IS NULL OR created_at >= $1)
			AND ($2::timestamptz IS NULL OR created_at <= $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	applications := &VendorApplicationAnalytics{}
	err := s.db.QueryRowContext(ctx, query, from, to).Scan(&applications.Pending, &applications.Approved, &applications.Rejected)
	if err != nil {
		return nil, err
	}

	return applications, nil
}
