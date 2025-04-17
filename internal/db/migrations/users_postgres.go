package migrations

import (
	"database/sql"
)

func CreateUsersTablePostgres(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		phone TEXT,
		first_name TEXT,
		last_name TEXT,
		password TEXT NOT NULL,
		status TEXT,
		background TEXT,
		avatar TEXT,
		birth_date TIMESTAMP,
		gender TEXT,
		role TEXT,
		created_at TIMESTAMP,
		created_by TEXT,
		updated_at TIMESTAMP,
		updated_by TEXT,
		deleted_at TIMESTAMP,
		deleted_by TEXT
	);`
	_, err := db.Exec(query)
	return err
}
