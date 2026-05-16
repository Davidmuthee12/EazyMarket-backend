package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Davidmuthee12/eazymarket/internals/store"
)

type AnalyticsEventPayload struct {
	SessionID string         `json:"session_id" validate:"omitempty,max=255"`
	EventType string         `json:"event_type" validate:"required,oneof=storefront_view product_view add_to_cart wishlist_add checkout_started order_created"`
	ProductID string         `json:"product_id" validate:"omitempty,uuid"`
	Metadata  map[string]any `json:"metadata" validate:"omitempty"`
}

func parseAnalyticsWindow(r *http.Request) (*time.Time, *time.Time, error) {
	query := r.URL.Query()

	var from *time.Time
	if value := query.Get("from"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, nil, err
		}
		from = &parsed
	}

	var to *time.Time
	if value := query.Get("to"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, nil, err
		}
		to = &parsed
	}

	return from, to, nil
}

func parseAnalyticsInterval(r *http.Request) string {
	switch r.URL.Query().Get("interval") {
	case "week":
		return "week"
	case "month":
		return "month"
	default:
		return "day"
	}
}

func parseAnalyticsLimit(r *http.Request, fallback int) int {
	limit := fallback
	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			limit = parsed
		}
	}

	if limit < 1 {
		return 1
	}
	if limit > 100 {
		return 100
	}
	return limit
}

// TrackStorefrontAnalyticsEvent godoc
//
//	@Summary		Track storefront analytics event
//	@Description	Records a storefront behavior event for the resolved vendor storefront. For local/dev clients, pass the vendor subdomain using X-Store-Subdomain or the store query parameter.
//	@Tags			storefront analytics
//	@Accept			json
//	@Produce		json
//	@Param			X-Store-Subdomain	header		string					false	"Vendor subdomain used when the request host is not a vendor subdomain"
//	@Param			store				query		string					false	"Vendor subdomain fallback for local/dev clients"
//	@Param			payload				body		AnalyticsEventPayload	true	"Analytics event payload"
//	@Success		201					{object}	store.AnalyticsEvent	"Analytics event tracked"
//	@Failure		400					{object}	error
//	@Failure		404					{object}	error
//	@Failure		500					{object}	error
//	@Router			/storefront/analytics/events [post]
func (app *application) trackStorefrontAnalyticsEventHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getStorefrontVendorFromCtx(r)
	if vendor == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	var payload AnalyticsEventPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(&payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	metadata, err := json.Marshal(payload.Metadata)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	event := &store.AnalyticsEvent{
		VendorID:  vendor.UserID,
		SessionID: payload.SessionID,
		EventType: payload.EventType,
		ProductID: payload.ProductID,
		Metadata:  metadata,
	}

	if err := app.store.Analytics.TrackEvent(r.Context(), event); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, event); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetVendorAnalyticsSummary godoc
