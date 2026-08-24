package news

import (
	"net/url"
	"strings"
	"testing"
)

func TestParseRSSAndRejectUnsafeArticleLinks(t *testing.T) {
	base, _ := url.Parse("https://club.example/feed.xml")
	items, err := parseFeed(strings.NewReader(`<?xml version="1.0"?><rss><channel>
		<item><title> First   team update </title><link>https://club.example/news/one</link><guid>one</guid><pubDate>Mon, 24 Aug 2026 08:00:00 +0200</pubDate></item>
		<item><title>Unsafe</title><link>javascript:alert(1)</link><guid>two</guid><pubDate>Mon, 24 Aug 2026 09:00:00 +0200</pubDate></item>
	</channel></rss>`), base)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one safe item, got %d", len(items))
	}
	if items[0].Title != "First team update" {
		t.Fatalf("unexpected title %q", items[0].Title)
	}
}

func TestParseAtom(t *testing.T) {
	base, _ := url.Parse("https://club.example/feed.xml")
	items, err := parseFeed(strings.NewReader(`<?xml version="1.0"?><feed xmlns="http://www.w3.org/2005/Atom"><entry><title>Match report</title><id>tag:club.example,2026:1</id><updated>2026-08-24T10:00:00+02:00</updated><link rel="alternate" href="/report"/></entry></feed>`), base)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].LinkURL != "https://club.example/report" {
		t.Fatalf("unexpected Atom items: %#v", items)
	}
}

func TestValidateRemoteURL(t *testing.T) {
	unsafe := []string{"http://localhost/feed", "http://127.0.0.1/feed", "file:///tmp/feed", "https://user:pass@example.com/feed"}
	for _, raw := range unsafe {
		parsed, _ := url.Parse(raw)
		if validateRemoteURL(parsed) == nil {
			t.Fatalf("expected %s to be rejected", raw)
		}
	}
	safe, _ := url.Parse("https://club.example/feed.xml")
	if err := validateRemoteURL(safe); err != nil {
		t.Fatalf("public feed rejected: %v", err)
	}
}
