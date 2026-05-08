package main

import (
	"net/http"

	"github.com/Davidmuthee12/eazymarket/internals/store"
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