//
//	@Summary		Get vendor analytics summary
//	@Description	Retrieves revenue, order, customer, and behavior metrics for the authenticated vendor.
//	@Tags			vendor analytics
//	@Produce		json
//	@Param			from	query		string	false	"Start datetime in RFC3339 format"
//	@Param			to		query		string	false	"End datetime in RFC3339 format"
//	@Success		200		{object}	store.AnalyticsSummary
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/analytics/summary [get]
func (app *application) getVendorAnalyticsSummaryHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getUserFromCtx(r)
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	summary, err := app.store.Analytics.GetVendorSummary(r.Context(), vendor.UUID, from, to)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, summary); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetVendorAnalyticsRevenue godoc
//
//	@Summary		Get vendor revenue analytics
//	@Description	Retrieves vendor revenue over time grouped by day, week, or month.
//	@Tags			vendor analytics
//	@Produce		json
//	@Param			from		query		string	false	"Start datetime in RFC3339 format"
//	@Param			to			query		string	false	"End datetime in RFC3339 format"
//	@Param			interval	query		string	false	"Grouping interval: day, week, or month"
//	@Success		200			{array}		store.RevenuePoint
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/analytics/revenue [get]
func (app *application) getVendorAnalyticsRevenueHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getUserFromCtx(r)
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	points, err := app.store.Analytics.GetRevenue(r.Context(), vendor.UUID, from, to, parseAnalyticsInterval(r))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, points); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetVendorTopProducts godoc
//
//	@Summary		Get vendor top products
//	@Description	Retrieves top vendor products by sales and behavior events.
//	@Tags			vendor analytics
//	@Produce		json
//	@Param			from	query		string	false	"Start datetime in RFC3339 format"
//	@Param			to		query		string	false	"End datetime in RFC3339 format"
//	@Param			limit	query		int		false	"Maximum number of products, capped at 100"
//	@Success		200		{array}		store.TopProductAnalytics
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/analytics/products/top [get]
func (app *application) getVendorTopProductsHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getUserFromCtx(r)
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	products, err := app.store.Analytics.GetTopProducts(r.Context(), vendor.UUID, from, to, parseAnalyticsLimit(r, 10))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, products); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetVendorOrderStatusAnalytics godoc
//
//	@Summary		Get vendor order status analytics
//	@Description	Retrieves order counts by status for the authenticated vendor.
//	@Tags			vendor analytics
//	@Produce		json
//	@Param			from	query		string	false	"Start datetime in RFC3339 format"
//	@Param			to		query		string	false	"End datetime in RFC3339 format"
//	@Success		200		{array}		store.StatusCount
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/analytics/orders/status [get]
func (app *application) getVendorOrderStatusAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getUserFromCtx(r)
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	statuses, err := app.store.Analytics.GetOrderStatus(r.Context(), vendor.UUID, from, to)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, statuses); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetVendorCustomerAnalytics godoc
//
//	@Summary		Get vendor customer analytics
//	@Description	Retrieves total and repeat customer counts for the authenticated vendor.
//	@Tags			vendor analytics
//	@Produce		json
//	@Param			from	query		string	false	"Start datetime in RFC3339 format"
//	@Param			to		query		string	false	"End datetime in RFC3339 format"
//	@Success		200		{object}	store.CustomerAnalytics
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/analytics/customers [get]
func (app *application) getVendorCustomerAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getUserFromCtx(r)
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	customers, err := app.store.Analytics.GetCustomerAnalytics(r.Context(), vendor.UUID, from, to)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, customers); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAdminAnalyticsSummary godoc
//
//	@Summary		Get admin analytics summary
//	@Description	Retrieves platform-wide revenue, order, user, vendor, and behavior metrics.
//	@Tags			admin analytics
//	@Produce		json
//	@Param			from	query		string	false	"Start datetime in RFC3339 format"
//	@Param			to		query		string	false	"End datetime in RFC3339 format"
//	@Success		200		{object}	store.AdminAnalyticsSummary
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/analytics/summary [get]
func (app *application) getAdminAnalyticsSummaryHandler(w http.ResponseWriter, r *http.Request) {
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	summary, err := app.store.Analytics.GetAdminSummary(r.Context(), from, to)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, summary); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAdminAnalyticsRevenue godoc
//
//	@Summary		Get admin revenue analytics
//	@Description	Retrieves platform revenue over time grouped by day, week, or month.
//	@Tags			admin analytics
//	@Produce		json
//	@Param			from		query		string	false	"Start datetime in RFC3339 format"
//	@Param			to			query		string	false	"End datetime in RFC3339 format"
//	@Param			interval	query		string	false	"Grouping interval: day, week, or month"
//	@Success		200			{array}		store.RevenuePoint
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/analytics/revenue [get]
func (app *application) getAdminAnalyticsRevenueHandler(w http.ResponseWriter, r *http.Request) {
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	points, err := app.store.Analytics.GetRevenue(r.Context(), "", from, to, parseAnalyticsInterval(r))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, points); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAdminTopVendors godoc
//
//	@Summary		Get admin top vendors
//	@Description	Retrieves top vendors by revenue and order volume.
//	@Tags			admin analytics
//	@Produce		json
//	@Param			from	query		string	false	"Start datetime in RFC3339 format"
//	@Param			to		query		string	false	"End datetime in RFC3339 format"
//	@Param			limit	query		int		false	"Maximum number of vendors, capped at 100"
//	@Success		200		{array}		store.TopVendorAnalytics
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/analytics/vendors/top [get]
func (app *application) getAdminTopVendorsHandler(w http.ResponseWriter, r *http.Request) {
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	vendors, err := app.store.Analytics.GetTopVendors(r.Context(), from, to, parseAnalyticsLimit(r, 10))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, vendors); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAdminTopProducts godoc
//
//	@Summary		Get admin top products
//	@Description	Retrieves platform-wide top products by sales and behavior events.
//	@Tags			admin analytics
//	@Produce		json
//	@Param			from	query		string	false	"Start datetime in RFC3339 format"
//	@Param			to		query		string	false	"End datetime in RFC3339 format"
//	@Param			limit	query		int		false	"Maximum number of products, capped at 100"
//	@Success		200		{array}		store.TopProductAnalytics
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/analytics/products/top [get]
func (app *application) getAdminTopProductsHandler(w http.ResponseWriter, r *http.Request) {
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	products, err := app.store.Analytics.GetTopProducts(r.Context(), "", from, to, parseAnalyticsLimit(r, 10))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, products); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAdminOrderStatusAnalytics godoc
//
//	@Summary		Get admin order status analytics
//	@Description	Retrieves platform-wide order counts by status.
//	@Tags			admin analytics
//	@Produce		json
//	@Param			from	query		string	false	"Start datetime in RFC3339 format"
//	@Param			to		query		string	false	"End datetime in RFC3339 format"
//	@Success		200		{array}		store.StatusCount
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/analytics/orders/status [get]
func (app *application) getAdminOrderStatusAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	statuses, err := app.store.Analytics.GetOrderStatus(r.Context(), "", from, to)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, statuses); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAdminUserAnalytics godoc
//
//	@Summary		Get admin user analytics
//	@Description	Retrieves platform user counts by status and role.
//	@Tags			admin analytics
//	@Produce		json
//	@Param			from	query		string	false	"Start datetime in RFC3339 format"
//	@Param			to		query		string	false	"End datetime in RFC3339 format"
//	@Success		200		{object}	store.UserAnalytics
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/analytics/users [get]
func (app *application) getAdminUserAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	users, err := app.store.Analytics.GetUserAnalytics(r.Context(), from, to)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, users); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAdminVendorApplicationAnalytics godoc
//
//	@Summary		Get admin vendor application analytics
//	@Description	Retrieves vendor application counts by status.
//	@Tags			admin analytics
//	@Produce		json
//	@Param			from	query		string	false	"Start datetime in RFC3339 format"
//	@Param			to		query		string	false	"End datetime in RFC3339 format"
//	@Success		200		{object}	store.VendorApplicationAnalytics
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/analytics/vendor-applications [get]
func (app *application) getAdminVendorApplicationAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	from, to, err := parseAnalyticsWindow(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	applications, err := app.store.Analytics.GetVendorApplicationAnalytics(r.Context(), from, to)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, applications); err != nil {
		app.internalServerError(w, r, err)
	}
}
