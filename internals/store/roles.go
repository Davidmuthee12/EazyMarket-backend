package store

import "database/sql"

type Role struct {
	ID          int64  `json:"role_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type RoleStore struct {
	db *sql.DB
}
