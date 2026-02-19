package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 1. Connection Check: Don't ignore the error here
	db, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
        AllowAllOrigins:  true, 
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization",  "content-type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	r.GET("/getvideos", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title FROM videos")
		if err != nil {
			c.JSON(500, gin.H{"error": "Database query failed"})
			return
		}
		defer rows.Close()

		var videos []gin.H
		for rows.Next() {
			var id int
			var title string
			rows.Scan(&id, &title)
			videos = append(videos, gin.H{"id": id, "title": title})
		}
		c.JSON(200, videos)
	})

	r.GET("/view/:id", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "No Authorization header provided"})
			return
		}

		// --- PASS THROUGH AUTH SERVICE ---
		// Image of microservices authentication flow
		
		authReq, _ := http.NewRequest("GET", "http://auth-service:5000/verify", nil)
		authReq.Header.Set("Authorization", token)

		authResp, err := http.DefaultClient.Do(authReq)
		if err != nil {
			log.Printf("Auth Service unreachable: %v", err)
			c.JSON(500, gin.H{"error": "Auth Service unreachable"})
			return
		}
		defer authResp.Body.Close()

		if authResp.StatusCode != 200 {
			log.Printf("Auth Service rejected token for ID %s", c.Param("id"))
			c.JSON(401, gin.H{"error": "Unauthorized Access"})
			return
		}

		// --- GET FILE PATH ---
		var path string
		err = db.QueryRow("SELECT file_path FROM videos WHERE id = ?", c.Param("id")).Scan(&path)
		if err != nil {
			log.Printf("Video ID %s not found: %v", c.Param("id"), err)
			c.JSON(404, gin.H{"error": "Video not found"})
			return
		}

		// --- FILE SYSTEM CHECK ---
		// Verify file exists on the shared volume before serving
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("File missing on disk at: %s", path)
			c.JSON(404, gin.H{"error": "Video file missing on storage"})
			return
		}

		// --- SERVE FILE ---
		c.Header("Content-Type", "video/mp4")
		c.Header("Content-Disposition", "inline")
		c.File(path)
	})

	r.Run(":5002")
}