package db

import (
	"fmt"
	"log"

	"github.com/haxxu/friend-flow-api/internal/config"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PG *gorm.DB

func InitPostgres(cfg *config.Config) (*gorm.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgreHost,
		cfg.PostgrePort,
		cfg.PostgreUser,
		cfg.PostgrePassword,
		cfg.PostgreDBName,
	)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL with GORM: %v", err)
		return nil, err
	}

	PG = db
	log.Println("✅ Connected to PostgreSQL with GORM")

	// 🔁 Auto-migrate PostgreSQL schema here
	AutoMigratePostgres(db)

	PG = db
	log.Println("✅ Connected to PostgreSQL")
	return db, nil
}
