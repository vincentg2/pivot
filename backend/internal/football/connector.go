package football

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

type FootballDataConnector struct {
	apiKey  string
	baseURL string
	client  *http.Client
	limiter *rate.Limiter
}

func NewFootballDataConnector(apiKey string) *FootballDataConnector {
	return &FootballDataConnector{apiKey: apiKey, baseURL: "https://api.football-data.org/v4", client: &http.Client{Timeout: 20 * time.Second}, limiter: rate.NewLimiter(rate.Every(7*time.Second), 1)}
}
func (c *FootballDataConnector) Enabled() bool { return c.apiKey != "" }

func (c *FootballDataConnector) FetchMatches(ctx context.Context, code string, from, to time.Time) (Season, []Match, error) {
	query := url.Values{"dateFrom": {from.Format("2006-01-02")}, "dateTo": {to.Format("2006-01-02")}}
	var payload struct {
		Matches []struct {
			ID                 int       `json:"id"`
			UTCDate            time.Time `json:"utcDate"`
			Status, Stage      string
			Matchday           *int
			Season             providerSeason
			HomeTeam, AwayTeam providerTeam
			Score              struct {
				FullTime struct{ Home, Away *int } `json:"fullTime"`
			} `json:"score"`
			Goals []struct {
				Minute, InjuryTime *int
				Type               string
				Team               struct{ ID int }
				Scorer             struct{ Name string }
			} `json:"goals"`
		} `json:"matches"`
	}
	if err := c.get(ctx, "/competitions/"+code+"/matches?"+query.Encode(), &payload); err != nil {
		return Season{}, nil, err
	}
	items := make([]Match, 0, len(payload.Matches))
	var season Season
	for _, source := range payload.Matches {
		season = source.Season.model()
		var goals []Goal
		if source.Goals != nil {
			goals = make([]Goal, 0, len(source.Goals))
			for _, goal := range source.Goals {
				if goal.Minute == nil || goal.Scorer.Name == "" {
					continue
				}
				goals = append(goals, Goal{Minute: *goal.Minute, InjuryTime: goal.InjuryTime, Type: goal.Type, TeamProviderID: goal.Team.ID, ScorerName: goal.Scorer.Name})
			}
		}
		items = append(items, Match{ProviderID: source.ID, Season: season, UTCDate: source.UTCDate, Status: source.Status, Stage: source.Stage, Matchday: source.Matchday, Home: source.HomeTeam.ref(), Away: source.AwayTeam.ref(), HomeProviderID: source.HomeTeam.ID, AwayProviderID: source.AwayTeam.ID, HomeScore: source.Score.FullTime.Home, AwayScore: source.Score.FullTime.Away, Goals: goals})
	}
	return season, items, nil
}

func (c *FootballDataConnector) FetchStanding(ctx context.Context, code string) (Season, []StandingRow, error) {
	var payload struct {
		Season    providerSeason `json:"season"`
		Standings []struct {
			Type  string
			Table []struct {
				Position, PlayedGames, Won, Draw, Lost, Points, GoalsFor, GoalsAgainst, GoalDifference int
				Team                                                                                   providerTeam
			}
		} `json:"standings"`
	}
	if err := c.get(ctx, "/competitions/"+code+"/standings", &payload); err != nil {
		return Season{}, nil, err
	}
	rows := make([]StandingRow, 0)
	for _, table := range payload.Standings {
		if table.Type != "TOTAL" || len(rows) > 0 {
			continue
		}
		for _, source := range table.Table {
			if source.Position < 1 {
				continue
			}
			rows = append(rows, StandingRow{Position: source.Position, Club: source.Team.ref(), ClubProviderID: source.Team.ID, Played: source.PlayedGames, Won: source.Won, Drawn: source.Draw, Lost: source.Lost, GoalsFor: source.GoalsFor, GoalsAgainst: source.GoalsAgainst, GoalDifference: source.GoalDifference, Points: source.Points})
		}
	}
	season := payload.Season.model()
	season.Current = true
	return season, rows, nil
}

func (c *FootballDataConnector) get(ctx context.Context, path string, destination any) error {
	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Auth-Token", c.apiKey)
	request.Header.Set("X-Unfold-Goals", "true")
	request.Header.Set("User-Agent", "Pivot/0.3")
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("football-data.org returned status %d", response.StatusCode)
	}
	return json.NewDecoder(response.Body).Decode(destination)
}

type providerSeason struct {
	ID                 int
	StartDate, EndDate string
	Current            bool
}

func (s providerSeason) model() Season {
	start, _ := time.Parse("2006-01-02", s.StartDate)
	end, _ := time.Parse("2006-01-02", s.EndDate)
	today := time.Now().UTC().Truncate(24 * time.Hour)
	current := s.Current || (!start.IsZero() && !end.IsZero() && !today.Before(start) && !today.After(end))
	return Season{ProviderID: s.ID, StartDate: start, EndDate: end, Current: current}
}

type providerTeam struct {
	ID                          int
	Name, ShortName, TLA, Crest string
}

func (t providerTeam) ref() ClubRef {
	return ClubRef{Name: t.Name, ShortName: t.ShortName, TLA: t.TLA, CrestURL: optionalString(t.Crest)}
}
func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
