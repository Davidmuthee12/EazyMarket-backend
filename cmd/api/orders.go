package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Davidmuthee12/eazymarket/internals/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateOrderPayload struct {
	ShippingAddress map[string]any `json:"shipping_address" validate:"omitempty"`
	Notes           string         `json:"notes" validate:"omitempty,max=500"`
}

type UpdateOrderStatusPayload struct {
	Status string `json:"status" validate:"required,oneof=confirmed shipped delivered cancelled refunded"`
}

// CreateOrder godoc
//
//	@Summary		Create an order
//	@Description	Creates an order from the authenticated user's active cart and clears the cart.
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateOrderPayload	false	"Order payload"
//	@Success		201		{object}	store.Order			"Order created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/orders [post]
func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}

	payload := CreateOrderPayload{}
	if r.Body != nil {
		if err := readJSON(w, r, &payload); err != nil {
			if !errors.Is(err, io.EOF) {
				app.badRequestResponse(w, r, err)
				return
			}
		}
	}

	if err := Validate.Struct(&payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	shippingAddress, err := json.Marshal(payload.ShippingAddress)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	order, err := app.store.Order.CreateFromCart(r.Context(), user.UUID, shippingAddress, payload.Notes)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrEmptyCart):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, order); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetOrders godoc
//
//	@Summary		Get orders
//	@Description	Retrieves the authenticated user's orders.
//	@Tags			orders
//	@Produce		json
//	@Success		200	{array}		store.Order	"Orders retrieved"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/orders [get]
func (app *application) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}

	orders, err := app.store.Order.GetAll(r.Context(), user.UUID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, orders); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetOrderByID godoc
//
//	@Summary		Get order by ID
//	@Description	Retrieves one order for the authenticated user.
//	@Tags			orders
//	@Produce		json
//	@Param			orderID	path		string		true	"Order ID"
//	@Success		200		{object}	store.Order	"Order retrieved"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/orders/{orderID} [get]
func (app *application) getOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}

	orderID := chi.URLParam(r, "orderID")
	if _, err := uuid.Parse(orderID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	order, err := app.store.Order.GetByID(r.Context(), user.UUID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, order); err != nil {
		app.internalServerError(w, r, err)
	}
}

// CancelOrder godoc
//
//	@Summary		Cancel an order
//	@Description	Cancels a pending or confirmed order for the authenticated user.
//	@Tags			orders
//	@Produce		json
//	@Param			orderID	path		string		true	"Order ID"
//	@Success		200		{object}	store.Order	"Order cancelled"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/orders/{orderID}/cancel [put]
func (app *application) cancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}

	orderID := chi.URLParam(r, "orderID")
	if _, err := uuid.Parse(orderID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	order, err := app.store.Order.Cancel(r.Context(), user.UUID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, order); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetVendorOrders godoc
//
//	@Summary		Get vendor orders
//	@Description	Retrieves orders containing products owned by the authenticated vendor.
//	@Tags			orders
//	@Produce		json
//	@Success		200	{array}		store.Order	"Vendor orders retrieved"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/orders [get]
func (app *application) getVendorOrdersHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getUserFromCtx(r)
	if vendor == nil {
		app.internalServerError(w, r, nil)
		return
	}

	orders, err := app.store.Order.GetVendorOrders(r.Context(), vendor.UUID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, orders); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetVendorOrderByID godoc
//
//	@Summary		Get vendor order by ID
//	@Description	Retrieves one order containing products owned by the authenticated vendor.
//	@Tags			orders
//	@Produce		json
//	@Param			orderID	path		string		true	"Order ID"
//	@Success		200		{object}	store.Order	"Vendor order retrieved"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/orders/{orderID} [get]
func (app *application) getVendorOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getUserFromCtx(r)
	if vendor == nil {
		app.internalServerError(w, r, nil)
		return
	}

	orderID := chi.URLParam(r, "orderID")
	if _, err := uuid.Parse(orderID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	order, err := app.store.Order.GetVendorOrderByID(r.Context(), vendor.UUID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, order); err != nil {
		app.internalServerError(w, r, err)
	}
}

// UpdateVendorOrderStatus godoc
//
//	@Summary		Update vendor order status
//	@Description	Updates the status of an order containing products owned by the authenticated vendor.
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			orderID	path		string						true	"Order ID"
//	@Param			payload	body		UpdateOrderStatusPayload	true	"Order status payload"
//	@Success		200		{object}	store.Order					"Order status updated"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/orders/{orderID}/status [put]
func (app *application) updateVendorOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getUserFromCtx(r)
	if vendor == nil {
		app.internalServerError(w, r, nil)
		return
	}

	orderID := chi.URLParam(r, "orderID")
	if _, err := uuid.Parse(orderID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var payload UpdateOrderStatusPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(&payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	order, err := app.store.Order.UpdateVendorOrderStatus(r.Context(), vendor.UUID, orderID, payload.Status)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, order); err != nil {
		app.internalServerError(w, r, err)
	}
}
