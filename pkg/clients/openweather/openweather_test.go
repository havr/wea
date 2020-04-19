package openweather_test

import (
	"context"
	"os"
	"testing"

	"github.com/havr/wea/pkg/clients/openweather"
	"github.com/stretchr/testify/require"
)

func TestDefaultClient_CurrentSituation_ReturnsResponse_InCaseOfSuccess(t *testing.T) {
	apiKey := os.Getenv("WE_OW_API_KEY") // keep in sync with the app env
	if apiKey == "" {
		t.Skip()
		return
	}

	cli := openweather.NewClient(apiKey)
	sit, err := cli.CurrentSituation(context.Background(), "Tallinn")
	require.NoError(t, err)
	require.NotEmpty(t, sit.Description)
	require.NotZero(t, sit.Temperature)
}

func TestDefaultClient_CurrentSituation_ReturnsError_InCaseOfCityNotFound(t *testing.T) {
	apiKey := os.Getenv("WE_OW_API_KEY") // keep in sync with the app env
	if apiKey == "" {
		t.Skip()
		return
	}

	cli := openweather.NewClient(apiKey)
	_, err := cli.CurrentSituation(context.Background(), "Ponyville")
	require.Equal(t, openweather.ErrLocationNotFound, err)
}
