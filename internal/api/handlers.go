package api

import (
	"mclogs-go/internal/cache"
	"mclogs-go/internal/config"
	"mclogs-go/internal/filter"
	"mclogs-go/internal/parser"
	"mclogs-go/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage storage.Storage
	cache   cache.Cache
	parser  *parser.Parser
	cfg     *config.Config
}

func NewHandler(s storage.Storage, c cache.Cache, p *parser.Parser, cfg *config.Config) *Handler {
	return &Handler{
		storage: s,
		cache:   c,
		parser:  p,
		cfg:     cfg,
	}
}

func (h *Handler) CreateLog(c *gin.Context) {
	var req struct {
		Content string `form:"content" json:"content" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	content := req.Content

	// Apply filters
	lenFilter := &filter.LengthFilter{MaxLength: h.cfg.Storage.MaxLength}
	content, _ = lenFilter.Filter(content)

	lineFilter := &filter.LinesFilter{MaxLines: h.cfg.Storage.MaxLines}
	content, _ = lineFilter.Filter(content)

	// Store log
	id, err := h.storage.Put(c.Request.Context(), content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"id":      id,
		"url":     "https://mclogs.lemwood.icu/" + id,
	})
}

func (h *Handler) GetLog(c *gin.Context) {
	id := c.Param("id")

	log, err := h.storage.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve log"})
		return
	}

	if log == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}

	// In a real app, you might want to cache the analysis result
	analysis := h.parser.Parse(log.Content)
	analysis.ID = log.ID

	c.JSON(http.StatusOK, analysis)
}

func (h *Handler) GetAIAnalysis(c *gin.Context) {
	id := c.Param("id")

	logData, err := h.storage.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve log"})
		return
	}

	if logData == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}

	// Currently, we use the rule engine for analysis. 
	// In the future, this can be integrated with actual AI models.
	analysis := h.parser.Parse(logData.Content)
	
	// Format to match the frontend's expected AI analysis response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"analysis": analysis,
		"summary": "AI analysis provided by rule engine.",
	})
}

func (h *Handler) GetRawLog(c *gin.Context) {
	id := c.Param("id")

	log, err := h.storage.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve log"})
		return
	}

	if log == nil {
		c.String(http.StatusNotFound, "Log not found")
		return
	}

	c.String(http.StatusOK, log.Content)
}
