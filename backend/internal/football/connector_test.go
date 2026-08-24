package football

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestFetchMatchesRequestsAndMapsGoals(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("X-Unfold-Goals") != "true" {
			t.Fatal("expected unfolded goals request header")
		}
		body := `{"matches":[{"id":99,"utcDate":"2026-08-23T19:00:00Z","status":"FINISHED","stage":"REGULAR_SEASON","matchday":1,"season":{"id":8,"startDate":"2026-08-01","endDate":"2027-05-30"},"homeTeam":{"id":1,"name":"Paris","shortName":"PSG","tla":"PSG"},"awayTeam":{"id":2,"name":"Rennes","shortName":"Rennes","tla":"REN"},"score":{"fullTime":{"home":2,"away":1}},"goals":[{"minute":12,"injuryTime":null,"type":"REGULAR","team":{"id":1},"scorer":{"name":"A. Striker"}},{"minute":90,"injuryTime":3,"type":"PENALTY","team":{"id":2},"scorer":{"name":"B. Forward"}}]}]}`
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})}

	connector := NewFootballDataConnector("test-key")
	connector.baseURL = "https://football-data.test"
	connector.client = client
	connector.limiter = rate.NewLimiter(rate.Inf, 1)
	_, matches, err := connector.FetchMatches(context.Background(), "FL1", time.Now(), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 || len(matches[0].Goals) != 2 {
		t.Fatalf("unexpected goals payload: %#v", matches)
	}
	if matches[0].Goals[1].ScorerName != "B. Forward" || matches[0].Goals[1].InjuryTime == nil || *matches[0].Goals[1].InjuryTime != 3 {
		t.Fatalf("unexpected second goal: %#v", matches[0].Goals[1])
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return fn(request) }
