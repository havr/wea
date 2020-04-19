package wiki_test

import (
	"context"
	"testing"

	"github.com/havr/wea/pkg/clients/wiki"
	"github.com/stretchr/testify/require"
)

func TestDefaultClient_SimpleExtract_ReturnsResponse_InCaseOfSuccess(t *testing.T) {
	cli := wiki.NewClient()
	resp, err := cli.SimpleExtract(context.Background(), "Wikipedia")
	require.NoError(t, err)
	require.NotEmpty(t, resp)
}

func TestDefaultClient_SimpleExtract_ReturnsError_InCaseOfEntryNotFound(t *testing.T) {
	cli := wiki.NewClient()
	_, err := cli.SimpleExtract(context.Background(), "ThisEntryDefinitelyDoesNotExist")
	require.Error(t, err, wiki.ErrEntryNotFound)
}
