package main

import (
	"net/http"

	"github.com/Davidmuthee12/eazymarket/internals/store"
)

type UpdateVendorProfile struct {
	Storename      string `json:"storename" validate:"required,max=150"`
	Subdomain      string `json:"subdomain" validate:"required,max=100"`
	Description    string `json:"description" validate:"required,max=250"`
	Logo_URL       string `json:"logo_url" validate:"omitempty,max=100"`
	Banner_URL     string `json:"banner_url" validate:"omitempty,max=100"`
	Business_Email string `json:"business_email" validate:"omitempty,max=255"`
	Business_Phone string `json:"business_phone" validate:"omitempty,max=20"`
	Address        string `json:"address" validate:"omitempty,max=100"`
}

// CreateVendorProfile godoc
//
//	@Summary		Creates a vendor profile
//	@Description	Creates a vendor profile for the authenticated vendor user
//	@Tags			vendor
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UpdateVendorProfile	true	"Vendor profile payload"
//	@Success		201		{object}	store.Vendor		"Vendor profile created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/profile [post]
func (app *application) vendorProfileHandler(w http.ResponseWriter, r *http.Request) {
	var Payload UpdateVendorProfile

	if err := readJSON(w, r, &Payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(Payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromCtx(r)

	vendor := &store.Vendor{
		Storename:      Payload.Storename,
		Subdomain:      Payload.Subdomain,
		Description:    Payload.Description,
		Logo_URL:       Payload.Logo_URL,
		Banner_URL:     Payload.Banner_URL,
		Business_Email: Payload.Business_Email,
		Business_Phone: Payload.Business_Phone,
		Address:        Payload.Address,
	}

	ctx := r.Context()

	// store the vendor profile
	err := app.store.Vendor.CreateVendorProfile(ctx, vendor, user.UUID)
	if err != nil {
		switch err {
		case store.ErrDuplicateStoreName:
			app.badRequestResponse(w, r, err)
		case store.ErrDuplicateSubdomain:
			app.badRequestResponse(w, r, err)
		case store.ErrDuplicateBusinessEmail:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, vendor); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetVendorProfile godoc
//
//	@Summary		Get vendor profile
//	@Description	Retrieves the vendor profile for the authenticated vendor user
//	@Tags			vendor
//	@Produce		json
//	@Success		200	{object}	store.Vendor	"Vendor profile retrieved"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/profile [get]
func (app *application) getVendorProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	ctx := r.Context()

	profile, err := app.store.Vendor.GetVendorByUUID(ctx, user.UUID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, profile); err != nil {
		app.internalServerError(w, r, err)
	}
}
