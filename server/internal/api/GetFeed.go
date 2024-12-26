package api

import (
	"net/http"
  "fmt"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Feed struct {
	Column1 string `db:"column1"`
	Column2 string `db:"column2"`
}

func GetFeed(db *sqlx.DB, c *gin.Context) {
  fmt.Println("Called")
	var feeds []Feed
	err := db.Select(&feeds, "SELECT column1, column2 FROM data")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, feeds)
}
