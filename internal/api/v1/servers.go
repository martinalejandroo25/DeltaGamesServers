package v1

import (
	"net/http"
	"strconv"

	"github.com/deltagames/deltagamesservers/internal/database"
	"github.com/deltagames/deltagamesservers/internal/gameserver"
	"github.com/gin-gonic/gin"
)

type ServerHandler struct {
	db      *database.DB
	manager *gameserver.Manager
}

func NewServerHandler(db *database.DB, manager *gameserver.Manager) *ServerHandler {
	return &ServerHandler{
		db:      db,
		manager: manager,
	}
}

type CreateServerRequest struct {
	Name        string            `json:"name" binding:"required"`
	GameID      string            `json:"gameId" binding:"required"`
	CustomPorts map[string]int    `json:"customPorts"`
	CustomEnv   map[string]string `json:"customEnv"`
}

func (h *ServerHandler) ListServers(c *gin.Context) {
	servers, err := h.db.ListServers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, servers)
}

func (h *ServerHandler) GetServer(c *gin.Context) {
	id := c.Param("id")
	server, err := h.db.GetServer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}
	c.JSON(http.StatusOK, server)
}

func (h *ServerHandler) CreateServer(c *gin.Context) {
	var req CreateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	server, err := h.manager.CreateServer(c.Request.Context(), req.Name, req.GameID, req.CustomPorts, req.CustomEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, server)
}

func (h *ServerHandler) StartServer(c *gin.Context) {
	id := c.Param("id")
	if err := h.manager.StartServer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Server started successfully"})
}

func (h *ServerHandler) StopServer(c *gin.Context) {
	id := c.Param("id")
	if err := h.manager.StopServer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Server stopped successfully"})
}

func (h *ServerHandler) DeleteServer(c *gin.Context) {
	id := c.Param("id")
	if err := h.manager.DeleteServer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Server deleted successfully"})
}

func (h *ServerHandler) GetLogs(c *gin.Context) {
	id := c.Param("id")
	tailStr := c.DefaultQuery("tail", "100")
	tail, _ := strconv.Atoi(tailStr)

	logs, err := h.manager.GetServerLogs(c.Request.Context(), id, tail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
