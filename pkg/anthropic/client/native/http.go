// Package client contains the HTTP client and related functionality for the anthropic package.
package native

import (
	"fmt"
	"net/http"

	"github.com/ehsanul/anthropic-go/v3/pkg/anthropic"
)

const (
	// AnthropicAPIVersion is the version of the Anthropics API that this client is compatible with.
	AnthropicAPIVersion = "2023-06-01"
	// AnthropicAPIMessagesBeta is the beta version of the Anthropics API that enables the messages endpoint.
	AnthropicAPIMessagesBeta = "messages-2023-12-15"
)

// doRequest sends an HTTP request and returns the response, handling any non-OK HTTP status codes.
func (c *Client) doRequest(request *http.Request) (*http.Response, error) {
	request.Header.Add("anthropic-version", AnthropicAPIVersion)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		mainErr := anthropic.MapHTTPStatusCodeToError(response.StatusCode)
		innerErr := anthropic.MapBodyToError(response.Body)
		if innerErr != nil {
			return nil, fmt.Errorf("%w: %w", mainErr, innerErr)
		}
		return nil, mainErr
	}

	return response, nil
}
