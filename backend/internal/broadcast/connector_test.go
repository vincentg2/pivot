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

func TestParseFootaoUsesSectionDateForFutureListings(t *testing.T) {
	document, err := html.Parse(strings.NewReader(`<html><body><section>
		<h2><a href="match-foot.php?v=vendredi-28-aout-2026&amp;jr=28&amp;ms=08&amp;an=2026">vendredi 28 août</a></h2>
		<div><time>20:45</time>
		<a><img class="im" alt="match Ligue 1+ foot programme soir"></a>
		<a href="lille-psg-chaine-tv-diffusion-heure" class="rc"><span itemprop="name">Lille · Paris PSG</span></a>
		<span class="ap"><a>Ligue 1</a></span></div>
	</section></body></html>`))
	if err != nil {
		t.Fatal(err)
	}
	location, _ := time.LoadLocation("Europe/Paris")
	from := time.Date(2026, 8, 24, 0, 0, 0, 0, location)
	items := parseFootao(document, from, from.AddDate(0, 2, 0))
	if len(items) != 1 {
		t.Fatalf("expected future listing, got %d", len(items))
	}
	if got := items[0].StartsAt.In(location).Format("2006-01-02 15:04"); got != "2026-08-28 20:45" {
		t.Fatalf("unexpected future start: %s", got)
	}
	if items[0].Channels[0] != "Ligue 1+" || items[0].HomeName != "Lille" || items[0].AwayName != "Paris PSG" {
		t.Fatalf("unexpected future listing: %#v", items[0])
	}

	changed := strings.ReplaceAll(`<html><body><section>
		<h2><a href="match-foot.php?v=vendredi-28-aout-2026&amp;jr=28&amp;ms=08&amp;an=2026">vendredi 28 août</a></h2>
		<div><time>21:00</time><a><img class="im" alt="match Ligue 1+ foot programme soir"></a>
		<a href="lille-psg-chaine-tv-diffusion-heure" class="rc"><span itemprop="name">Lille · Paris PSG</span></a></div>
	</section></body></html>`, "\t", "")
	changedDocument, err := html.Parse(strings.NewReader(changed))
	if err != nil {
		t.Fatal(err)
	}
	changedItems := parseFootao(changedDocument, from, from.AddDate(0, 2, 0))
	if len(changedItems) != 1 || changedItems[0].SourceKey != items[0].SourceKey {
		t.Fatal("a schedule change must retain the Footao source identity")
	}
}

func TestCollectionWindowSpansTwoCalendarMonths(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Paris")
	from, to := collectionWindow(time.Date(2026, 8, 24, 15, 30, 0, 0, location))
	if got := from.Format("2006-01-02 15:04"); got != "2026-08-24 00:00" {
		t.Fatalf("unexpected collection start: %s", got)
	}
	if got := to.Format("2006-01-02 15:04"); got != "2026-10-24 00:00" {
		t.Fatalf("unexpected collection end: %s", got)
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
	if nameScore("Paris PSG", "Paris Saint-Germain FC") != 3 {
		t.Fatal("expected Paris PSG to match Paris Saint-Germain FC")
	}
}
