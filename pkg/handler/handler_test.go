package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/havr/wea/pkg/clients/openweather"
	"github.com/havr/wea/pkg/clients/wiki"
	"github.com/havr/wea/pkg/handler"
	"github.com/havr/wea/pkg/service"
	"github.com/stretchr/testify/require"
)

func TestHandler_CityInformation_InCaseOfCityNotFound(t *testing.T) {
	testCases := []struct {
		Name          string
		Err           error
		ExpectCode    int
		ExpectMessage string
	}{
		{
			Name:          "OWCityNotFound",
			Err:           openweather.ErrLocationNotFound,
			ExpectCode:    http.StatusNotFound,
			ExpectMessage: openweather.ErrLocationNotFound.Error(),
		},
		{
			Name:          "WikiEntryNotFound",
			Err:           wiki.ErrEntryNotFound,
			ExpectCode:    http.StatusNotFound,
			ExpectMessage: wiki.ErrEntryNotFound.Error(),
		},
		{
			Name:          "InternalError",
			Err:           errors.New("hello"),
			ExpectCode:    http.StatusInternalServerError,
			ExpectMessage: handler.ErrInternal,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			st := WeaStub{
				Err: testCase.Err,
			}

			h := handler.New(st)
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/city-information?name=Hello", nil)
			h.ServeHTTP(resp, req)

			require.Equal(t, testCase.ExpectCode, resp.Code)
			var got handler.Error
			require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &got))
			require.Equal(t, got.Error, testCase.ExpectMessage)
		})
	}
}

func TestHandler_CityInformation_ReturnsBadRequest_InCaseOfNoCitySpecified(t *testing.T) {
	st := WeaStub{
		Err: errors.New("shouldn't be returned at all"),
	}

	h := handler.New(st)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/city-information", nil)
	h.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	var got handler.Error
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &got))
	require.Equal(t, handler.ErrNoCityNameProvided, got.Error)
}

func TestHandler_CityInformation_ReturnsResponse_InCaseOfSuccess(t *testing.T) {
	st := WeaStub{
		Response: service.LocationDescriptionWithWeather{
			TemperatureCelsius:  30.0,
			WeatherDescription:  "nice",
			LocationDescription: "hello",
		},
	}

	h := handler.New(st)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/city-information?name=Hello", nil)
	h.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var got handler.CityInformationResponse
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &got))
	require.Equal(t, st.Response.LocationDescription, got.Description)
	require.Equal(t, st.Response.WeatherDescription, got.WeatherSituation)
	require.Equal(t, st.Response.TemperatureCelsius, got.Temperature)
}

type WeaStub struct {
	Response service.LocationDescriptionWithWeather
	Err      error
}

func (w WeaStub) LocationDescriptionWithWeather(ctx context.Context, cityName string) (*service.LocationDescriptionWithWeather, error) {
	if w.Err != nil {
		return nil, w.Err
	}

	return &w.Response, nil
}
