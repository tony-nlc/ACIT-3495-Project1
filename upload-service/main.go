package main

import (
	"bytes"
	"database/sql"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", os.Getenv("DB_DSN"))
	r := gin.Default()

	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"}, 
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	r.POST("/upload", func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		// --- PASS THROUGH AUTH SERVICE ---
		authReq, _ := http.NewRequest("GET", "http://auth-service:5000/verify", nil)
		authReq.Header.Set("Authorization", token)
		authResp, err := http.DefaultClient.Do(authReq)
		if err != nil || authResp.StatusCode != 200 {
			c.JSON(401, gin.H{"error": "Auth Service rejected user"})
			return
		}

		// Proceed to File Service and MySQL...
		header, err := c.FormFile("video")
        if err != nil {
            c.JSON(400, gin.H{"error": "No video file found in request"})
            return
        }

        // 2. You MUST open the header to get the actual file data
        file, err := header.Open()
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to open video file"})
            return
        }
        defer file.Close() // Good practice to close it when the function finishes

        // 3. Prepare to pass it to the File Service
        body := &bytes.Buffer{}
        writer := multipart.NewWriter(body)
        
        // Use header.Filename here (it's now a valid string field)
        part, err := writer.CreateFormFile("video", header.Filename)
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to create multipart form"})
            return
        }

        // 4. Copy the actual file content (the io.Reader) into the part
        io.Copy(part, file)
        writer.Close()
		
		req, _ := http.NewRequest("POST", "http://file-service:5001/save", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", token)
		resp, _ := http.DefaultClient.Do(req)

		if resp.StatusCode == 200 {
			db.Exec("INSERT INTO videos (name, path) VALUES (?, ?)", header.Filename, "/storage/"+header.Filename)
			c.JSON(200, gin.H{"status": "Upload Complete"})
		}
	})
	r.Run(":5003")
}