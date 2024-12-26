package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
  "fmt"
)

func NewRouter(db *sqlx.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running!"})
	})

	router.POST("/data", func(c *gin.Context) {
		var jsonData map[string]interface{}
		if err := c.BindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		fmt.Printf("JSON Data: %+v\n", jsonData)

		query, args, err := sqlx.Named("INSERT INTO feeds (column1, column2) VALUES (:key1, :key2)", jsonData)
    fmt.Println(query)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare query"})
			return
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"received": jsonData})
	})

	return router
}
