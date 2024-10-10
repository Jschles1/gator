package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Jschles1/gator/internal/database"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
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

func parseRSSDate(dateString string) (time.Time, error) {
	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	return time.Parse(layout, dateString)
}

func getNullableString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func outputDescription(ns sql.NullString) string {
	html := getNullableString(ns)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html
	}
	return strings.TrimSpace(doc.Text())
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Handle the case where no feeds are available
			return fmt.Errorf("no feeds available to fetch")
		}
		return err
	}
	markFetchedParams := database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	}
	feedToFetch, err := s.db.MarkFeedFetched(context.Background(), markFetchedParams)
	if err != nil {
		return err
	}
	feed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}
	for _, item := range feed.Channel.Item {
		if item.Title == "" {
			continue
		}
		pubDate, err := parseRSSDate(item.PubDate)
		if err != nil {
			fmt.Println(fmt.Errorf("error parsing publish date: %w", err))
			continue
		}
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			FeedID:      feedToFetch.ID,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			PublishedAt: pubDate,
			Title:       item.Title,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			Url:         item.Link,
		}
		post, err := s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Println(fmt.Errorf("error creating post: %w", err))
			continue
		}
		fmt.Println("Post \"" + post.Title + "\" was created!")
	}
	return nil
}
