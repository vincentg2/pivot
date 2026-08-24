package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FootballDataConnector struct {
	apiKey string
	client *http.Client
}

func NewFootballDataConnector(apiKey string) *FootballDataConnector {
	return &FootballDataConnector{apiKey: apiKey, client: &http.Client{Timeout: 15 * time.Second}}
}

func (c *FootballDataConnector) Enabled() bool { return c.apiKey != "" }

func (c *FootballDataConnector) FetchCompetition(ctx context.Context, code string) (Competition, []Club, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.football-data.org/v4/competitions/"+code+"/teams", nil)
	if err != nil {
		return Competition{}, nil, err
	}
	request.Header.Set("X-Auth-Token", c.apiKey)
	request.Header.Set("User-Agent", "Pivot/0.2")
	response, err := c.client.Do(request)
	if err != nil {
		return Competition{}, nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Competition{}, nil, fmt.Errorf("football-data.org returned status %d", response.StatusCode)
	}
	var payload struct {
		Competition struct {
			Code, Name, Emblem string
			Area               struct{ Name string }
		} `json:"competition"`
		Teams []struct {
			ID                                          int
			Name, ShortName, TLA, Crest, Website, Venue string
		} `json:"teams"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return Competition{}, nil, err
	}
	competition := Competition{ProviderCode: payload.Competition.Code, Name: payload.Competition.Name, Country: payload.Competition.Area.Name, EmblemURL: optionalString(payload.Competition.Emblem)}
	clubs := make([]Club, 0, len(payload.Teams))
	for _, team := range payload.Teams {
		clubs = append(clubs, Club{ProviderID: team.ID, Name: team.Name, ShortName: team.ShortName, TLA: team.TLA, CrestURL: optionalString(team.Crest), WebsiteURL: optionalString(team.Website), Venue: team.Venue})
	}
	return competition, clubs, nil
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
