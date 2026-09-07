package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/deltagames/deltagamesservers/internal/config"
	"github.com/deltagames/deltagamesservers/internal/database"
	"github.com/deltagames/deltagamesservers/internal/gameserver"
	"github.com/deltagames/deltagamesservers/internal/podman"
)

func TestHealthEndpoint(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "deltagames_api_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Conn.Close()

	cfg := &config.Config{
		ServerPort:   8080,
		DatabasePath: dbPath,
		TemplatesDir: "../../containers/templates",
		DataDir:      filepath.Join(tempDir, "data"),
	}

	podmanClient := podman.NewClient("")
	manager := gameserver.NewManager(db, podmanClient, cfg)
	router := SetupRouter(cfg, db, manager)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
