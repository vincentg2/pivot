package broadcast

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/html"
)

func TestParseFootao(t *testing.T) {
	document, err := html.Parse(strings.NewReader(`<html><body><section id="j240826">
		<div><time itemprop="startDate" content="2026-08-24T18:30">18:30</time>
		<a><img class="im" alt="match DAZN foot programme soir"></a>
		<a href="/bologne-lazio" class="rc"><span itemprop="name">Bologne · Lazio</span></a>
		<span class="ap"><a>Serie A</a></span></div>
		<div><time itemprop="startDate" content="2026-08-24T21:00">21:00</time>
		<a><img class="im" alt="match Canal+ Foot foot programme soir"></a>
		<a href="/psg-om" class="rc"><span itemprop="name">PSG · Marseille</span></a><span class="ap">»</span></div>
	</section></body></html>`))
	if err != nil {
		t.Fatal(err)
	}
	location, _ := time.LoadLocation("Europe/Paris")
	from := time.Date(2026, 8, 24, 0, 0, 0, 0, location)
	items := parseFootao(document, from, from)
	if len(items) != 2 {
		t.Fatalf("expected 2 listings, got %d", len(items))
	}
	if items[0].HomeName != "Bologne" || items[0].AwayName != "Lazio" {
		t.Fatalf("unexpected teams: %#v", items[0])
	}
	if items[0].CompetitionName != "Serie A" || items[1].CompetitionName != "Serie A" {
		t.Fatalf("competition inheritance failed: %#v", items)
	}
	if len(items[1].Channels) != 1 || items[1].Channels[0] != "Canal+ Foot" {
		t.Fatalf("unexpected channels: %#v", items[1].Channels)
	}
}

func TestChooseMatchRequiresBothTeams(t *testing.T) {
	first := uuid.New()
	second := uuid.New()
	item := ImportedListing{HomeName: "Paris Saint-Germain", AwayName: "Olympique de Marseille"}
	selected := chooseMatch(item, []MatchCandidate{{ID: first, HomeName: "Paris Saint-Germain FC", AwayName: "Olympique de Marseille"}, {ID: second, HomeName: "Paris FC", AwayName: "Lyon"}})
	if selected == nil || *selected != first {
		t.Fatalf("expected first match, got %v", selected)
	}
	if selected = chooseMatch(ImportedListing{HomeName: "Paris", AwayName: "Unknown"}, []MatchCandidate{{ID: first, HomeName: "Paris FC", AwayName: "Lyon"}}); selected != nil {
		t.Fatalf("ambiguous listing must stay external")
	}
}

func TestNormalizeFrenchClubAliases(t *testing.T) {
	if nameScore("Bologne", "Bologna FC 1909") < 2 {
		t.Fatal("expected Bologne to match Bologna")
	}
	if nameScore("La Corogne", "RC Deportivo La Coruña") < 2 {
		t.Fatal("expected La Corogne to match Deportivo La Coruña")
	}
}
