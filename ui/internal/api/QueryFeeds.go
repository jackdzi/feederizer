package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


	type Feed struct {
		FeedName  string `json:"feed_name"`
		FeedTitle string `json:"feed_title"`
		Link      string `json:"link"`
		PubDate   string `json:"pub_date"`
		Content   string `json:"content"`
	}


func GetSubscribedFeedItems() error {
	req, _ := http.NewRequest("GET", "http://localhost:8080/feed", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var feeds []Feed

	err = json.Unmarshal(body, &feeds)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Println("Feed Name:", feed.FeedName)
		fmt.Println("Feed Title:", feed.FeedTitle)
		fmt.Println("Link:", feed.Link)
		fmt.Println("Publication Date:", feed.PubDate)
		fmt.Println("Content:", feed.Content)
		fmt.Println("-----")
	}

	return nil
}
