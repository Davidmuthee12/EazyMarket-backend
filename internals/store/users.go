package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type password struct {
	text *string
	hash []byte
}

type User struct {
	ID         int64    `json:"-"`
	UUID       string   `json:"id"`
	UserName   string   `json:"username"`
	Phone      int64    `json:"phone"`
	Avatar_Url string   `json:"avatar_url"`
	RoleID     int64    `json:"role_id"`
	Role       Role     `json:"role"`
	Email      string   `json:"email"`
	Password   password `json:"-"`
	CreatedAt  string   `json:"created_at"`
	IsActive   bool     `json:"is_active"`
	Status     string   `json:"status"`
}

type UserStore struct {
	db *sql.DB
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users(username, email, phone, password, role_id)
		VALUES ($1, $2, $3, $4, (SELECT id FROM roles WHERE name = $5))
		RETURNING id, uuid, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	role := user.Role.Name
	if role == "" {
		role = "user"
	}

	err := tx.QueryRowContext(
		ctx,
		query,
		user.UserName,
		user.Email,
		user.Phone,
		user.Password.hash,
		role,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.CreatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_email_key":
				return ErrDuplicateEmail
			case "users_username_key":
				return ErrDuplicateUsername
			}
		}

		return err
	}

	return nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	// transaction wrapper
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// create the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// create the user invite
		if err := s.CreateUserInvitation(ctx, tx, token, invitationExp, user.UUID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) CreateUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userUUID string) error {
	query := `
		INSERT INTO user_invitation (token, user_uuid, expiry) VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userUUID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT users.id, users.uuid, username, email, password, created_at, roles.*
		FROM users
		JOIN roles ON  (users.role_id = roles.id)
		WHERE users.id = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.UserName,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) GetByUUID(ctx context.Context, userUUID string) (*User, error) {
	query := `
		SELECT users.id, users.uuid, username, email, password, created_at, roles.*
		FROM users
		JOIN roles ON  (users.role_id = roles.id)
		WHERE users.uuid = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userUUID,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.UserName,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, uuid, username, email, password, created_at FROM users
		WHERE email = $1 AND is_active = true
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.UUID,
		&user.UserName,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// 1. find the user this token belongs to
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2. Update the status of user
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}
		// 3. clean/deleting invitations
		if err := s.deleteUserInvitations(ctx, tx, user.UUID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.GetByID(ctx, userID)
		if err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, user.UUID); err != nil {
			return err
		}

		if err := s.delete(ctx, tx, userID); err != nil {
			return err
		}
		return nil

	})
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
	 SELECT u.id, u.uuid, u.username, u.email, u.created_at, u.is_active
	 FROM users u
	 JOIN user_invitation ui on u.uuid = ui.user_uuid
	 WHERE ui.token = $1 AND ui.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.UUID,
		&user.UserName,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id  = $4`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.UserName, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userUUID string) error {
	query := `DELETE FROM user_invitation WHERE user_uuid = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userUUID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetAllUsers(ctx context.Context) ([]*User, error) {
	query := `
		SELECT
			u.id,
			u.uuid,
			u.username,
			u.email,
			COALESCE(u.avatar_url, ''),
			u.role_id,
			u.is_active,
			u.created_at,
			r.id,
			r.name,
			r.level,
			COALESCE(r.description, '')
		FROM users u
		JOIN roles r ON u.role_id = r.id
		ORDER BY u.id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID,
			&user.UUID,
			&user.UserName,
			&user.Email,
			&user.Avatar_Url,
			&user.RoleID,
			&user.IsActive,
			&user.CreatedAt,
			&user.Role.ID,
			&user.Role.Name,
			&user.Role.Level,
			&user.Role.Description,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNotFound
	}

	return users, nil
}

func (s *UserStore) UpdateRole(ctx context.Context, userUUID string) error {
	query := `
		INSERT INTO role_upgrade_requests (user_id, requested_role_id)
		VALUES ($1, (SELECT id FROM roles WHERE name = 'vendor'))
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userUUID)
	return err
}

func (s *UserStore) UpdateRoleRequest(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET role = 'vendor', role_id = (SELECT id FROM roles WHERE name = 'vendor')
		WHERE uuid = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UserStore) GetUpgradeRequests(ctx context.Context) ([]*User, error) {
	query := `
		SELECT
			u.id,
			u.uuid,
			u.username,
			u.email,
			COALESCE(u.avatar_url, ''),
			u.role_id,
			u.is_active,
			u.created_at AS joined,
			r.created_at AS requested_at,
			r.status
		FROM users u
		JOIN role_upgrade_requests r ON u.uuid = r.user_id
		ORDER BY u.id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		var requestedAt time.Time

		err := rows.Scan(
			&user.ID,
			&user.UUID,
			&user.UserName,
			&user.Email,
			&user.Avatar_Url,
			&user.RoleID,
			&user.IsActive,
			&user.CreatedAt,
			&requestedAt,
			&user.Status,
		)
		if err != nil {
			return nil, err
		}

		requests = append(requests, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(requests) == 0 {
		return nil, ErrNotFound
	}

	return requests, nil
}
