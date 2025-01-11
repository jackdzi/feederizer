package db

import (
	"github.com/jackdzi/feederizer/server/internal/api"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	whitelist := make(map[string]bool)
	whitelist["::1"] = true
	whitelist["172.20.0.1"] = true // Docker local machine (I think)

	router.Use(IPWhiteList(whitelist))
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running!"})
	})

	router.POST("/init", func(c *gin.Context) {
		api.DbInit(db, c)
	})

	router.POST("/credentials", func(c *gin.Context) {
		api.Login(db, c)
	})

	router.DELETE("/user", func(c *gin.Context) {
		api.DeleteAll(db, c)
	})

	router.POST("/user", func(c *gin.Context) {
		api.AddUser(db, c)
	})

	router.GET("/feed", func(c *gin.Context) {
		api.GetFeed(db, c)
	})

	return router
}
