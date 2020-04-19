package service_test

import (
	"context"
	"testing"

	"github.com/havr/wea/pkg/clients/openweather"
	"github.com/havr/wea/pkg/clients/wiki"
	"github.com/havr/wea/pkg/service"
	"github.com/havr/wea/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestWea_CityDescriptionWithWeather_ReturnsResponse_InCaseOfSuccess(t *testing.T) {
	ow := &OWClientStub{
		Response: openweather.WeatherSituation{
			Temperature: types.Temperature(300.0),
			Description: "nice",
		},
	}
	wi := &WikiClientStub{
		Response: "hello",
	}
	wea := service.NewWea(ow, wi)
	desc, err := wea.LocationDescriptionWithWeather(context.Background(), "hello")
	require.NoError(t, err)
	require.Equal(t, desc.TemperatureCelsius, ow.Response.Temperature.Celsius())
	require.Equal(t, desc.WeatherDescription, ow.Response.Description)
	require.Equal(t, desc.LocationDescription, wi.Response)
}

func TestWea_LocationDescriptionWithWeather_ReturnsOWError_InCaseOfBothErrors(t *testing.T) {
	ow := &OWClientStub{
		Err: openweather.ErrLocationNotFound,
	}
	wi := &WikiClientStub{
		Err: wiki.ErrEntryNotFound,
	}
	wea := service.NewWea(ow, wi)
	_, err := wea.LocationDescriptionWithWeather(context.Background(), "hello")
	require.Error(t, ow.Err, err)
}

func TestWea_LocationDescriptionWithWeather_ReturnsError_InCaseOfOWError(t *testing.T) {
	ow := &OWClientStub{
		Err: openweather.ErrLocationNotFound,
	}
	wi := &WikiClientStub{
		Response: "hello",
	}
	wea := service.NewWea(ow, wi)
	_, err := wea.LocationDescriptionWithWeather(context.Background(), "hello")
	require.Error(t, ow.Err, err)
}

func TestWea_LocationDescriptionWithWeather_ReturnsError_InCaseOfWikiError(t *testing.T) {
	ow := &OWClientStub{
		Response: openweather.WeatherSituation{
			Temperature: types.Temperature(300.0),
			Description: "nice",
		},
	}
	wi := &WikiClientStub{
		Err: wiki.ErrEntryNotFound,
	}
	wea := service.NewWea(ow, wi)
	_, err := wea.LocationDescriptionWithWeather(context.Background(), "hello")
	require.Error(t, wi.Err, err)
}

type OWClientStub struct {
	Err      error
	Response openweather.WeatherSituation
}

func (s *OWClientStub) CurrentSituation(ctx context.Context, cityName string) (*openweather.WeatherSituation, error) {
	if s.Err != nil {
		return nil, s.Err
	}

	return &s.Response, nil
}

type WikiClientStub struct {
	Err      error
	Response string
}

func (s *WikiClientStub) SimpleExtract(ctx context.Context, name string) (string, error) {
	if s.Err != nil {
		return "", s.Err
	}

	return s.Response, nil
}
