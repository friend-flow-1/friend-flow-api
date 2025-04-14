package main

import (
	"log"
	"net/http"

	"github.com/haxxu/friend-flow-api/internal"
	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/haxxu/friend-flow-api/internal/db"
	"github.com/haxxu/friend-flow-api/internal/routes"
)

func main() {
	cfg := config.LoadConfig()

	userHandler, err := internal.InitializeModules(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize modules: %v", err)
	}

	log.Fatal(http.ListenAndServe(":8080", routes.SetupRoutes(userHandler)))

	defer db.Session.Close()

	select {}
}
