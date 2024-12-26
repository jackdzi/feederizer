package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
  whitelist := make(map[string]bool)
  whitelist["::1"] = true
  whitelist["172.20.0.1"] = true

  router.Use(IPWhiteList(whitelist))
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running!"})
	})

  router.POST("/delete", func(c *gin.Context) {
    DeleteAll(db, c)
  })

	router.POST("/insert", func(c *gin.Context) {
		Insert(db, c)
	})

	router.GET("/feed", func(c *gin.Context) {
		GetFeed(db, c)
	})

	return router
}
