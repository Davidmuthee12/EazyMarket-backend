package main

import (
	"net/http"

	"github.com/Davidmuthee12/eazymarket/internals/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AddCategory struct {
	Name      string `json:"name" validate:"required,max=100"`
	Slug      string `json:"slug" validate:"required,max=100"`
	Image_URL string `json:"image_url" validate:"omitempty,max=100"`
}

// AddCategory godoc
//
//	@Summary		Add a category
//	@Description	Creates a new product category. Requires admin privileges.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		AddCategory		true	"Category payload"
//	@Success		201		{object}	store.Category	"Category created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/categories [post]
func (app *application) addCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	var Payload AddCategory

	if err := readJSON(w, r, &Payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(&Payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category := &store.Category{
		Name:      Payload.Name,
		Slug:      Payload.Slug,
		Image_URL: Payload.Image_URL,
	}

	ctx := r.Context()

	// store the category
	err := app.store.Category.AddCategory(ctx, category)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusCreated, category); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetCategories godoc
//
//	@Summary		Get all categories
//	@Description	Retrieves all product categories
//	@Tags			categories
//	@Produce		json
//	@Success		200	{array}		store.Category	"Categories retrieved"
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/categories [get]
func (app *application) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := app.store.Category.GetCategories(ctx)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, categories); err != nil {
		app.internalServerError(w, r, err)
	}
}

// DeleteCategory godoc
//
//	@Summary		Delete a category
//	@Description	Deletes a product category by ID. Requires admin privileges.
//	@Tags			categories
//	@Param			id	path		string	true	"Category ID"
//	@Success		200	{object}	nil
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/categories/{id} [delete]
func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(categoryID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	err := app.store.Category.DeleteCategory(ctx, categoryID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, err)
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

// UpdateCategory godoc
//
//	@Summary		Update a category
//	@Description	Updates an existing product category. Requires admin privileges.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Category ID"
//	@Param			payload	body		AddCategory		true	"Category update payload"
//	@Success		200		{object}	store.Category	"Category updated"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/admin/categories/{id} [put]
func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(categoryID); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var payload AddCategory

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(&payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category := &store.Category{
		ID:        categoryID,
		Name:      payload.Name,
		Slug:      payload.Slug,
		Image_URL: payload.Image_URL,
	}

	ctx := r.Context()

	err := app.store.Category.UpdateCategory(ctx, category)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, category); err != nil {
		app.internalServerError(w, r, err)
	}
}
