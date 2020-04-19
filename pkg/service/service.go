package service

import (
	"context"
	"sync"

	"github.com/havr/wea/pkg/clients/openweather"
	"github.com/havr/wea/pkg/clients/wiki"
)

// Wea is a generic interface for the application business logic.
type Wea interface {
	LocationDescriptionWithWeather(ctx context.Context, cityName string) (*LocationDescriptionWithWeather, error)
}

// NewWea creates the application business logic instance using the given clients.
func NewWea(ow openweather.Client, wiki wiki.Client) Wea {
	return &defaultWea{
		ow:   ow,
		wiki: wiki,
	}
}

type defaultWea struct {
	ow   openweather.Client
	wiki wiki.Client
}

// LocationDescriptionWithWeather is the struct for location description joined with its weather situation.
type LocationDescriptionWithWeather struct {
	TemperatureCelsius  float64
	WeatherDescription  string
	LocationDescription string
}

// LocationDescriptionWithWeather fetches weather situation and Wikipedia description for the given location and returns them joined together.
func (w defaultWea) LocationDescriptionWithWeather(ctx context.Context, locationName string) (*LocationDescriptionWithWeather, error) {
	var (
		wg          sync.WaitGroup
		description string
		wikiError   error
		weather     *openweather.WeatherSituation
		owError     error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		weather, owError = w.ow.CurrentSituation(ctx, locationName)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		description, wikiError = w.wiki.SimpleExtract(ctx, locationName)
	}()

	wg.Wait()

	if owError != nil {
		return nil, owError
	}
	if wikiError != nil {
		return nil, wikiError
	}

	return &LocationDescriptionWithWeather{
		TemperatureCelsius:  weather.Temperature.Celsius(),
		WeatherDescription:  weather.Description,
		LocationDescription: description,
	}, nil
}
