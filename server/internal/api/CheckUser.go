package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CheckUser(db *sqlx.DB, c *gin.Context) string {
	var user struct {
		Name  string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return ""
	}

	query := "SELECT id FROM users WHERE name = ? AND password = ?"
	rows, err := db.Queryx(query, user.Name, user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return ""
	}
	defer rows.Close()

	if rows.Next() {
		c.JSON(200, gin.H{"message": "User found"})
    return user.Name
	} else {
		c.JSON(404, gin.H{"error": "User not found"})
    return ""
	}
}
