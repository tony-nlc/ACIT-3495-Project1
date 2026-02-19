package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
        AllowAllOrigins:  true, 
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization",  "content-type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	r.POST("/upload", func(c *gin.Context) {
    token := c.GetHeader("Authorization") // Assuming this is "Bearer <token>"

    // 1. Auth check
    authReq, _ := http.NewRequest("GET", "http://auth-service:5000/verify", nil)
    authReq.Header.Set("Authorization", token)
    authResp, err := http.DefaultClient.Do(authReq)
    if err != nil || authResp.StatusCode != 200 {
        c.JSON(401, gin.H{"error": "Auth Service rejected user"})
        return
    }

    // 2. Prepare File
    header, _ := c.FormFile("video")
    file, _ := header.Open()
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, _ := writer.CreateFormFile("video", header.Filename)
    io.Copy(part, file)
    writer.Close()

    // 3. Call File Service
    req, _ := http.NewRequest("POST", "http://file-service:5001/save", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Set("Authorization", token) // Ensure this has "Bearer " prefix

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        c.JSON(500, gin.H{"error": "File service unreachable"})
        return
    }
    defer resp.Body.Close()

    // 4. Check response from File Service
    if resp.StatusCode != 200 {
        c.JSON(resp.StatusCode, gin.H{"error": "File service failed to save"})
        return
    }

    c.JSON(200, gin.H{"status": "Upload Complete", "message": "Saved via file-service"})
})
	r.Run(":5003")
}