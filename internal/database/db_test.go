package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/deltagames/deltagamesservers/internal/database/models"
)

func TestDB_Operations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "deltagames_db_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Conn.Close()

	server := &models.GameServer{
		ID:          "gmod-test-123",
		Name:        "Test Garry's Mod",
		GameID:      "gmod",
		ContainerID: "mock-container-456",
		Status:      "stopped",
		PortMap:     map[string]int{"27015/udp": 27015},
		EnvVars:     map[string]string{"STEAM_APP_ID": "4020"},
		VolumePath:  "/tmp/test-volume",
		AutoStart:   false,
	}

	if err := db.SaveServer(server); err != nil {
		t.Fatalf("SaveServer failed: %v", err)
	}

	fetched, err := db.GetServer("gmod-test-123")
	if err != nil {
		t.Fatalf("GetServer failed: %v", err)
	}
	if fetched.Name != "Test Garry's Mod" {
		t.Errorf("Expected name 'Test Garry's Mod', got '%s'", fetched.Name)
	}

	servers, err := db.ListServers()
	if err != nil {
		t.Fatalf("ListServers failed: %v", err)
	}
	if len(servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(servers))
	}

	if err := db.DeleteServer("gmod-test-123"); err != nil {
		t.Fatalf("DeleteServer failed: %v", err)
	}

	serversAfter, _ := db.ListServers()
	if len(serversAfter) != 0 {
		t.Errorf("Expected 0 servers after deletion, got %d", len(serversAfter))
	}
}
