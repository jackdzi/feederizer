package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// TODO: Add route


func UnsubscribeFromFeed(db *sqlx.DB, c *gin.Context) {
	var subscription struct {
		UserID int `json:"user_id"`
		FeedID int `json:"feed_id"`
	}

	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("DELETE FROM subscriptions WHERE user_id = ? AND feed_id = ?",
		subscription.UserID, subscription.FeedID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to unsubscribe: " + err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Successfully unsubscribed from feed"})
}
