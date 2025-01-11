package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Login(db *sqlx.DB, c *gin.Context) {
	var user struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var storedPassword string
	query := "SELECT password FROM users WHERE name = ?"
	if err := db.Get(&storedPassword, query, user.Name); err != nil {
		c.JSON(500, gin.H{"error": "User not found"})
		return
	}

	if storedPassword != user.Password {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(200, gin.H{"message": "Login successful"})
}
