package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("record not found")
	ErrConflict          = errors.New("resource already exists")
	ErrEmptyCart         = errors.New("cart is empty")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users interface {
		GetAllUsers(context.Context) ([]*User, error)
		GetByID(context.Context, int64) (*User, error)
		GetByUUID(context.Context, string) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Create(context.Context, *sql.Tx, *User) error
		CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
		UpdateRole(context.Context, string) error
		GetUpgradeRequests(context.Context) ([]*User, error)
		UpdateRoleRequest(ctx context.Context, userID, reviewerID string) error
		RejectRequest(ctx context.Context, userID, reviewerID string) error
	}
	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}

	Vendor interface {
		CreateVendorProfile(ctx context.Context, Vendor *Vendor, userUUID string) error
		GetVendorByUUID(ctx context.Context, userID string) (*Vendor, error)
	}

	Category interface {
		AddCategory(ctx context.Context, category *Category) error
		GetCategories(context.Context) ([]Category, error)
		DeleteCategory(ctx context.Context, categoryID string) error
		UpdateCategory(ctx context.Context, category *Category) error
	}

	Product interface {
		CreateProduct(ctx context.Context, product *Products, vendorID string) error
		GetAllProduct(ctx context.Context, vendorID string) ([]Products, error)
		GetProductByUUID(ctx context.Context, productID string) (*Products, error)
		UpdateProduct(ctx context.Context, product *Products, vendorID string) error
		DeleteProduct(ctx context.Context, productID string, vendorID string) error
	}

	Cart interface {
		AddItem(ctx context.Context, userID, productID string, quantity int) (*CartItem, error)
		GetCart(ctx context.Context, userID string) (*Cart, error)
		UpdateItem(ctx context.Context, userID, productID string, quantity int) (*CartItem, error)
		RemoveItem(ctx context.Context, userID, productID string) error
		ClearCart(ctx context.Context, userID string) error
	}

	Order interface {
		CreateFromCart(ctx context.Context, userID string, shippingAddress []byte, notes string) (*Order, error)
		GetAll(ctx context.Context, userID string) ([]Order, error)
		GetByID(ctx context.Context, userID, orderID string) (*Order, error)
		Cancel(ctx context.Context, userID, orderID string) (*Order, error)
		GetVendorOrders(ctx context.Context, vendorID string) ([]Order, error)
		GetVendorOrderByID(ctx context.Context, vendorID, orderID string) (*Order, error)
		UpdateVendorOrderStatus(ctx context.Context, vendorID, orderID, status string) (*Order, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users:    &UserStore{db},
		Roles:    &RoleStore{db},
		Vendor:   &VenderStore{db},
		Category: &CategoryStore{db},
		Product:  &ProductStore{db},
		Cart:     &CartStore{db},
		Order:    &OrderStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
