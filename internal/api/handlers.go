package api

import (
	"log"
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
	log.Printf("[API] GetLog request for ID: %s", id)

	logData, err := h.storage.Get(c.Request.Context(), id)
	if err != nil {
		log.Printf("[API] Error retrieving log %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve log"})
		return
	}

	if logData == nil {
		log.Printf("[API] Log not found: %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}

	// In a real app, you might want to cache the analysis result
	analysis := h.parser.Parse(logData.Content)
	analysis.ID = logData.ID

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
	
	// Format the analysis result into a markdown string for the frontend's AI display
	var markdown string
	markdown += "## 日志分析摘要\n\n"
	if len(analysis.Information) > 0 {
		markdown += "### 基本信息\n"
		for _, info := range analysis.Information {
			markdown += "- **" + info.Label + "**: " + info.Value + "\n"
		}
		markdown += "\n"
	}

	if len(analysis.Problems) > 0 {
		markdown += "### 发现的问题\n"
		for _, prob := range analysis.Problems {
			severityEmoji := "⚠️"
			if prob.Severity == "error" || prob.Severity == "critical" {
				severityEmoji = "❌"
			}
			markdown += "#### " + severityEmoji + " " + prob.Message + "\n"
			if len(prob.Solutions) > 0 {
				markdown += "**建议方案:**\n"
				for _, sol := range prob.Solutions {
					markdown += "- " + sol.Message + "\n"
				}
			}
			markdown += "\n"
		}
	} else {
		markdown += "✅ 未发现明显的问题。\n"
	}

	// Format to match the frontend's expected AI analysis response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"analysis": markdown,
	})
}

func (h *Handler) GetRawLog(c *gin.Context) {
	id := c.Param("id")
	log.Printf("[API] GetRawLog request for ID: %s", id)

	logData, err := h.storage.Get(c.Request.Context(), id)
	if err != nil {
		log.Printf("[API] Error retrieving raw log %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve log"})
		return
	}

	if logData == nil {
		log.Printf("[API] Raw log not found: %s", id)
		c.String(http.StatusNotFound, "Log not found")
		return
	}

	c.String(http.StatusOK, logData.Content)
}
