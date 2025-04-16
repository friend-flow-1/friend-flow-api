package main

import (
	"log"
	"net/http"

	"github.com/haxxu/friend-flow-api/internal"
	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/haxxu/friend-flow-api/internal/db"
)

func main() {
	cfg := config.LoadConfig()

	app, err := internal.InitializeApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	log.Fatal(http.ListenAndServe(":"+cfg.ApiPort, app.Router))

	defer db.Session.Close()

	select {}
}
