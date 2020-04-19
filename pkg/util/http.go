package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPError is intended to represent a HTTP response with non 2xx status code.
type HTTPError struct {
	StatusCode int
	Body       []byte
}

// Error returns the error details (status code and the response body) as a string.
func (e HTTPError) Error() string {
	return fmt.Sprintf("%s: %s", http.StatusText(e.StatusCode), string(e.Body))
}

// GetJSON performs the GET request to the given URL joined with the given query params.
// In case of success it unmarshals the response into the target.
// In case of request failure it returns an annotated original error.
// In case if response doesn't contain 2xx status code, it returns an HTTPError with actual status code and a raw response.
func GetJSON(ctx context.Context, target interface{}, url string, query url.Values) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"?"+query.Encode(), nil)
	if err != nil {
		return fmt.Errorf("create request with context: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return HTTPError{
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("unmarshal response JSON: %w", err)
	}

	return nil
}
