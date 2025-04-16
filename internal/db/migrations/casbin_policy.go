package migrations

import "github.com/gocql/gocql"

func CreateCasbinPolicyTable(session *gocql.Session) error {
	query := `
	CREATE TABLE IF NOT EXISTS casbin_policy (
		no TEXT PRIMARY KEY,
		ptype TEXT,
		v1 TEXT,
		v2 TEXT,
		v3 TEXT,
		v4 TEXT
	)`
	return session.Query(query).Exec()
}
