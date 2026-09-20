// Package sortedtravel is the official Go client for the Sorted Travel REST API.
//
// Homepage: https://sorted.travel
package sortedtravel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	DefaultBaseURL = "https://sorted.travel"
	DefaultTimeout = 30 * time.Second
	Version        = "0.1.1"
	userAgent      = "sorted-travel-go/" + Version + " (+https://sorted.travel)"
)

// SortedTravelError is the base error type for this SDK.
type SortedTravelError struct {
	Message string
}

func (err *SortedTravelError) Error() string {
	return err.Message
}

// APIError is returned for non-2xx HTTP responses.
type APIError struct {
	Status int
	Body   interface{}
}

func (err *APIError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", err.Status, truncateBody(err.Body, 300))
}

// TransportRequest is the input for an injectable HTTP transport.
type TransportRequest struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    []byte
}

// TransportResponse is the output from an injectable HTTP transport.
type TransportResponse struct {
	Status      int
	ContentType string
	Text        string
}

// Transport performs one HTTP round trip for tests or custom backends.
type Transport func(ctx context.Context, request TransportRequest, timeout time.Duration) (TransportResponse, error)

// Client is a thin wrapper around the Sorted Travel REST API.
type Client struct {
	apiKey      string
	baseURL     string
	timeout     time.Duration
	transport   Transport
	envProvided bool
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithAPIKey sets the Bearer token used for Authorization.
func WithAPIKey(apiKey string) ClientOption {
	return func(client *Client) {
		client.apiKey = apiKey
	}
}

// WithBaseURL overrides the API origin.
func WithBaseURL(baseURL string) ClientOption {
	return func(client *Client) {
		client.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithTimeout overrides the per-request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.timeout = timeout
	}
}

// WithTransport replaces the default net/http transport.
func WithTransport(transport Transport) ClientOption {
	return func(client *Client) {
		client.transport = transport
	}
}

// WithEnv loads SORTED_TRAVEL_API_KEY and SORTED_TRAVEL_BASE_URL from env.
// When set, process environment variables are not read.
func WithEnv(env map[string]string) ClientOption {
	return func(client *Client) {
		client.envProvided = true
		if apiKey := env["SORTED_TRAVEL_API_KEY"]; apiKey != "" {
			client.apiKey = apiKey
		}
		if baseURL := env["SORTED_TRAVEL_BASE_URL"]; baseURL != "" {
			client.baseURL = strings.TrimRight(baseURL, "/")
		}
	}
}

// NewClient builds a REST client using options and process environment defaults.
func NewClient(options ...ClientOption) *Client {
	client := &Client{
		baseURL:   DefaultBaseURL,
		timeout:   DefaultTimeout,
		transport: defaultTransport,
	}
	for _, option := range options {
		option(client)
	}
	if !client.envProvided {
		if client.apiKey == "" {
			client.apiKey = os.Getenv("SORTED_TRAVEL_API_KEY")
		}
		if configuredBase := os.Getenv("SORTED_TRAVEL_BASE_URL"); configuredBase != "" && client.baseURL == DefaultBaseURL {
			client.baseURL = strings.TrimRight(configuredBase, "/")
		}
	}
	return client
}

// Status calls GET /api/v1/status.
func (client *Client) Status(ctx context.Context) (interface{}, error) {
	return client.Get(ctx, "/api/v1/status", nil)
}

// Sandbox calls GET /sandbox.
func (client *Client) Sandbox(ctx context.Context) (interface{}, error) {
	return client.Get(ctx, "/sandbox", nil)
}

// ListDestinations calls GET /api/v1/destinations with cursor pagination.
func (client *Client) ListDestinations(ctx context.Context, cursor *string, limit *int) (interface{}, error) {
	params := map[string]interface{}{}
	if cursor != nil {
		params["cursor"] = *cursor
	}
	if limit != nil {
		params["limit"] = *limit
	}
	return client.Get(ctx, "/api/v1/destinations", params)
}

// CreateAPIKey calls POST /api/v1/api-keys.
func (client *Client) CreateAPIKey(ctx context.Context) (interface{}, error) {
	return client.Post(ctx, "/api/v1/api-keys", map[string]interface{}{})
}

// CreateJob calls POST /api/v1/jobs.
func (client *Client) CreateJob(ctx context.Context, operation string, body map[string]interface{}) (interface{}, error) {
	payload := map[string]interface{}{"operation": operation}
	for key, value := range body {
		payload[key] = value
	}
	return client.Post(ctx, "/api/v1/jobs", payload)
}

// GetJob calls GET /api/v1/jobs/{job_id}.
func (client *Client) GetJob(ctx context.Context, jobID string) (interface{}, error) {
	encodedJobID := url.PathEscape(jobID)
	return client.Get(ctx, "/api/v1/jobs/"+encodedJobID, nil)
}

// Get issues a GET request to a host-relative REST path.
func (client *Client) Get(ctx context.Context, path string, params map[string]interface{}) (interface{}, error) {
	return client.request(ctx, http.MethodGet, path, params, nil)
}

// Post issues a POST request with a JSON body to a host-relative REST path.
func (client *Client) Post(ctx context.Context, path string, body interface{}) (interface{}, error) {
	return client.request(ctx, http.MethodPost, path, nil, body)
}

func (client *Client) request(
	ctx context.Context,
	method string,
	path string,
	params map[string]interface{},
	body interface{},
) (interface{}, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, errors.New("REST paths must start with '/'")
	}
	requestURL := client.baseURL + path
	if len(params) > 0 {
		query := url.Values{}
		for key, value := range params {
			query.Set(key, stringifyQueryValue(value))
		}
		requestURL += "?" + query.Encode()
	}
	headers := client.headers(nil)
	var encodedBody []byte
	if body != nil {
		headers["Content-Type"] = "application/json"
		marshaledBody, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			return nil, marshalErr
		}
		encodedBody = marshaledBody
	}
	response, err := client.transport(ctx, TransportRequest{
		URL:     requestURL,
		Method:  method,
		Headers: headers,
		Body:    encodedBody,
	}, client.timeout)
	if err != nil {
		return nil, err
	}
	return decodeResponse(response)
}

