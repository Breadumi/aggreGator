package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rss_xml, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("Raw XML:", string(rss_xml)) // temporary debug line

	var rss RSSFeed
	err = xml.Unmarshal(rss_xml, &rss)
	if err != nil {
		return nil, err
	}

	// clean the XML text for titles and descriptions
	rss.Channel.Title = cleanXML(rss.Channel.Title)
	rss.Channel.Description = cleanXML(rss.Channel.Description)
	for i := range rss.Channel.Item {
		rss.Channel.Item[i].Title = cleanXML(rss.Channel.Item[i].Title)
		rss.Channel.Item[i].Description = cleanXML(rss.Channel.Item[i].Description)
	}

	return &rss, nil
}

func cleanXML(s string) string {
	return html.UnescapeString(s)
}
