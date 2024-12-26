// package rss
//
// import (
// 	"github.com/mmcdole/gofeed"
// 	"log"
// )
//
// func FetchRSSFeed(url string) ([]*gofeed.Item, error) {
// 	parser := gofeed.NewParser()
// 	feed, err := parser.ParseURL(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return feed.Items, nil
// }
