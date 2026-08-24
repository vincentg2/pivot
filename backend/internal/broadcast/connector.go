package broadcast

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"golang.org/x/net/html"
)

const footaoURL = "https://www.footao.tv/"

type FootaoConnector struct {
	enabled   bool
	userAgent string
	client    *http.Client
}

func NewFootaoConnector(enabled bool, userAgent string) *FootaoConnector {
	return &FootaoConnector{enabled: enabled, userAgent: userAgent, client: &http.Client{Timeout: 20 * time.Second}}
}

func (c *FootaoConnector) Enabled() bool { return c.enabled }

func (c *FootaoConnector) Fetch(ctx context.Context, from, to time.Time) ([]ImportedListing, error) {
	if !c.enabled {
		return nil, ErrConnectorDisabled
	}
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt*2) * time.Second):
			}
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, footaoURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", c.userAgent)
		req.Header.Set("Accept", "text/html,application/xhtml+xml")
		response, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		if response.StatusCode == http.StatusTooManyRequests || response.StatusCode >= 500 {
			io.Copy(io.Discard, io.LimitReader(response.Body, 4096))
			response.Body.Close()
			lastErr = fmt.Errorf("Footao returned %s", response.Status)
			continue
		}
		if response.StatusCode != http.StatusOK {
			response.Body.Close()
			return nil, fmt.Errorf("Footao returned %s", response.Status)
		}
		document, parseErr := html.Parse(io.LimitReader(response.Body, 5<<20))
		response.Body.Close()
		if parseErr != nil {
			return nil, fmt.Errorf("parse Footao response: %w", parseErr)
		}
		return parseFootao(document, from, to), nil
	}
	return nil, lastErr
}

func parseFootao(document *html.Node, from, to time.Time) []ImportedListing {
	location, _ := time.LoadLocation("Europe/Paris")
	items := make([]ImportedListing, 0)
	for _, section := range elements(document, "section") {
		currentCompetition := ""
		for child := section.FirstChild; child != nil; child = child.NextSibling {
			if child.Type != html.ElementNode || child.Data != "div" {
				continue
			}
			timeNode := firstElement(child, "time", func(n *html.Node) bool { return attr(n, "content") != "" })
			matchNode := firstElement(child, "a", func(n *html.Node) bool { return hasClass(n, "rc") })
			if timeNode == nil || matchNode == nil {
				continue
			}
			startsAt, err := time.ParseInLocation("2006-01-02T15:04", attr(timeNode, "content"), location)
			if err != nil || startsAt.Before(startOfDay(from, location)) || startsAt.After(endOfDay(to, location)) {
				continue
			}
			competitionNode := firstElement(child, "span", func(n *html.Node) bool { return hasClass(n, "ap") })
			if competitionNode != nil {
				value := cleanText(textContent(competitionNode))
				if value != "" && value != "»" {
					currentCompetition = value
				}
			}
			nameNode := firstElement(matchNode, "span", func(n *html.Node) bool { return attr(n, "itemprop") == "name" })
			label := cleanText(textContent(nameNode))
			if label == "" {
				label = cleanText(textContent(matchNode))
			}
			home, away := splitTeams(label)
			channels := make([]string, 0)
			for _, img := range elements(child, "img") {
				if !hasClass(img, "im") {
					continue
				}
				if channel := channelFromAlt(attr(img, "alt")); channel != "" && !contains(channels, channel) {
					channels = append(channels, channel)
				}
			}
			if len(channels) == 0 {
				continue
			}
			kind := "live"
			classes := attr(child, "class") + " " + attr(timeNode, "class")
			if strings.Contains(classes, "df") {
				kind = "delayed"
			}
			if strings.Contains(classes, "rd") {
				kind = "replay"
			}
			href := attr(matchNode, "href")
			var sourceURL *string
			if href != "" {
				if resolved, err := url.Parse(href); err == nil {
					base, _ := url.Parse(footaoURL)
					absoluteURL := base.ResolveReference(resolved)
					if (absoluteURL.Scheme == "http" || absoluteURL.Scheme == "https") && (absoluteURL.Host == "footao.tv" || absoluteURL.Host == "www.footao.tv") {
						absolute := absoluteURL.String()
						sourceURL = &absolute
					}
				}
			}
			keyMaterial := startsAt.Format(time.RFC3339) + "|" + href + "|" + label
			sum := sha256.Sum256([]byte(keyMaterial))
			items = append(items, ImportedListing{SourceKey: hex.EncodeToString(sum[:]), StartsAt: startsAt, HomeName: home, AwayName: away, Label: label, CompetitionName: currentCompetition, Kind: kind, Channels: channels, SourceURL: sourceURL})
		}
	}
	return items
}

func elements(root *html.Node, name string) []*html.Node {
	items := make([]*html.Node, 0)
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == name {
			items = append(items, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(root)
	return items
}
func firstElement(root *html.Node, name string, predicate func(*html.Node) bool) *html.Node {
	var found *html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if found != nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == name && predicate(n) {
			found = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(root)
	return found
}
func attr(node *html.Node, name string) string {
	if node == nil {
		return ""
	}
	for _, a := range node.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}
func hasClass(node *html.Node, class string) bool {
	for _, value := range strings.Fields(attr(node, "class")) {
		if value == class {
			return true
		}
	}
	return false
}
func textContent(node *html.Node) string {
	if node == nil {
		return ""
	}
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
			b.WriteByte(' ')
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(node)
	return b.String()
}
func cleanText(value string) string { return strings.Join(strings.Fields(value), " ") }
func splitTeams(label string) (string, string) {
	parts := strings.SplitN(label, "·", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}
func channelFromAlt(value string) string {
	const prefix = "match "
	const suffix = " foot programme soir"
	if !strings.HasPrefix(value, prefix) || !strings.HasSuffix(value, suffix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(value, prefix), suffix))
}
func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
func startOfDay(value time.Time, location *time.Location) time.Time {
	v := value.In(location)
	return time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, location)
}
func endOfDay(value time.Time, location *time.Location) time.Time {
	return startOfDay(value, location).AddDate(0, 0, 1).Add(-time.Nanosecond)
}

func normalizeName(value string) string {
	value = strings.ToLower(strings.NewReplacer("é", "e", "è", "e", "ê", "e", "ë", "e", "à", "a", "â", "a", "ä", "a", "á", "a", "ç", "c", "î", "i", "ï", "i", "í", "i", "ô", "o", "ö", "o", "ó", "o", "ù", "u", "û", "u", "ü", "u", "ú", "u", "ñ", "n").Replace(value))
	normalized := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, value)
	aliases := map[string]string{
		"bologne":    "bologna",
		"lacorogne":  "deportivolacoruna",
		"intermilan": "internazionale",
	}
	if alias, ok := aliases[normalized]; ok {
		return alias
	}
	return normalized
}
