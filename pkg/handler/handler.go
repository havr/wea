package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/havr/wea/pkg/clients/openweather"
	"github.com/havr/wea/pkg/clients/wiki"
	"github.com/havr/wea/pkg/service"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// ErrInternal is a generic error to represent internal server faults.
const ErrInternal = "internal server error"

// ErrNoCityNameProvided serves to indicate that the location name isn't provided.
// I'd suggest to replace it by some generic `ValidationError` error in production.
const ErrNoCityNameProvided = "required parameter 'name' isn`t provided"

// New returns a HTTP handler that serves the given Wea.
func New(wea service.Wea) http.Handler {
	h := handler{wea: wea}
	r := httprouter.New()
	r.GET("/city-information", h.cityInformation)
	return r
}

type handler struct {
	wea service.Wea
}

// CityInformationResponse defines the structure of a JSON response for the the city information request.
type CityInformationResponse struct {
	Description      string  `json:"description"`
	WeatherSituation string  `json:"weatherSituation"`
	Temperature      float64 `json:"temperature"`
}

func (h *handler) cityInformation(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	cityName := r.URL.Query().Get("name")
	if cityName == "" {
		writeError(w, http.StatusBadRequest, ErrNoCityNameProvided)
		return
	}

	sit, err := h.wea.LocationDescriptionWithWeather(r.Context(), cityName)
	if err != nil {
		if errors.Is(err, wiki.ErrEntryNotFound) || errors.Is(err, openweather.ErrLocationNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, context.Canceled) {
			return
		}

		writeInternalError(w)
		zap.L().Error("get city description with weather", zap.String("cityName", cityName), zap.Error(err))
		return
	}

	writeResponse(w, CityInformationResponse{
		Description:      sit.LocationDescription,
		WeatherSituation: sit.WeatherDescription,
		Temperature:      sit.TemperatureCelsius,
	})
}

func writeResponse(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeInternalError(w)
		zap.L().Error("encode response", zap.Reflect("response", resp), zap.Error(err))
	}
}

func writeInternalError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, ErrInternal)
}

// Error defines the structure of an API error response.
type Error struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(Error{
		Error: message,
	})
	if err != nil {
		zap.L().Error("unable to marshal error response", zap.Error(err))
		return
	}
	w.WriteHeader(statusCode)
	if _, err := w.Write(resp); err != nil {
		zap.L().Error("unable to write error response", zap.Error(err))
		return
	}
}
