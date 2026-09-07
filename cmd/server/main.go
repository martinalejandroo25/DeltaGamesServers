package main

import (
	"context"
	"fmt"
	"log"

	"github.com/deltagames/deltagamesservers/internal/api"
	"github.com/deltagames/deltagamesservers/internal/config"
	"github.com/deltagames/deltagamesservers/internal/database"
	"github.com/deltagames/deltagamesservers/internal/gameserver"
	"github.com/deltagames/deltagamesservers/internal/podman"
)

func main() {
	log.Println("==================================================")
	log.Println("  DeltaGamesServers Engine - Starting Up...")
	log.Println("==================================================")

	// 1. Load Configuration
	cfg := config.LoadConfig()
	log.Printf("[Config] Database Path: %s", cfg.DatabasePath)
	log.Printf("[Config] Podman Socket: %s", cfg.PodmanSocketPath)
	log.Printf("[Config] Templates Dir: %s", cfg.TemplatesDir)
	log.Printf("[Config] Data Dir     : %s", cfg.DataDir)

	// 2. Initialize Database
	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("[DB ERROR] Failed to initialize SQLite database: %v", err)
	}
	log.Println("[DB] SQLite database initialized successfully (WAL mode).")

	// 3. Connect to Podman REST API Socket
	podmanClient := podman.NewClient(cfg.PodmanSocketPath)
	if err := podmanClient.Ping(context.Background()); err != nil {
		log.Printf("[PODMAN WARNING] Podman socket ping failed: %v", err)
		log.Printf("[PODMAN WARNING] Ensure 'systemctl --user enable --now podman.socket' is active.")
	} else {
		log.Println("[PODMAN] Connection to Podman REST Socket verified.")
	}

	// 4. Initialize Game Server Lifecycle Manager
	manager := gameserver.NewManager(db, podmanClient, cfg)

	// 5. Setup Gin Router
	router := api.SetupRouter(cfg, db, manager)

	// 6. Start HTTP Server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("[HTTP] Listening and serving HTTP on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("[HTTP ERROR] Server failed to start: %v", err)
	}
}
