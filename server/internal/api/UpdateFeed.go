package api

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// TODO: add routes

func UpdateFeed(db *sqlx.DB, c *gin.Context) {
	userID := c.Param("userID")

	var feeds []Feed
	query := `SELECT f.id, f.feed_name, f.api_url
			FROM feeds f
			JOIN subscriptions s ON f.id = s.feed_id
			WHERE s.user_id = ?`

	if err := db.Select(&feeds, query, userID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve feeds: " + err.Error()})
		return
	}

	totalNewItems := 0
	for _, feed := range feeds {
		newItems, err := fetchAndStoreItems(db, feed)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to refresh feed: " + err.Error()})
			return
		}
		totalNewItems += newItems
	}

	c.JSON(200, gin.H{"message": "Feeds refreshed successfully", "new_items": totalNewItems})
}

func fetchAndStoreItems(db *sqlx.DB, feed Feed) (int, error) {
	resp, err := http.Get(feed.ApiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return 0, err
	}

	newItems := 0
	tx, err := db.Beginx()
	if err != nil {
		return 0, err
	}

	for _, item := range rss.Channel.Items {
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			pubDate, err = time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				pubDate = time.Now()
			}
		}

		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM feed_items WHERE parent_feed_id = ? AND link = ?)`
		if err := tx.Get(&exists, checkQuery, feed.ID, item.Link); err != nil {
			tx.Rollback()
			return 0, err
		}

		if !exists {
			content := item.Content
			if content == "" {
				content = item.Description
			}

			insertQuery := `INSERT INTO feed_items
							(parent_feed_id, feed_title, link, pub_date, content)
							VALUES (?, ?, ?, ?, ?)`
			_, err := tx.Exec(insertQuery, feed.ID, item.Title, item.Link, pubDate, content)
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			newItems++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newItems, nil
}
