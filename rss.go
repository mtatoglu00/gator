package main

import (
	"context"
	"encoding/xml"
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
	if err != nil {
		log.Println("error when requesting feed", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("error when receiving response", err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response", err)
		return nil, err
	}

	defer resp.Body.Close()

	var data *RSSFeed

	err = xml.Unmarshal(body, &data)
	if err != nil {
		log.Println("error unmarshaling response body", err)
		return nil, err
	}

	data.Channel.Title = html.UnescapeString(data.Channel.Title)
	data.Channel.Description = html.UnescapeString(data.Channel.Description)
	for i, v := range data.Channel.Item {
		data.Channel.Item[i].Description = html.UnescapeString(v.Description)
		data.Channel.Item[i].Title = html.UnescapeString(v.Title)
	}

	return data, nil
}
