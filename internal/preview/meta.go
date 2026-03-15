package preview

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// PageMeta holds extracted page metadata.
type PageMeta struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Favicon     string      `json:"favicon"`
	OG          OpenGraph   `json:"open_graph"`
	Twitter     TwitterCard `json:"twitter_card"`
}

// OpenGraph holds Open Graph metadata.
type OpenGraph struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
	URL         string `json:"url,omitempty"`
	Type        string `json:"type,omitempty"`
	SiteName    string `json:"site_name,omitempty"`
}

// TwitterCard holds Twitter Card metadata.
type TwitterCard struct {
	Card        string `json:"card,omitempty"`
	Site        string `json:"site,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
}

// Fetch retrieves and parses metadata from a URL.
func Fetch(rawURL string, timeoutSec int) (*PageMeta, error) {
	client := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	meta := &PageMeta{}
	extractMeta(doc, meta)

	// Default favicon
	if meta.Favicon == "" {
		meta.Favicon = strings.TrimRight(rawURL, "/") + "/favicon.ico"
	}

	return meta, nil
}

func extractMeta(n *html.Node, meta *PageMeta) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if n.FirstChild != nil {
				meta.Title = strings.TrimSpace(n.FirstChild.Data)
			}
		case "meta":
			handleMetaTag(n, meta)
		case "link":
			handleLinkTag(n, meta)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractMeta(c, meta)
	}
}

func handleMetaTag(n *html.Node, meta *PageMeta) {
	var name, property, content string
	for _, a := range n.Attr {
		switch strings.ToLower(a.Key) {
		case "name":
			name = strings.ToLower(a.Val)
		case "property":
			property = strings.ToLower(a.Val)
		case "content":
			content = a.Val
		}
	}

	if content == "" {
		return
	}

	// Standard meta tags
	switch name {
	case "description":
		meta.Description = content
	}

	// Open Graph
	switch property {
	case "og:title":
		meta.OG.Title = content
	case "og:description":
		meta.OG.Description = content
	case "og:image":
		meta.OG.Image = content
	case "og:url":
		meta.OG.URL = content
	case "og:type":
		meta.OG.Type = content
	case "og:site_name":
		meta.OG.SiteName = content
	}

	// Twitter Card
	key := name
	if key == "" {
		key = property
	}
	switch key {
	case "twitter:card":
		meta.Twitter.Card = content
	case "twitter:site":
		meta.Twitter.Site = content
	case "twitter:title":
		meta.Twitter.Title = content
	case "twitter:description":
		meta.Twitter.Description = content
	case "twitter:image":
		meta.Twitter.Image = content
	}
}

func handleLinkTag(n *html.Node, meta *PageMeta) {
	var rel, href string
	for _, a := range n.Attr {
		switch strings.ToLower(a.Key) {
		case "rel":
			rel = strings.ToLower(a.Val)
		case "href":
			href = a.Val
		}
	}

	if href == "" {
		return
	}

	switch rel {
	case "icon", "shortcut icon":
		meta.Favicon = href
	}
}
