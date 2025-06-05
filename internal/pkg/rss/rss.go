package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"

	"github.com/charankamal20/gator/internal/pkg/models"
)


func FetchFeed(ctx context.Context, feedUrl string) (*models.RSSFeed, error) {

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "gator")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch feed: %s", resp.Status)
	}

	var result models.RSSFeed
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RSS feed: %w", err)
	}

	result.Channel.Title = html.UnescapeString(result.Channel.Title)
	result.Channel.Description = html.UnescapeString(result.Channel.Description)

	for i, item := range result.Channel.Item {
		result.Channel.Item[i].Title = html.UnescapeString(item.Title)
		result.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &result, nil
}
