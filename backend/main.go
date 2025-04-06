package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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
		api.GET("/videos", handleGetVideos(db))
		api.GET("/videos/:id", handleGetVideo(db))
		
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