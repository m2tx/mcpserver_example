package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

var (
	reDDGLink    = regexp.MustCompile(`<a[^>]+class="result__a"[^>]+href="([^"]+)"[^>]*>([\s\S]*?)</a>`)
	reDDGSnippet = regexp.MustCompile(`<a[^>]+class="result__snippet"[^>]*>([\s\S]*?)</a>`)
	reStripTags  = regexp.MustCompile(`<[^>]+>`)
)

type SearchResult struct {
	Title   string
	URL     string
	Snippet string
}

// defaultMinInterval is the minimum time between DuckDuckGo requests to avoid rate limiting.
const defaultMinInterval = time.Second

type DuckDuckGoClient struct {
	httpClient  *http.Client
	mu          sync.Mutex
	lastRequest time.Time
	minInterval time.Duration
}

func NewDuckDuckGoClient() *DuckDuckGoClient {
	return &DuckDuckGoClient{
		httpClient:  &http.Client{},
		minInterval: defaultMinInterval,
	}
}

func (c *DuckDuckGoClient) wait(ctx context.Context) error {
	c.mu.Lock()
	now := time.Now()
	elapsed := now.Sub(c.lastRequest)
	var wait time.Duration
	if elapsed < c.minInterval {
		wait = c.minInterval - elapsed
		c.lastRequest = now.Add(wait)
	} else {
		c.lastRequest = now
	}
	c.mu.Unlock()

	if wait == 0 {
		return nil
	}
	select {
	case <-time.After(wait):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *DuckDuckGoClient) Search(ctx context.Context, query string, count int) ([]SearchResult, error) {
	if err := c.wait(ctx); err != nil {
		return nil, err
	}

	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return extractResults(string(body), count), nil
}

func extractResults(html string, count int) []SearchResult {
	matches := reDDGLink.FindAllStringSubmatch(html, count+5)
	snippetMatches := reDDGSnippet.FindAllStringSubmatch(html, count+5)

	results := make([]SearchResult, 0, min(len(matches), count))
	for i := range min(len(matches), count) {
		urlStr := matches[i][1]
		title := strings.TrimSpace(stripTags(matches[i][2]))

		if strings.Contains(urlStr, "uddg=") {
			if u, err := url.QueryUnescape(urlStr); err == nil {
				if _, after, ok := strings.Cut(u, "uddg="); ok {
					urlStr = after
				}
			}
		}

		var snippet string
		if i < len(snippetMatches) {
			snippet = strings.TrimSpace(stripTags(snippetMatches[i][1]))
		}

		results = append(results, SearchResult{
			Title:   title,
			URL:     urlStr,
			Snippet: snippet,
		})
	}

	return results
}

func stripTags(s string) string {
	return reStripTags.ReplaceAllString(s, "")
}
