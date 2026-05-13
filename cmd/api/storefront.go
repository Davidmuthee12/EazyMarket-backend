package main

import (
	"net/http"

	"github.com/Davidmuthee12/eazymarket/internals/store"
	"github.com/go-chi/chi/v5"
)

func getStorefrontVendorFromCtx(r *http.Request) *store.Vendor {
	vendor, _ := r.Context().Value(storefrontVendorCtx).(*store.Vendor)
	return vendor
}

// GetStorefront godoc
//
//	@Summary		Get storefront profile
//	@Description	Retrieves the approved vendor storefront resolved from the request subdomain. For local/dev clients, pass the vendor subdomain using X-Store-Subdomain or the store query parameter.
//	@Tags			storefront
//	@Produce		json
//	@Param			X-Store-Subdomain	header		string			false	"Vendor subdomain used when the request host is not a vendor subdomain"
//	@Param			store				query		string			false	"Vendor subdomain fallback for local/dev clients"
//	@Success		200					{object}	store.Vendor	"Storefront profile retrieved"
//	@Failure		404					{object}	error			"Storefront not found"
//	@Failure		500					{object}	error
//	@Router			/storefront [get]
func (app *application) getStorefrontHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getStorefrontVendorFromCtx(r)
	if vendor == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, vendor); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetStorefrontProducts godoc
//
//	@Summary		Get storefront products
//	@Description	Retrieves published products for the approved vendor storefront resolved from the request subdomain. For local/dev clients, pass the vendor subdomain using X-Store-Subdomain or the store query parameter.
//	@Tags			storefront
//	@Produce		json
//	@Param			X-Store-Subdomain	header		string			false	"Vendor subdomain used when the request host is not a vendor subdomain"
//	@Param			store				query		string			false	"Vendor subdomain fallback for local/dev clients"
//	@Success		200					{array}		store.Products	"Published storefront products retrieved"
//	@Failure		404					{object}	error			"Storefront not found"
//	@Failure		500					{object}	error
//	@Router			/storefront/products [get]
func (app *application) getStorefrontProductsHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getStorefrontVendorFromCtx(r)
	if vendor == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	products, err := app.store.Product.GetPublishedProductsByVendor(r.Context(), vendor.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, products); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetStorefrontProductBySlug godoc
//
//	@Summary		Get storefront product by slug
//	@Description	Retrieves one published product from the approved vendor storefront resolved from the request subdomain. For local/dev clients, pass the vendor subdomain using X-Store-Subdomain or the store query parameter.
//	@Tags			storefront
//	@Produce		json
//	@Param			X-Store-Subdomain	header		string			false	"Vendor subdomain used when the request host is not a vendor subdomain"
//	@Param			store				query		string			false	"Vendor subdomain fallback for local/dev clients"
//	@Param			slug				path		string			true	"Product slug"
//	@Success		200					{object}	store.Products	"Published storefront product retrieved"
//	@Failure		404					{object}	error			"Storefront or product not found"
//	@Failure		500					{object}	error
//	@Router			/storefront/products/{slug} [get]
func (app *application) getStorefrontProductBySlugHandler(w http.ResponseWriter, r *http.Request) {
	vendor := getStorefrontVendorFromCtx(r)
	if vendor == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	slug := chi.URLParam(r, "slug")
	product, err := app.store.Product.GetPublishedProductBySlug(r.Context(), vendor.UserID, slug)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, product); err != nil {
		app.internalServerError(w, r, err)
	}
}
