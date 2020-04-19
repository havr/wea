package openweather

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/havr/wea/pkg/types"
	"github.com/havr/wea/pkg/util"
)

// ErrLocationNotFound is returned when the given location doesn't exist in the OpenWeather database.
var ErrLocationNotFound = errors.New("location not found")

type defaultClient struct {
	apiKey string
}

// Client is the generic interface for OpenWeather clients to implement.
type Client interface {
	CurrentSituation(ctx context.Context, cityName string) (*WeatherSituation, error)
}

// NewClient creates an OpenWeather client that uses its HTTP API using the provided apiKey.
func NewClient(apiKey string) Client {
	return &defaultClient{
		apiKey: apiKey,
	}
}

// WeatherSituation describes a current weather situation.
type WeatherSituation struct {
	Temperature types.Temperature
	Description string
}

type weatherResponse struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

// CurrentSituation queries the current weather situation for the given location.
// It returns ErrLocationNotFound if the given location isn't found. In case of other errors it returns an annotated util.HTTPError.
func (c *defaultClient) CurrentSituation(ctx context.Context, locationName string) (*WeatherSituation, error) {
	var response weatherResponse
	if err := util.GetJSON(ctx, &response, "https://api.openweathermap.org/data/2.5/weather", url.Values{
		"q":     []string{locationName},
		"appid": []string{c.apiKey},
	}); err != nil {
		var httperr util.HTTPError
		if errors.As(err, &httperr) && httperr.StatusCode == http.StatusNotFound {
			return nil, ErrLocationNotFound
		}
		return nil, fmt.Errorf("get current weather: %w", err)
	}

	return &WeatherSituation{
		Temperature: types.Temperature(response.Main.Temp),
		Description: response.Weather[0].Description,
	}, nil
}
