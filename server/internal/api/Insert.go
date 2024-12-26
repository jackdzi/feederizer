package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Insert(db *sqlx.DB, c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	query, args, err := sqlx.Named("INSERT INTO data (column1, column2) VALUES (:key1, :key2)", jsonData)
	query = db.Rebind(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare query"})
		return
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
		return
	}
  fmt.Println("c")

	c.JSON(http.StatusOK, gin.H{"received": jsonData})
}
