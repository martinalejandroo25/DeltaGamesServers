package podman

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	socketPath string
}

type ContainerCreateSpec struct {
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	Env          map[string]string `json:"env,omitempty"`
	PortMappings []PortMapping     `json:"port_mappings,omitempty"`
	Mounts       []MountSpec       `json:"mounts,omitempty"`
	Entrypoint   []string          `json:"entrypoint,omitempty"`
	Cmd          []string          `json:"command,omitempty"`
}

type PortMapping struct {
	HostPort      uint16 `json:"host_port"`
	ContainerPort uint16 `json:"container_port"`
	Protocol      string `json:"protocol"`
}

type MountSpec struct {
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Type        string `json:"Type"`
}

type ContainerCreateResponse struct {
	ID       string   `json:"Id"`
	Warnings []string `json:"Warnings"`
}

type ContainerInspectResponse struct {
	ID    string `json:"Id"`
	State struct {
		Status   string `json:"Status"` // running, exited, stopped, created
		Running  bool   `json:"Running"`
		ExitCode int    `json:"ExitCode"`
	} `json:"State"`
}

func NewClient(socketPath string) *Client {
	// If socketPath doesn't exist, we fallback or try to detect
	if socketPath == "" {
		socketPath = fmt.Sprintf("/run/user/%d/podman/podman.sock", os.Getuid())
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", socketPath, 10*time.Second)
		},
	}

	return &Client{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		socketPath: socketPath,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://d/v4.0.0/libpod/_ping", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("podman ping failed (socket: %s): %w", c.socketPath, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from podman ping: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) CreateContainer(ctx context.Context, spec ContainerCreateSpec) (string, error) {
	// Prepare libpod payload
	payload := map[string]interface{}{
		"name":  spec.Name,
		"image": spec.Image,
		"env":   spec.Env,
	}

	if len(spec.PortMappings) > 0 {
		var ports []map[string]interface{}
		for _, p := range spec.PortMappings {
			ports = append(ports, map[string]interface{}{
				"host_port":      p.HostPort,
				"container_port": p.ContainerPort,
				"protocol":       p.Protocol,
			})
		}
		payload["port_mappings"] = ports
	}

	if len(spec.Mounts) > 0 {
		var mounts []map[string]interface{}
		for _, m := range spec.Mounts {
			mounts = append(mounts, map[string]interface{}{
				"source":      m.Source,
				"destination": m.Destination,
				"type":        m.Type,
			})
		}
		payload["mounts"] = mounts
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://d/v4.0.0/libpod/containers/create", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("podman create container request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("podman create container failed [%d]: %s", resp.StatusCode, string(body))
	}

	var res ContainerCreateResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	return res.ID, nil
}

func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	endpoint := fmt.Sprintf("http://d/v4.0.0/libpod/containers/%s/start", url.PathEscape(containerID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotModified {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to start container [%d]: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *Client) StopContainer(ctx context.Context, containerID string, timeoutSec int) error {
	endpoint := fmt.Sprintf("http://d/v4.0.0/libpod/containers/%s/stop?t=%d", url.PathEscape(containerID), timeoutSec)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotModified {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to stop container [%d]: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *Client) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	endpoint := fmt.Sprintf("http://d/v4.0.0/libpod/containers/%s?force=%t&v=true", url.PathEscape(containerID), force)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to remove container [%d]: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *Client) InspectContainer(ctx context.Context, containerID string) (*ContainerInspectResponse, error) {
	endpoint := fmt.Sprintf("http://d/v4.0.0/libpod/containers/%s/json", url.PathEscape(containerID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("container inspect failed [%d]", resp.StatusCode)
	}

	var inspect ContainerInspectResponse
	if err := json.NewDecoder(resp.Body).Decode(&inspect); err != nil {
		return nil, err
	}
	return &inspect, nil
}

func (c *Client) GetLogs(ctx context.Context, containerID string, stdout, stderr bool, tail int) (string, error) {
	endpoint := fmt.Sprintf("http://d/v4.0.0/libpod/containers/%s/logs?stdout=%t&stderr=%t&tail=%d", url.PathEscape(containerID), stdout, stderr, tail)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.ToValidUTF8(string(body), ""), nil
}
