package main

import (
	"fmt"
	"log"
	"mclogs-go/internal/api"
	"mclogs-go/internal/cache"
	"mclogs-go/internal/config"
	"mclogs-go/internal/parser"
	"mclogs-go/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Storage
	var store storage.Storage
	if cfg.Storage.Postgres.Enabled {
		store, err = storage.NewPostgresStorage(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize Postgres storage: %v", err)
		}
		log.Println("Initialized Postgres storage")
	} else if cfg.Storage.MongoDB.Enabled {
		store, err = storage.NewMongoStorage(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize MongoDB storage: %v", err)
		}
		log.Println("Initialized MongoDB storage")
	} else {
		log.Fatal("No storage driver enabled in config")
	}

	// Initialize Cache
	var c cache.Cache
	if cfg.Cache.Enabled {
		if cfg.Cache.Driver == "redis" {
			c = cache.NewRedisCache(cfg)
		}
	}

	// Initialize Parser
	p, err := parser.NewParser(cfg.Patterns)
	if err != nil {
		log.Fatalf("Failed to initialize Parser: %v", err)
	}

	// Set up Router
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize Handlers
	h := api.NewHandler(store, c, p, cfg)

	// Routes
	v1 := r.Group("/1")
	{
		v1.POST("/log", h.CreateLog)
		v1.GET("/log/:id", h.GetLog)
		v1.GET("/insights/:id", h.GetLog)
		v1.GET("/ai-analysis/:id", h.GetAIAnalysis) // Add this line to fix 404 for AI analysis
		v1.GET("/raw/:id", h.GetRawLog)
	}

	// Start Server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
