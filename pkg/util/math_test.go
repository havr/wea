package util_test

import (
	"testing"

	"github.com/havr/wea/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestRound(t *testing.T) {
	require.Equal(t, 123.0, util.Round(123.45678, 0))
	require.Equal(t, 123.46, util.Round(123.45678, 2))
}
