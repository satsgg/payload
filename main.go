package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"payload/backend/video"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("your-secret-key") // In production, use a secure key from environment variables

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
	// In production, serve the built frontend files
	if _, err := os.Stat("./dist"); err == nil {
		r.Static("/assets", "./dist/assets")
		r.StaticFile("/", "./dist/index.html")
		r.NoRoute(func(c *gin.Context) {
			c.File("./dist/index.html")
		})
	}

	// API routes
	api := r.Group("/api")
	{
		// Public routes
		api.GET("/videos", handleListVideos(db))
		api.GET("/videos/:id", handleGetVideo(db))
		api.GET("/videos/:id/*filepath", handleStreamVideo(db))
		
		// Admin routes
		admin := api.Group("/admin")
		{
			admin.POST("/login", handleLogin(db))
			admin.GET("/verify", authMiddleware(), handleVerifyToken())
			admin.Use(authMiddleware())
			{
				admin.POST("/videos", handleUploadVideo(db))
				admin.PUT("/videos/:id", handleUpdateVideo(db))
				admin.DELETE("/videos/:id", handleDeleteVideo(db))
			}
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
		);

		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			is_admin BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		-- Insert default admin user if not exists
		INSERT OR IGNORE INTO users (username, password_hash, is_admin)
		VALUES ('admin', '$2a$10$jM3CsbO3jaxAx760kmOBbOR8YgPpTU093EGPg8/W4tc3iM7axOM/6', TRUE);
	`)

	return db, err
}

// Handlers will be implemented here
func handleListVideos(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Handling list videos request")
		
		rows, err := db.Query("SELECT id, title, description, duration, thumbnail, created_at FROM videos ORDER BY created_at DESC")
		if err != nil {
			log.Printf("Database query error: %v", err)
			c.JSON(500, gin.H{"error": "Failed to fetch videos", "details": err.Error()})
			return
		}
		defer rows.Close()

		// Initialize an empty array
		videos := make([]gin.H, 0)

		for rows.Next() {
			var id, title, description, thumbnail string
			var duration int
			var createdAt time.Time
			err := rows.Scan(&id, &title, &description, &duration, &thumbnail, &createdAt)
			if err != nil {
				log.Printf("Row scanning error: %v", err)
				c.JSON(500, gin.H{"error": "Failed to scan video", "details": err.Error()})
				return
			}
			videos = append(videos, gin.H{
				"id":            id,
				"title":         title,
				"description":   description,
				"duration":      duration,
				"thumbnail_path": thumbnail,
				"created_at":    createdAt.Format(time.RFC3339),
			})
		}

		if err = rows.Err(); err != nil {
			log.Printf("Error after row iteration: %v", err)
			c.JSON(500, gin.H{"error": "Error processing video list", "details": err.Error()})
			return
		}

		log.Printf("Successfully fetched %d videos", len(videos))
		c.JSON(200, videos)
	}
}

func handleGetVideo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		videoID := c.Param("id")
		if videoID == "" {
			c.JSON(400, gin.H{"error": "Video ID is required"})
			return
		}

		var id, title, description string
		var duration int
		var createdAt time.Time
		err := db.QueryRow(
			"SELECT id, title, description, duration, created_at FROM videos WHERE id = ?",
			videoID,
		).Scan(&id, &title, &description, &duration, &createdAt)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Video not found"})
			} else {
				c.JSON(500, gin.H{"error": "Failed to fetch video"})
			}
			return
		}

		c.JSON(200, gin.H{
			"id":          id,
			"title":       title,
			"description": description,
			"duration":    duration,
			"created_at":  createdAt.Format(time.RFC3339),
		})
	}
}

func handleUploadVideo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get form data
		title := c.PostForm("title")
		description := c.PostForm("description")
		file, err := c.FormFile("video")
		if err != nil {
			c.JSON(400, gin.H{"error": "No video file provided"})
			return
		}

		// Generate unique ID for the video
		videoID := generateVideoID()

		// Create video directory
		videoDir := filepath.Join(config.VideoStorePath, videoID)
		if err := os.MkdirAll(videoDir, 0755); err != nil {
			c.JSON(500, gin.H{"error": "Failed to create video directory"})
			return
		}

		// Save uploaded file
		uploadPath := filepath.Join(videoDir, "original.mp4")
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save video file"})
			return
		}

		// Initialize video processor
		processor, err := video.NewVideoProcessor()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to initialize video processor"})
			return
		}

		// Get video duration
		duration, err := processor.GetVideoDuration(uploadPath)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get video duration"})
			return
		}

		// Generate thumbnail
		_, err = processor.GenerateThumbnail(uploadPath, videoDir)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate thumbnail"})
			return
		}

		// Convert to HLS
		_, err = processor.ConvertToHLS(uploadPath, videoDir)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to convert video to HLS"})
			return
		}

		// Save video info to database
		_, err = db.Exec(`
			INSERT INTO videos (
				id, title, description, filename, duration, thumbnail, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?)
		`, videoID, title, description, "playlist.m3u8", duration, "thumbnail.jpg", time.Now())

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to save video info"})
			return
		}

		c.JSON(200, gin.H{
			"id": videoID,
			"title": title,
			"description": description,
			"duration": duration,
		})
	}
}

func generateVideoID() string {
	// Generate a random 8-character string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
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
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func handleLogin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		var passwordHash string
		var isAdmin bool
		err := db.QueryRow("SELECT password_hash, is_admin FROM users WHERE username = ?", req.Username).Scan(&passwordHash, &isAdmin)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// Create JWT token
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Username: req.Username,
			IsAdmin:  isAdmin,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(200, gin.H{"token": tokenString})
	}
}

func handleVerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// The authMiddleware already verified the token
		c.JSON(200, gin.H{"status": "valid"})
	}
}

func handleStreamVideo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		videoID := c.Param("id")
		requestedFile := c.Param("filepath")
		
		if videoID == "" {
			c.JSON(400, gin.H{"error": "Video ID is required"})
			return
		}

		// Construct the full path to the video file
		videoPath := filepath.Join(config.VideoStorePath, videoID, requestedFile)
		
		// Check if file exists
		if _, err := os.Stat(videoPath); os.IsNotExist(err) {
			c.JSON(404, gin.H{"error": "Video file not found"})
			return
		}

		// Set appropriate headers based on file type
		if strings.HasSuffix(requestedFile, ".m3u8") {
			c.Header("Content-Type", "application/x-mpegURL")
		} else if strings.HasSuffix(requestedFile, ".ts") {
			c.Header("Content-Type", "video/mp2t")
		}
		
		// Stream the file
		c.File(videoPath)
	}
}