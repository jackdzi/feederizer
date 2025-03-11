package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


//TODO: FIX
func GetSubscribedFeedItems(db *sqlx.DB, c *gin.Context, user string) {
	query := `
select
    f.feed_name,
    fi.feed_title,
    fi.link,
    fi.pub_date,
    fi.content
from users u
join subscriptions s on u.id = s.user_id
join feeds f on s.feed_id = f.id
join feed_items fi on f.id = fi.parent_feed_id
where u.name = ?
order by fi.pub_date desc;
	`
	fmt.Println(user)
	var feedItems []FeedItem
	err := db.Select(&feedItems, query, user)
	if err != nil {
		fmt.Println("Error fetching feed items:", err) // Add logging for debugging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feed items"})
		return
	}
  fmt.Println(feedItems)

	if len(feedItems) == 0 {
		c.JSON(http.StatusOK, []FeedItem{}) // Return empty array instead of null
		return
	}
	c.JSON(http.StatusOK, feedItems)
}
