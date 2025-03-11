package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func DeleteAll(db *sqlx.DB, c *gin.Context) {
	tx, err := db.Beginx()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to begin transaction: " + err.Error()})
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Transaction panicked"})
		}
	}()

	statements := []string{
		"DROP TABLE IF EXISTS feed_items;",
		"DROP TABLE IF EXISTS subscriptions;",
		"DROP TABLE IF EXISTS feeds;",
		"DROP TABLE IF EXISTS users;",
	}

	for _, stmt := range statements {
		_, err := tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to execute statement: " + err.Error()})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "All tables deleted successfully"})
}
