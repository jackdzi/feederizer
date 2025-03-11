package api

import "time"

// To parse RSS
type RSS struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	PubDate     string        `xml:"pubDate"`
	Description string        `xml:"description"`
	Content     string        `xml:"encoded"`
}

type Feed struct {
	ID       int    `db:"id"`
	FeedName string `db:"feed_name"`
	ApiURL   string `db:"api_url"`
}

type FeedItem struct {
	ID           int       `db:"id"`
	ParentFeedID int       `db:"parent_feed_id"`
	FeedTitle    string    `db:"feed_title"`
	Link         string    `db:"link"`
	PubDate      time.Time `db:"pub_date"`
	Content      string    `db:"content"`
}
