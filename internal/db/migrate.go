package db

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/haxxu/friend-flow-api/internal/db/migrations"
)

func AutoMigrate(session *gocql.Session) {
	migrations := []func(*gocql.Session) error{
		migrations.CreateUsersTable,
	}

	for _, migrate := range migrations {
		if err := migrate(session); err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}
	}

	log.Println("✅ All tables migrated successfully")
}
