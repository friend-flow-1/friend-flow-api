package db

import (
	"log"

	"github.com/gocql/gocql"
)

func CreateIndexesScylla(session *gocql.Session) {
	queries := []string{
		`CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);`,
		// Add more index queries here
		// `CREATE INDEX IF NOT EXISTS another_idx ON table (column);`,
	}

	for _, q := range queries {
		if err := session.Query(q).Exec(); err != nil {
			log.Printf("❌ Failed to create index: %s\nError: %v", q, err)
		} else {
			log.Printf("✅ Successfully created index: %s", q)
		}
	}
}
