package wiki

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/havr/wea/pkg/util"
)

// ErrEntryNotFound is returned when the given Wikipedia entry with doesn't exist.
var ErrEntryNotFound = errors.New("wikipedia entry not found")

// Client is the generic interface for Wikipedia clients.
type Client interface {
	SimpleExtract(ctx context.Context, name string) (string, error)
}

// NewClient creates a Wikipedia API client.
func NewClient() Client {
	return &defaultClient{}
}

type defaultClient struct{}

type extractResponse struct {
	Query struct {
		Pages map[string]struct {
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

// SimpleExtract returns the introduction paragraph for the given topic. It doesn't resolve ambiguities, so it may return the
// '$topic may refer to:' article instead of the concrete one.
// It returns ErrEntryNotFound if the topic isn't found.
func (d *defaultClient) SimpleExtract(ctx context.Context, topicName string) (string, error) {
	var response extractResponse
	if err := util.GetJSON(ctx, &response, "https://en.wikipedia.org/w/api.php", url.Values{
		"exintro":     []string{"true"},
		"explaintext": []string{"true"},
		"action":      []string{"query"},
		"prop":        []string{"extracts"},
		"titles":      []string{topicName},
		"format":      []string{"json"},
	}); err != nil {
		return "", fmt.Errorf("get wiki extract: %w", err)
	}

	if _, ok := response.Query.Pages["-1"]; ok {
		return "", ErrEntryNotFound
	}

	var firstDescription string
	for _, p := range response.Query.Pages {
		firstDescription = p.Extract
		break
	}

	return firstDescription, nil
}
