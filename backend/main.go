package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Port           int
	DBPath         string
	VideoStorePath string
	LNDHost        string
	TLSCertPath    string
	MacaroonPath   string
}

var config Config

func init() {
	flag.IntVar(&config.Port, "port", 8080, "Port to run the server on")
	flag.StringVar(&config.DBPath, "db", "videos.db", "Path to SQLite database")
	flag.StringVar(&config.VideoStorePath, "videos", "videos", "Path to store video files")
	flag.StringVar(&config.LNDHost, "lnd-host", "localhost:10009", "LND gRPC host")
	flag.StringVar(&config.TLSCertPath, "tls-cert", "", "Path to LND TLS certificate")
	flag.StringVar(&config.MacaroonPath, "macaroon", "", "Path to LND macaroon file")
}

func main() {
	flag.Parse()

	// Ensure video storage directory exists
	if err := os.MkdirAll(config.VideoStorePath, 0755); err != nil {
		log.Fatal("Failed to create video storage directory:", err)
	}

	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Serve static files for the frontend
	r.Static("/assets", "./dist/assets")
	r.StaticFile("/", "./dist/index.html")

	// API routes
	api := r.Group("/api")
	{
		api.GET("/videos", handleGetVideos(db))
		api.GET("/videos/:id", handleGetVideo(db))
		
		// Admin routes
		admin := api.Group("/admin")
		admin.Use(authMiddleware())
		{
			admin.POST("/videos", handleUploadVideo(db))
			admin.PUT("/videos/:id", handleUpdateVideo(db))
			admin.DELETE("/videos/:id", handleDeleteVideo(db))
		}
	}

	// Start server
	addr := fmt.Sprintf(":%d", config.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS videos (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			filename TEXT NOT NULL,
			duration INTEGER,
			thumbnail TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)

	return db, err
}

// Handlers will be implemented here
func handleGetVideos(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}

func handleGetVideo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}

func handleUploadVideo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}

func handleUpdateVideo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}

func handleDeleteVideo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}