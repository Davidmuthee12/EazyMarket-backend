package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrDuplicateCategoryName = errors.New("a category with a similar name already exists")
	ErrDuplicateSlug         = errors.New("a category with a similar slug name exists")
)

type Category struct {
	ID        string `json:"-"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Parent_ID string `json:"parent_id"`
	Image_URL string `json:"image_url"`
	CreatedAt string `json:"created_at"`
}

type CategoryStore struct {
	db *sql.DB
}

func (s *CategoryStore) AddCategory(ctx context.Context, category *Category) error {
	query := `
		INSERT INTO categories(
			name,
			slug,
			image_url
		)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.Image_URL,
	).Scan(
		&category.ID,
		&category.CreatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "category_name_key":
				return ErrDuplicateCategoryName
			case "category_slug_key":
				return ErrDuplicateSlug
			}
		}
		return err
	}

	return nil
}
