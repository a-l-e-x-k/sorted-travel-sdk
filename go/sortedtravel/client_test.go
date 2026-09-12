package sortedtravel

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestStatusUsesPublicPath(t *testing.T) {
	var captured TransportRequest
	client := NewClient(
		WithEnv(map[string]string{}),
		WithTransport(func(ctx context.Context, request TransportRequest, timeout time.Duration) (TransportResponse, error) {
			captured = request
			return TransportResponse{
				Status:      200,
				ContentType: "application/json",
				Text:        `{"status":"OK"}`,
			}, nil
		}),
	)
	payload, err := client.Status(context.Background())
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	expectedPayload := map[string]interface{}{"status": "OK"}
	if !reflect.DeepEqual(payload, expectedPayload) {
		t.Fatalf("unexpected payload: %#v", payload)
	}
	if captured.URL != "https://sorted.travel/api/v1/status" {
		t.Fatalf("unexpected url: %s", captured.URL)
	}
	if captured.Method != "GET" {
		t.Fatalf("unexpected method: %s", captured.Method)
	}
	if _, present := captured.Headers["Authorization"]; present {
		t.Fatalf("unexpected authorization header: %s", captured.Headers["Authorization"])
	}
}

func TestListDestinationsEncodesCursor(t *testing.T) {
	var capturedURL string
	client := NewClient(
		WithEnv(map[string]string{}),
		WithTransport(func(ctx context.Context, request TransportRequest, timeout time.Duration) (TransportResponse, error) {
			capturedURL = request.URL
			return TransportResponse{
				Status:      200,
				ContentType: "application/json",
				Text:        `{"data":[],"has_more":false}`,
			}, nil
		}),
	)
	cursor := "abc+1"
	limit := 5
	_, err := client.ListDestinations(context.Background(), &cursor, &limit)
	if err != nil {
		t.Fatalf("ListDestinations returned error: %v", err)
	}
	if !strings.Contains(capturedURL, "cursor=abc%2B1") {
		t.Fatalf("cursor not encoded in url: %s", capturedURL)
	}
	if !strings.Contains(capturedURL, "limit=5") {
		t.Fatalf("limit missing from url: %s", capturedURL)
	}
}

func TestAPIKeySetsBearerHeader(t *testing.T) {
	var capturedHeaders map[string]string
	client := NewClient(
		WithAPIKey("st_sandbox_test"),
		WithEnv(map[string]string{}),
		WithTransport(func(ctx context.Context, request TransportRequest, timeout time.Duration) (TransportResponse, error) {
			capturedHeaders = request.Headers
			return TransportResponse{
				Status:      200,
				ContentType: "application/json",
				Text:        `{}`,
			}, nil
		}),
	)
	_, err := client.Status(context.Background())
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if capturedHeaders["Authorization"] != "Bearer st_sandbox_test" {
		t.Fatalf("unexpected authorization header: %s", capturedHeaders["Authorization"])
	}
}

func TestNon2xxRaisesAPIError(t *testing.T) {
	client := NewClient(
		WithEnv(map[string]string{}),
		WithTransport(func(ctx context.Context, request TransportRequest, timeout time.Duration) (TransportResponse, error) {
			return TransportResponse{
				Status:      429,
				ContentType: "application/json",
				Text:        `{"error":"rate_limited"}`,
			}, nil
		}),
	)
	_, err := client.Status(context.Background())
	apiError, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %#v", err)
	}
	if apiError.Status != 429 {
		t.Fatalf("unexpected status: %d", apiError.Status)
	}
	expectedBody := map[string]interface{}{"error": "rate_limited"}
	if !reflect.DeepEqual(apiError.Body, expectedBody) {
		t.Fatalf("unexpected body: %#v", apiError.Body)
	}
}
