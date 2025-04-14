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
		full_name TEXT,
		password TEXT,
		status TEXT,
		background TEXT,
		avatar TEXT,
		birth_date TIMESTAMP,
		gender TEXT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP
	)`
	return session.Query(query).Exec()
}
