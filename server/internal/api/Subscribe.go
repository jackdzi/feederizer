package api

import(
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// TODO: Add route

func SubscribeToFeed(db *sqlx.DB, c *gin.Context) {
	var subscription struct {
		UserID int `json:"user_id"`
		FeedID int `json:"feed_id"`
	}

	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if subscription already exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM subscriptions WHERE user_id = ? AND feed_id = ?)`
	if err := db.Get(&exists, checkQuery, subscription.UserID, subscription.FeedID); err != nil {
		c.JSON(500, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	if exists {
		c.JSON(400, gin.H{"error": "Subscription already exists"})
		return
	}
