package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", feedURL, nil)

	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("error reading response")
	}

	var data *RSSFeed

	err = xml.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln("error unmarshaling response body")
	}

	data.Channel.Title = html.UnescapeString(data.Channel.Title)
	data.Channel.Description = html.UnescapeString(data.Channel.Description)
	for _, v := range data.Channel.Item {
		v.Description = html.UnescapeString(v.Description)
		v.Title = html.UnescapeString(v.Title)
	}

	return data, nil
}

func Agg() {
	feed, err := FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		log.Fatalln("error reading feed")
	}

	fmt.Println(feed.Channel.Title)
	fmt.Println(feed.Channel.Description)
	for _, v := range feed.Channel.Item {
		fmt.Println(v.Title)
		fmt.Println(v.Description)
	}
}
