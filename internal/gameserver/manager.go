package gameserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/deltagames/deltagamesservers/internal/config"
	"github.com/deltagames/deltagamesservers/internal/database"
	"github.com/deltagames/deltagamesservers/internal/database/models"
	"github.com/deltagames/deltagamesservers/internal/podman"
)

type Manager struct {
	db           *database.DB
	podmanClient *podman.Client
	cfg          *config.Config
}

func NewManager(db *database.DB, podmanClient *podman.Client, cfg *config.Config) *Manager {
	return &Manager{
		db:           db,
		podmanClient: podmanClient,
		cfg:          cfg,
	}
}

func (m *Manager) CreateServer(ctx context.Context, name, gameID string, customPorts map[string]int, customEnv map[string]string) (*models.GameServer, error) {
	// 1. Load template
	templatePath := filepath.Join(m.cfg.TemplatesDir, gameID+".json")
	tmplData, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("game template '%s' not found: %w", gameID, err)
	}

	var tmpl models.GameTemplate
	if err := json.Unmarshal(tmplData, &tmpl); err != nil {
		return nil, fmt.Errorf("invalid template format: %w", err)
	}

	// 2. Prepare Unique Server ID & Storage Directory
	serverID := fmt.Sprintf("%s-%d", gameID, time.Now().Unix())
	volumePath := filepath.Join(m.cfg.DataDir, serverID)
	if err := os.MkdirAll(volumePath, 0777); err != nil {
		return nil, fmt.Errorf("failed to create server data dir: %w", err)
	}

	// 3. Prepare Environment
	envVars := make(map[string]string)
	for k, v := range tmpl.Environment {
		envVars[k] = v
	}
	envVars["STEAM_APP_ID"] = tmpl.SteamAppID
	envVars["START_BINARY"] = tmpl.StartBinary
	envVars["START_PARAMS"] = tmpl.DefaultParams
	for k, v := range customEnv {
		envVars[k] = v
	}

	// 4. Build Port Mappings
	var portMappings []podman.PortMapping
	finalPortMap := make(map[string]int)

	for _, p := range tmpl.Ports {
		proto := p.Protocol
		if proto == "" {
			proto = "udp"
		}
		hostPort := p.ContainerPort
		if cp, ok := customPorts[fmt.Sprintf("%d/%s", p.ContainerPort, proto)]; ok {
			hostPort = cp
		}

		portKey := fmt.Sprintf("%d/%s", p.ContainerPort, proto)
		finalPortMap[portKey] = hostPort

		portMappings = append(portMappings, podman.PortMapping{
			HostPort:      uint16(hostPort),
			ContainerPort: uint16(p.ContainerPort),
			Protocol:      proto,
		})
	}

	// 5. Build Mounts
	mounts := []podman.MountSpec{
		{
			Source:      volumePath,
			Destination: "/home/steam/game",
			Type:        "bind",
		},
	}

	// 6. Request Podman Container Creation
	containerName := fmt.Sprintf("deltagames-%s", serverID)
	createSpec := podman.ContainerCreateSpec{
		Name:         containerName,
		Image:        tmpl.BaseImage,
		Env:          envVars,
		PortMappings: portMappings,
		Mounts:       mounts,
	}

	containerID, err := m.podmanClient.CreateContainer(ctx, createSpec)
	if err != nil {
		// Log error but proceed or mark error in DB
		containerID = ""
	}

	// 7. Save to DB
	server := &models.GameServer{
		ID:          serverID,
		Name:        name,
		GameID:      gameID,
		ContainerID: containerID,
		Status:      "stopped",
		PortMap:     finalPortMap,
		EnvVars:     envVars,
		VolumePath:  volumePath,
		AutoStart:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := m.db.SaveServer(server); err != nil {
		return nil, fmt.Errorf("failed to save server record: %w", err)
	}

	return server, nil
}

func (m *Manager) StartServer(ctx context.Context, id string) error {
	server, err := m.db.GetServer(id)
	if err != nil {
		return fmt.Errorf("server not found: %w", err)
	}

	if server.ContainerID == "" {
		return fmt.Errorf("server has no container ID assigned")
	}

	if err := m.podmanClient.StartContainer(ctx, server.ContainerID); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	server.Status = "running"
	return m.db.SaveServer(server)
}

func (m *Manager) StopServer(ctx context.Context, id string) error {
	server, err := m.db.GetServer(id)
	if err != nil {
		return fmt.Errorf("server not found: %w", err)
	}

	if server.ContainerID == "" {
		return fmt.Errorf("server has no container ID assigned")
	}

	if err := m.podmanClient.StopContainer(ctx, server.ContainerID, 15); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	server.Status = "stopped"
	return m.db.SaveServer(server)
}

func (m *Manager) DeleteServer(ctx context.Context, id string) error {
	server, err := m.db.GetServer(id)
	if err != nil {
		return err
	}

	if server.ContainerID != "" {
		_ = m.podmanClient.StopContainer(ctx, server.ContainerID, 5)
		_ = m.podmanClient.RemoveContainer(ctx, server.ContainerID, true)
	}

	return m.db.DeleteServer(id)
}

func (m *Manager) GetServerLogs(ctx context.Context, id string, tail int) (string, error) {
	server, err := m.db.GetServer(id)
	if err != nil {
		return "", err
	}
	if server.ContainerID == "" {
		return "Contenedor no inicializado.", nil
	}
	return m.podmanClient.GetLogs(ctx, server.ContainerID, true, true, tail)
}
