package main

import (
	"net/http"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func main() {
	r := gin.Default()

	// Login: Issues the token
	r.POST("/login", func(c *gin.Context) {
		var login struct{ User, Pass string }
		c.BindJSON(&login)
		if login.User == "admin" && login.Pass == "password" {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": login.User})
			str, _ := token.SignedString(jwtKey)
			c.JSON(200, gin.H{"token": str})
			return
		}
		c.JSON(401, gin.H{"error": "unauthorized"})
	})

	// Verify: Other services call this to validate a user
	r.GET("/verify", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusOK)
	})

	r.Run(":5000")
}