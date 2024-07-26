package anthropic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrAnthropicInvalidRequest = errors.New("invalid request")
	ErrAnthropicUnauthorized   = errors.New("unauthorized: there's an issue with your API key")
	ErrAnthropicForbidden      = errors.New("forbidden: your API key does not have permission to use the specified resource")
	ErrAnthropicRateLimit      = errors.New("your account has hit a rate limit")
	ErrAnthropicInternalServer = errors.New("an unexpected error has occurred internal to Anthropic's systems")

	ErrAnthropicApiKeyRequired = errors.New("apiKey is required")
)

// mapHTTPStatusCodeToError maps an HTTP status code to an error.
func MapHTTPStatusCodeToError(code int) error {
	switch code {
	case http.StatusBadRequest:
		return ErrAnthropicInvalidRequest
	case http.StatusUnauthorized:
		return ErrAnthropicUnauthorized
	case http.StatusForbidden:
		return ErrAnthropicForbidden
	case http.StatusTooManyRequests:
		return ErrAnthropicRateLimit
	case http.StatusInternalServerError:
		return ErrAnthropicInternalServer
	default:
		return errors.New("unknown error occurred")
	}
}

type ErrorResponse struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func MapBodyToError(body io.Reader) error {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("anthropic-go internal error reading response body: %w", err)
	}
	var errorResponse ErrorResponse
	err = json.Unmarshal(bodyBytes, &errorResponse)
	if err != nil {
		return fmt.Errorf("anthropic-go internal error parsing response body: %w", err)
	}

	return fmt.Errorf("%s: %s", errorResponse.Error.Type, errorResponse.Error.Message)
}
