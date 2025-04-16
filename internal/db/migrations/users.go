package migrations

import "github.com/gocql/gocql"

func CreateUsersTable(session *gocql.Session) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email TEXT,
		phone TEXT,
		first_name TEXT,
		last_name TEXT,
		password TEXT,
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
	)`
	return session.Query(query).Exec()
}
