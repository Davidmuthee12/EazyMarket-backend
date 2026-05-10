package main

import (
	"net/http"

	"github.com/Davidmuthee12/eazymarket/internals/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CartItemPayload struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type CartQuantityPayload struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

// AddCartItem godoc
//
//	@Summary		Add item to cart
//	@Description	Adds a product to the authenticated user's active cart. Creates the cart automatically if it does not exist.
//	@Tags			cart
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CartItemPayload	true	"Cart item payload"
//	@Success		201		{object}	store.CartItem	"Cart item added"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/cart/items [post]
func (app *application) addCartItemHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}

	var payload CartItemPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(&payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	item, err := app.store.Cart.AddItem(ctx, user.UUID, payload.ProductID, payload.Quantity)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, item); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetCart godoc
//
//	@Summary		Get cart
//	@Description	Retrieves the authenticated user's active cart.
//	@Tags			cart
//	@Produce		json
//	@Success		200	{object}	store.Cart	"Cart retrieved"
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/cart [get]
func (app *application) getCartHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}

	cart, err := app.store.Cart.GetCart(r.Context(), user.UUID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, cart); err != nil {
		app.internalServerError(w, r, err)
	}
}

// UpdateCartItem godoc
//
//	@Summary		Update cart item
//	@Description	Updates the quantity of a product in the authenticated user's active cart.
//	@Tags			cart
//	@Accept			json
//	@Produce		json
//	@Param			productID	path		string				true	"Product ID"
//	@Param			payload		body		CartQuantityPayload	true	"Cart item quantity payload"
//	@Success		200			{object}	store.CartItem		"Cart item updated"
//	@Failure		400			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/cart/items/{productID} [put]
func (app *application) updateCartItemHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}

	productID := chi.URLParam(r, "productID")
	if _, err := uuid.Parse(productID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var payload CartQuantityPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(&payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	item, err := app.store.Cart.UpdateItem(r.Context(), user.UUID, productID, payload.Quantity)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, item); err != nil {
		app.internalServerError(w, r, err)
	}
}

// RemoveCartItem godoc
//
//	@Summary		Remove cart item
//	@Description	Removes a product from the authenticated user's active cart.
//	@Tags			cart
//	@Param			productID	path		string	true	"Product ID"
//	@Success		200			{object}	nil		"Cart item removed"
//	@Failure		400			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/cart/items/{productID} [delete]
func (app *application) removeCartItemHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}

	productID := chi.URLParam(r, "productID")
	if _, err := uuid.Parse(productID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err := app.store.Cart.RemoveItem(r.Context(), user.UUID, productID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// ClearCart godoc
//
//	@Summary		Clear cart
//	@Description	Removes all items from the authenticated user's active cart.
//	@Tags			cart
//	@Success		200	{object}	nil	"Cart cleared"
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/cart [delete]
func (app *application) clearCartHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}

	if err := app.store.Cart.ClearCart(r.Context(), user.UUID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
