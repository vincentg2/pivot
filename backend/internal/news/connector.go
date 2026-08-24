package news

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
)

type RSSConnector struct {
	client    *http.Client
	userAgent string
}

func NewRSSConnector(userAgent string) *RSSConnector {
	dialer := &net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}
	transport := &http.Transport{DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}
		ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err
		}
		for _, resolved := range ips {
			if unsafeIP(resolved.IP) {
				continue
			}
			return dialer.DialContext(ctx, network, net.JoinHostPort(resolved.IP.String(), port))
		}
		return nil, ErrUnsafeFeedURL
	}}
	client := &http.Client{Transport: transport, Timeout: 15 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) > 5 {
			return fmt.Errorf("too many feed redirects")
		}
		return validateRemoteURL(req.URL)
	}}
	return &RSSConnector{client: client, userAgent: userAgent}
}

func (c *RSSConnector) Fetch(ctx context.Context, feed Feed) ([]ImportedItem, error) {
	parsed, err := url.Parse(feed.URL)
	if err != nil {
		return nil, ErrUnsafeFeedURL
	}
	if err = validateRemoteURL(parsed); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/rss+xml,application/atom+xml,application/xml,text/xml")
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("official feed returned %s", response.Status)
	}
	return parseFeed(io.LimitReader(response.Body, 2<<20), response.Request.URL)
}
func parseFeed(reader io.Reader, baseURL *url.URL) ([]ImportedItem, error) {
	var document feedDocument
	if err := xml.NewDecoder(reader).Decode(&document); err != nil {
		return nil, fmt.Errorf("parse official feed: %w", err)
	}
	items := make([]ImportedItem, 0)
	rssItems := document.Channel.Items
	if len(rssItems) == 0 {
		rssItems = document.Items
	}
	for _, item := range rssItems {
		published := parseFeedTime(first(item.PubDate, item.Date))
		link := resolveArticleURL(baseURL, item.Link)
		if published.IsZero() || link == "" || strings.TrimSpace(item.Title) == "" {
			continue
		}
		id := first(strings.TrimSpace(item.GUID), link)
		items = append(items, ImportedItem{ExternalID: stableID(id), Title: truncate(clean(item.Title), 300), PublishedAt: published, LinkURL: link})
	}
	for _, entry := range document.Entries {
		published := parseFeedTime(first(entry.Published, entry.Updated))
		link := ""
		for _, candidate := range entry.Links {
			if candidate.Rel == "" || candidate.Rel == "alternate" {
				link = candidate.Href
				break
			}
		}
		link = resolveArticleURL(baseURL, link)
		if published.IsZero() || link == "" || strings.TrimSpace(entry.Title) == "" {
			continue
		}
		items = append(items, ImportedItem{ExternalID: stableID(first(strings.TrimSpace(entry.ID), link)), Title: truncate(clean(entry.Title), 300), PublishedAt: published, LinkURL: link})
	}
	return items, nil
}

type feedDocument struct {
	Channel struct {
		Items []rssItem `xml:"item"`
	} `xml:"channel"`
	Items   []rssItem   `xml:"item"`
	Entries []atomEntry `xml:"entry"`
}
type rssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	GUID    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
	Date    string `xml:"date"`
}
type atomEntry struct {
	Title     string `xml:"title"`
	ID        string `xml:"id"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
	Links     []struct {
		Href string `xml:"href,attr"`
		Rel  string `xml:"rel,attr"`
	} `xml:"link"`
}

func validateRemoteURL(value *url.URL) error {
	if value == nil || (value.Scheme != "http" && value.Scheme != "https") || value.Hostname() == "" || value.User != nil {
		return ErrUnsafeFeedURL
	}
	if ip := net.ParseIP(value.Hostname()); ip != nil && unsafeIP(ip) {
		return ErrUnsafeFeedURL
	}
	host := strings.ToLower(value.Hostname())
	if host == "localhost" || strings.HasSuffix(host, ".localhost") || strings.HasSuffix(host, ".local") {
		return ErrUnsafeFeedURL
	}
	return nil
}
func unsafeIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() || !ip.IsGlobalUnicast() {
		return true
	}
	address, ok := netip.AddrFromSlice(ip)
	if !ok {
		return true
	}
	address = address.Unmap()
	for _, raw := range []string{"100.64.0.0/10", "192.0.0.0/24", "192.0.2.0/24", "198.18.0.0/15", "198.51.100.0/24", "203.0.113.0/24", "2001:db8::/32"} {
		if netip.MustParsePrefix(raw).Contains(address) {
			return true
		}
	}
	return false
}
func resolveArticleURL(baseURL *url.URL, raw string) string {
	value, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return ""
	}
	if baseURL != nil {
		value = baseURL.ResolveReference(value)
	}
	if validateRemoteURL(value) != nil {
		return ""
	}
	return value.String()
}
func parseFeedTime(value string) time.Time {
	layouts := []string{time.RFC1123Z, time.RFC1123, time.RFC3339, time.RFC3339Nano, time.RFC822Z, time.RFC822, "Mon, 2 Jan 2006 15:04:05 MST"}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, strings.TrimSpace(value)); err == nil {
			return parsed
		}
	}
	return time.Time{}
}
func first(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
func stableID(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
func clean(value string) string { return strings.Join(strings.Fields(value), " ") }
func truncate(value string, max int) string {
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max])
}
