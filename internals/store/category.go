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

func (s *CategoryStore) GetCategories(ctx context.Context) ([]Category, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			image_url,
			created_at
		FROM categories
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []Category{}

	for rows.Next() {
		var category Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Image_URL,
			&category.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, ErrNotFound
	}

	return categories, nil
}

func (s *CategoryStore) DeleteCategory(ctx context.Context, categoryID string) error {
	query := `
		DELETE FROM categories
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		categoryID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryStore) UpdateCategory(ctx context.Context, category *Category) error {
	query := `
		UPDATE categories
		SET name = $1, slug = $2, image_url = $3
		WHERE id = $4
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.Image_URL,
		category.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
