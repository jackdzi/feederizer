package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func DeleteAll(db *sqlx.DB, c *gin.Context) {
	fmt.Println("Called")
	query := "DELETE FROM data"
	_, err := db.Exec(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "All rows deleted successfully"})
}