func (client *Client) headers(extra map[string]string) map[string]string {
	headers := map[string]string{
		"User-Agent": userAgent,
		"Accept":     "application/json",
	}
	for key, value := range extra {
		headers[key] = value
	}
	if client.apiKey != "" {
		headers["Authorization"] = "Bearer " + client.apiKey
	}
	return headers
}

func defaultTransport(ctx context.Context, request TransportRequest, timeout time.Duration) (TransportResponse, error) {
	httpClient := &http.Client{Timeout: timeout}
	var bodyReader io.Reader
	if len(request.Body) > 0 {
		bodyReader = bytes.NewReader(request.Body)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, request.Method, request.URL, bodyReader)
	if err != nil {
		return TransportResponse{}, err
	}
	for key, value := range request.Headers {
		httpRequest.Header.Set(key, value)
	}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return TransportResponse{}, err
	}
	defer httpResponse.Body.Close()
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return TransportResponse{}, err
	}
	return TransportResponse{
		Status:      httpResponse.StatusCode,
		ContentType: httpResponse.Header.Get("Content-Type"),
		Text:        string(responseBody),
	}, nil
}

func decodeResponse(response TransportResponse) (interface{}, error) {
	value := parseBody(response.Text)
	if response.Status < 200 || response.Status >= 300 {
		return nil, &APIError{Status: response.Status, Body: value}
	}
	return value, nil
}

func parseBody(text string) interface{} {
	if text == "" {
		return text
	}
	var value interface{}
	if err := json.Unmarshal([]byte(text), &value); err != nil {
		return text
	}
	return value
}

func stringifyQueryValue(value interface{}) string {
	switch typed := value.(type) {
	case bool:
		if typed {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprint(value)
	}
}

func truncateBody(value interface{}, limit int) string {
	if limit <= 0 {
		limit = 300
	}
	text := fmt.Sprint(value)
	if len(text) <= limit {
		return text
	}
	return text[:limit-1] + "…"
}
