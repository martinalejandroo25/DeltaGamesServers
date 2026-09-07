package models

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Role      string    `json:"role"` // admin, manager, viewer
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GameServer struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	GameID      string            `json:"gameId"`
	ContainerID string            `json:"containerId"`
	Status      string            `json:"status"` // stopped, starting, running, stopping, error
	PortMap     map[string]int    `json:"portMap"` // e.g., {"27015/udp": 27015}
	EnvVars     map[string]string `json:"envVars"`
	VolumePath  string            `json:"volumePath"`
	AutoStart   bool              `json:"autoStart"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type GameTemplate struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	SteamAppID    string            `json:"steamAppId"`
	BaseImage     string            `json:"baseImage"`
	StartBinary   string            `json:"startBinary"`
	DefaultParams string            `json:"defaultParams"`
	Environment   map[string]string `json:"environment"`
	Ports         []PortSpec        `json:"ports"`
	Volumes       []VolumeSpec      `json:"volumes"`
}

type PortSpec struct {
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol"` // UDP or TCP
	Description   string `json:"description"`
}

type VolumeSpec struct {
	ContainerPath string `json:"containerPath"`
	Description   string `json:"description"`
}

type ServerMetric struct {
	ServerID  string    `json:"serverId"`
	CPUUsage  float64   `json:"cpuUsage"`  // percentage
	MemoryMB  float64   `json:"memoryMb"`  // in MB
	MemoryLimit float64 `json:"memoryLimitMb"`
	NetRxBytes uint64   `json:"netRxBytes"`
	NetTxBytes uint64   `json:"netTxBytes"`
	Timestamp time.Time `json:"timestamp"`
}
