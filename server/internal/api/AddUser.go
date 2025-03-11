package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AddUser(db *sqlx.DB, c *gin.Context) {
	var user struct {
		Name  string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO users (name, password) VALUES (?, ?)"
	_, err := db.Exec(query, user.Name, user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User added successfully"})
}
