package v1

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/deltagames/deltagamesservers/internal/config"
	"github.com/deltagames/deltagamesservers/internal/database/models"
	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	cfg *config.Config
}

func NewTemplateHandler(cfg *config.Config) *TemplateHandler {
	return &TemplateHandler{cfg: cfg}
}

func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	files, err := os.ReadDir(h.cfg.TemplatesDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read templates directory"})
		return
	}

	var templates []models.GameTemplate
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			content, err := os.ReadFile(filepath.Join(h.cfg.TemplatesDir, file.Name()))
			if err != nil {
				continue
			}
			var tmpl models.GameTemplate
			if err := json.Unmarshal(content, &tmpl); err == nil {
				templates = append(templates, tmpl)
			}
		}
	}

	c.JSON(http.StatusOK, templates)
}
