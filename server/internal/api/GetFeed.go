package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Feed struct {
	Id string `db:"id"`
	Name string `db:"name"`
  Password string `db:"password"`
}

func GetFeed(db *sqlx.DB, c *gin.Context) {
	var feeds []Feed
	err := db.Select(&feeds, "SELECT id, name, password FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, feeds)
}
