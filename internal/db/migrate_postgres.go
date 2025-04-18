package db

import (
	"log"

	"github.com/haxxu/friend-flow-api/internal/modules/user"
	"gorm.io/gorm"
)

func AutoMigratePostgres(db *gorm.DB) error {
	// ✅ Enable uuid-ossp
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		log.Fatalf("❌ Failed to enable uuid-ossp extension: %v", err)
	}

	err := db.AutoMigrate(
		&user.User{}, // Add other models here as needed
	)
	if err != nil {
		log.Printf("❌ Auto migration failed: %v", err)
		return err
	}

	log.Println("✅ PostgreSQL auto migration complete")
	return nil
}
