package db

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
  "feederizer/server/internal/api"
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

  router.POST("/init", func(c *gin.Context) {
    api.DbInit(db, c)
  })

  router.POST("/delete", func(c *gin.Context) {
    api.DeleteAll(db, c)
  })

	router.POST("/adduser", func(c *gin.Context) {
		api.AddUser(db, c)
	})

	router.GET("/feed", func(c *gin.Context) {
		api.GetFeed(db, c)
	})

	return router
}
