package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" 
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
	db        *sql.DB
)

func init() {
	// Best Practice: The DSN is constructed from environment variables
	dsn := os.Getenv("DB_DSN") 
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()
	storagePath := "/storage"

	r.Use(cors.New(cors.Config{
        AllowAllOrigins:  true, 
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization",  "content-type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))


	r.POST("/save", authMiddleware(), func(c *gin.Context) {
		file, err := c.FormFile("video")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		// 1. Save file to Disk
		filePath := fmt.Sprintf("%s/%s", storagePath, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
    log.Printf("DISK ERROR: %v", err) // <--- Add this
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
    return
}

		// 2. Save metadata to MySQL
		query := "INSERT INTO videos (title, file_path) VALUES (?, ?)"
		_, err = db.Exec(query, file.Filename, filePath)
		if err != nil {
    log.Printf("DATABASE ERROR: %v", err) // <--- Add this
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save metadata to DB"})
    return
}

		c.JSON(http.StatusOK, gin.H{"message": "File and Metadata saved successfully"})
	})

	r.Run(":5001")
}