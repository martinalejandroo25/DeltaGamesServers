package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deltagames/deltagamesservers/internal/database/models"
	_ "modernc.org/sqlite"
)

type DB struct {
	Conn *sql.DB
}

func InitDB(dbPath string) (*DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	conn, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{Conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS game_servers (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		game_id TEXT NOT NULL,
		container_id TEXT,
		status TEXT NOT NULL,
		port_map TEXT,
		env_vars TEXT,
		volume_path TEXT,
		auto_start BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Conn.Exec(schema)
	return err
}

func (db *DB) SaveServer(server *models.GameServer) error {
	portMapJSON, _ := json.Marshal(server.PortMap)
	envVarsJSON, _ := json.Marshal(server.EnvVars)

	query := `
	INSERT INTO game_servers (id, name, game_id, container_id, status, port_map, env_vars, volume_path, auto_start, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT(id) DO UPDATE SET
		name=excluded.name,
		container_id=excluded.container_id,
		status=excluded.status,
		port_map=excluded.port_map,
		env_vars=excluded.env_vars,
		volume_path=excluded.volume_path,
		auto_start=excluded.auto_start,
		updated_at=CURRENT_TIMESTAMP;
	`
	_, err := db.Conn.Exec(query, server.ID, server.Name, server.GameID, server.ContainerID, server.Status, string(portMapJSON), string(envVarsJSON), server.VolumePath, server.AutoStart)
	return err
}

func (db *DB) GetServer(id string) (*models.GameServer, error) {
	query := `SELECT id, name, game_id, container_id, status, port_map, env_vars, volume_path, auto_start, created_at, updated_at FROM game_servers WHERE id = ?`
	row := db.Conn.QueryRow(query, id)

	var s models.GameServer
	var portMapStr, envVarsStr string
	err := row.Scan(&s.ID, &s.Name, &s.GameID, &s.ContainerID, &s.Status, &portMapStr, &envVarsStr, &s.VolumePath, &s.AutoStart, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal([]byte(portMapStr), &s.PortMap)
	_ = json.Unmarshal([]byte(envVarsStr), &s.EnvVars)
	return &s, nil
}

func (db *DB) ListServers() ([]models.GameServer, error) {
	query := `SELECT id, name, game_id, container_id, status, port_map, env_vars, volume_path, auto_start, created_at, updated_at FROM game_servers`
	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []models.GameServer
	for rows.Next() {
		var s models.GameServer
		var portMapStr, envVarsStr string
		if err := rows.Scan(&s.ID, &s.Name, &s.GameID, &s.ContainerID, &s.Status, &portMapStr, &envVarsStr, &s.VolumePath, &s.AutoStart, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(portMapStr), &s.PortMap)
		_ = json.Unmarshal([]byte(envVarsStr), &s.EnvVars)
		servers = append(servers, s)
	}
	return servers, nil
}

func (db *DB) DeleteServer(id string) error {
	_, err := db.Conn.Exec("DELETE FROM game_servers WHERE id = ?", id)
	return err
}
