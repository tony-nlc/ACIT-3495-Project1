package main

import (
	"bytes"
	"database/sql"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", os.Getenv("DB_DSN"))
	r := gin.Default()

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
		file, header, _ := c.FormFile("video")
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("video", header.Filename)
		io.Copy(part, file)
		writer.Close()
		
		req, _ := http.NewRequest("POST", "http://file-service:5001/save", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, _ := http.DefaultClient.Do(req)

		if resp.StatusCode == 200 {
			db.Exec("INSERT INTO videos (name, path) VALUES (?, ?)", header.Filename, "/storage/"+header.Filename)
			c.JSON(200, gin.H{"status": "Upload Complete"})
		}
	})
	r.Run(":5003")
}