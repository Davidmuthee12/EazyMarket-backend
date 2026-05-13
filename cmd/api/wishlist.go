package main

import (
	"errors"
	"net/http"

	"github.com/Davidmuthee12/eazymarket/internals/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// AddToWishlist godoc
//
//	@Summary		Add product to wishlist
//	@Description	Adds a published product from the resolved storefront vendor to the authenticated user's storefront wishlist.
//	@Tags			wishlist
//	@Produce		json
//	@Param			X-Store-Subdomain	header		string			false	"Vendor subdomain used when the request host is not a vendor subdomain"
//	@Param			store				query		string			false	"Vendor subdomain fallback for local/dev clients"
//	@Param			productID			path		string			true	"Product ID"
//	@Success		201					{object}	store.Wishlist	"Product added to wishlist"
//	@Failure		400					{object}	error
//	@Failure		401					{object}	error
//	@Failure		404					{object}	error
//	@Failure		409					{object}	error
//	@Failure		500					{object}	error
//	@Security		ApiKeyAuth
//	@Router			/wishlist/{productID} [post]
func (app *application) addToWishlistHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}
	vendor := getStorefrontVendorFromCtx(r)
	if vendor == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	productID := chi.URLParam(r, "productID")
	if _, err := uuid.Parse(productID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	product, err := app.store.Wishlist.AddToWishList(ctx, user.UUID, vendor.UserID, productID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictResponse(w, r, err)
			return
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusCreated, product); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetUserWishlist godoc
//
//	@Summary		Get wishlist
//	@Description	Retrieves the authenticated user's wishlist with product details for the resolved storefront vendor.
//	@Tags			wishlist
//	@Produce		json
//	@Param			X-Store-Subdomain	header		string			false	"Vendor subdomain used when the request host is not a vendor subdomain"
//	@Param			store				query		string			false	"Vendor subdomain fallback for local/dev clients"
//	@Success		200					{array}		store.Wishlist	"Wishlist retrieved"
//	@Failure		401					{object}	error
//	@Failure		500					{object}	error
//	@Security		ApiKeyAuth
//	@Router			/wishlist [get]
func (app *application) getUserWishlistHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}
	vendor := getStorefrontVendorFromCtx(r)
	if vendor == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	ctx := r.Context()

	wishlist, err := app.store.Wishlist.GetUserWishlist(ctx, user.UUID, vendor.UserID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, wishlist); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetWishlistByID godoc
//
//	@Summary		Get wishlist item
//	@Description	Retrieves one product from the authenticated user's wishlist for the resolved storefront vendor.
//	@Tags			wishlist
//	@Produce		json
//	@Param			X-Store-Subdomain	header		string			false	"Vendor subdomain used when the request host is not a vendor subdomain"
//	@Param			store				query		string			false	"Vendor subdomain fallback for local/dev clients"
//	@Param			productID			path		string			true	"Product ID"
//	@Success		200					{object}	store.Wishlist	"Wishlist item retrieved"
//	@Failure		400					{object}	error
//	@Failure		401					{object}	error
//	@Failure		404					{object}	error
//	@Failure		500					{object}	error
//	@Security		ApiKeyAuth
//	@Router			/wishlist/{productID} [get]
func (app *application) getWishlistByIDHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}
	vendor := getStorefrontVendorFromCtx(r)
	if vendor == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	productID := chi.URLParam(r, "productID")
	if _, err := uuid.Parse(productID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	wishlist, err := app.store.Wishlist.GetWishlistByID(ctx, user.UUID, vendor.UserID, productID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, wishlist); err != nil {
		app.internalServerError(w, r, err)
	}
}

// DeleteWishlistItem godoc
//
//	@Summary		Delete product from wishlist
//	@Description	Removes a product from the authenticated user's wishlist for the resolved storefront vendor.
//	@Tags			wishlist
//	@Param			X-Store-Subdomain	header		string	false	"Vendor subdomain used when the request host is not a vendor subdomain"
//	@Param			store				query		string	false	"Vendor subdomain fallback for local/dev clients"
//	@Param			productID			path		string	true	"Product ID"
//	@Success		200					{object}	nil		"Product removed from wishlist"
//	@Failure		400					{object}	error
//	@Failure		401					{object}	error
//	@Failure		404					{object}	error
//	@Failure		500					{object}	error
//	@Security		ApiKeyAuth
//	@Router			/wishlist/{productID} [delete]
func (app *application) deleteWishlistItemHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if user == nil {
		app.internalServerError(w, r, nil)
		return
	}
	vendor := getStorefrontVendorFromCtx(r)
	if vendor == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	productID := chi.URLParam(r, "productID")
	if _, err := uuid.Parse(productID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	err := app.store.Wishlist.DeleteFromWishlist(ctx, user.UUID, vendor.UserID, productID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
