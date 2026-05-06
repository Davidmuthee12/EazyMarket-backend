package main

import (
	"net/http"

	"github.com/Davidmuthee12/eazymarket/internals/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userUUID	path		string	true	"user UUID"
//	@Success		200			{object}	store.User
//	@Failure		400			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userUUID}/ [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userUUID := chi.URLParam(r, "userUUID")
	if _, err := uuid.Parse(userUUID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := app.getUser(ctx, userUUID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetAllUsers godoc
//
//	@Summary		Fetches all users
//	@Description	Fetches all users
//	@Tags			users
//	@Produce		json
//	@Success		200	{array}		store.User
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/ [get]
func (app *application) getAllUsersHandlers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := app.store.Users.GetAllUsers(ctx)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, users); err != nil {
		app.internalServerError(w, r, err)
	}
}

// UpdateRole godoc
//
//	@Summary		Request a vendor role upgrade
//	@Description	Submits a role upgrade request to vendor for the given user
//	@Tags			users
//	@Produce		json
//	@Param			userUUID	path	string	true	"user UUID"
//	@Success		200			"Request submitted"
//	@Failure		400			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userUUID}/upgrade-to-vendor [post]
func (app *application) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	userUUID := chi.URLParam(r, "userUUID")
	if _, err := uuid.Parse(userUUID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := app.getUser(ctx, userUUID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.store.Users.UpdateRole(ctx, user.UUID); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}

}

// GetVendorRequests godoc
//
//	@Summary		Fetches vendor upgrade requests
//	@Description	Fetches all user requests to upgrade role to vendor
//	@Tags			users
//	@Produce		json
//	@Success		200	{array}		store.User
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/vendor-request [get]
func (app *application) vendorRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requests, err := app.store.Users.GetUpgradeRequests(ctx)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, requests); err != nil {
		app.internalServerError(w, r, err)
	}
}

// ApproveVendor godoc
//
//	@Summary		Approves a vendor upgrade request
//	@Description	Approves a pending vendor role upgrade request for the given user
//	@Tags			admin
//	@Produce		json
//	@Param			userUUID	path	string	true	"user UUID"
//	@Success		200			"Vendor request approved"
//	@Failure		400			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/vendor-request/{userUUID}/approve [put]
func (app *application) approveVendorHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userUUID")
	if _, err := uuid.Parse(userID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	// Get the authenticated admin user from context
	reviewer := getUserFromCtx(r)
	if reviewer == nil {
		app.internalServerError(w, r, nil)
		return
	}

	err := app.store.Users.UpdateRoleRequest(ctx, userID, reviewer.UUID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}

}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
