package main

import (
	"database/sql"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", os.Getenv("DB_DSN"))
	r := gin.Default()

	r.GET("/getvideos", func(c *gin.Context) {
		rows, _ := db.Query("SELECT id, title FROM videos")
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

		// --- PASS THROUGH AUTH SERVICE ---
		authReq, _ := http.NewRequest("GET", "http://auth-service:5000/verify", nil)
		authReq.Header.Set("Authorization", token)
		authResp, err := http.DefaultClient.Do(authReq)
		if err != nil || authResp.StatusCode != 200 {
			c.JSON(401, gin.H{"error": "Unauthorized Access"})
			return
		}

		var path string
		db.QueryRow("SELECT path FROM videos WHERE id = ?", c.Param("id")).Scan(&path)
		c.File(path)
	})
	r.Run(":5002")
}