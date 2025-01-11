package api

import (
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"log"
)

func DbInit(db *sqlx.DB, c *gin.Context) {

	schema := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		password TEXT NULL OR NOT NULL
	);

	CREATE TABLE IF NOT EXISTS feeds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		feed_name TEXT NOT NULL,
		api_url TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS subscriptions (
		user_id INTEGER,
		feed_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS feed_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		parent_feed_id INTEGER,
		feed_title TEXT NOT NULL,
		link TEXT NOT NULL,
		pub_date DATETIME NOT NULL,
		content TEXT,
		FOREIGN KEY (parent_feed_id) REFERENCES feeds(id) ON DELETE CASCADE
	);`

	_, execErr := db.Exec(schema)
	if execErr != nil {
		log.Fatalln(execErr)
	}

	c.JSON(200, gin.H{"message": "Database initialized"})
}
