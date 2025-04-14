package db

import (
	_ "fmt"
	"log"
	"time"

	"github.com/haxxu/friend-flow-api/internal/config"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitScylla(cfg *config.Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.DBHost)
	cluster.Keyspace = cfg.DBKeyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second

	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect ScyllaDB: %v", err)
	}
	log.Println("✅ Connected to ScyllaDB")

	return Session, nil
}
