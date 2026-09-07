package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	ServerPort       int
	DatabasePath     string
	PodmanSocketPath string
	TemplatesDir     string
	DataDir          string
	JWTSecret        string
}

func LoadConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}

	uid := os.Getuid()
	defaultSocket := filepath.Join("/run", "user", strconv.Itoa(uid), "podman", "podman.sock")

	return &Config{
		ServerPort:       getEnvAsInt("PORT", 8080),
		DatabasePath:     getEnv("DB_PATH", filepath.Join(homeDir, ".deltagames", "deltagames.db")),
		PodmanSocketPath: getEnv("PODMAN_SOCKET", defaultSocket),
		TemplatesDir:     getEnv("TEMPLATES_DIR", "containers/templates"),
		DataDir:          getEnv("DATA_DIR", filepath.Join(homeDir, ".deltagames", "servers")),
		JWTSecret:        getEnv("JWT_SECRET", "deltagames-secret-key-change-in-production"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}
