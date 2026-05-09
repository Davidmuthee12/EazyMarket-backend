package main

import (
	"net/http"

	"github.com/Davidmuthee12/eazymarket/internals/store"
)

type ProductPayload struct {
	Name           string  `json:"name" validate:"required,max=100"`
	Slug           string  `json:"slug" validate:"required,max=100"`
	Description    string  `json:"description" validate:"required,max=250"`
	Category_ID    string  `json:"category_id" validate:"omitempty,uuid"`
	Price          float64 `json:"price" validate:"required"`
	Compare_Price  float64 `json:"compare_price" validate:"omitempty"`
	Stock_Quantity int     `json:"stock_quantity" validate:"required"`
	SKU            string  `json:"sku" validate:"required,max=50"`
	Weight         float64 `json:"weight" validate:"omitempty"`
}

// CreateProduct godoc
//
//	@Summary		Create a product
//	@Description	Creates a new product for the authenticated vendor.
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		ProductPayload	true	"Product payload"
//	@Success		201		{object}	store.Products	"Product created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/products [post]
func (app *application) postProductsHandler(w http.ResponseWriter, r *http.Request) {
	var Payload ProductPayload
	vendor := getUserFromCtx(r)
	if vendor == nil {
		app.internalServerError(w, r, nil)
		return
	}

	if err := readJSON(w, r, &Payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(&Payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product := &store.Products{
		Name:           Payload.Name,
		Slug:           Payload.Slug,
		Description:    Payload.Description,
		Category_ID:    Payload.Category_ID,
		Price:          Payload.Price,
		Compare_Price:  Payload.Compare_Price,
		Stock_Quantity: Payload.Stock_Quantity,
		SKU:            Payload.SKU,
		Weight:         Payload.Weight,
	}

	ctx := r.Context()

	// store the product
	err := app.store.Product.CreateProduct(ctx, product, vendor.UUID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, err)
			return
		case store.ErrDuplicateProductName, store.ErrDuplicateProductSlug:
			app.badRequestResponse(w, r, err)
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

// GetAllProducts godoc
//
//	@Summary		Get all products
//	@Description	Retrieves all products for the authenticated vendor.
//	@Tags			products
//	@Produce		json
//	@Success		200	{array}		store.Products	"Products retrieved"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/vendor/products [get]
func (app *application) getAllProducts(w http.ResponseWriter, r *http.Request) {
	vendor := getUserFromCtx(r)
	if vendor == nil {
		app.internalServerError(w, r, nil)
		return
	}

	ctx := r.Context()

	products, err := app.store.Product.GetAllProduct(ctx, vendor.UUID)
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

	if err := app.jsonResponse(w, http.StatusOK, products); err != nil {
		app.internalServerError(w, r, err)
	}
}
