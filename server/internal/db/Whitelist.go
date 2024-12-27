package db

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IPWhiteList(whitelist map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !whitelist[c.ClientIP()] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
        "ip": clientIP,
				"status":  http.StatusForbidden,
				"message": "Permission denied",
			})
			return
		}
	}
}
