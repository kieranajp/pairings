package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockReadCloser is a mock implementation of io.ReadCloser
type mockReadCloser struct {
	io.Reader
}

func (m *mockReadCloser) Close() error {
	return nil
}

// mockHTTPClient is a mock implementation of http.Client for testing
type mockHTTPClient struct {
	doFunc func(*http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestNewGeminiClient(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		model   string
		wantErr bool
	}{
		{
			name:    "valid configuration",
			apiKey:  "test-key",
			model:   "gemini-2.0-flash",
			wantErr: false,
		},
		{
			name:    "empty model uses default",
			apiKey:  "test-key",
			model:   "",
			wantErr: false,
		},
		{
			name:    "empty API key",
			apiKey:  "",
			model:   "gemini-2.0-flash",
			wantErr: false, // We don't validate API key in constructor
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewGeminiClient(tt.apiKey, tt.model)
			if client == nil {
				t.Error("NewGeminiClient returned nil")
			}
			if tt.model == "" && client.model != defaultModel {
				t.Errorf("NewGeminiClient() model = %v, want %v", client.model, defaultModel)
			}
		})
	}
}

func TestGeminiClient_Complete(t *testing.T) {
	tests := []struct {
		name           string
		prompt         string
		mockResponse   string
		mockStatusCode int
		mockErr        error
		wantResponse   string
		wantErr        bool
	}{
		{
			name:           "successful response",
			prompt:         "test prompt",
			mockResponse:   `{"candidates":[{"content":{"parts":[{"text":"test response"}]}}]}`,
			mockStatusCode: http.StatusOK,
			mockErr:        nil,
			wantResponse:   "test response",
			wantErr:        false,
		},
		{
			name:           "API error",
			prompt:         "test prompt",
			mockResponse:   "",
			mockStatusCode: http.StatusInternalServerError,
			mockErr:        nil,
			wantResponse:   "",
			wantErr:        true,
		},
		{
			name:           "HTTP client error",
			prompt:         "test prompt",
			mockResponse:   "",
			mockStatusCode: 0,
			mockErr:        errors.New("connection error"),
			wantResponse:   "",
			wantErr:        true,
		},
		{
			name:           "empty response",
			prompt:         "test prompt",
			mockResponse:   `{"candidates":[]}`,
			mockStatusCode: http.StatusOK,
			mockErr:        nil,
			wantResponse:   "",
			wantErr:        true,
		},
		{
			name:           "malformed response",
			prompt:         "test prompt",
			mockResponse:   `{"invalid": "response"}`,
			mockStatusCode: http.StatusOK,
			mockErr:        nil,
			wantResponse:   "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &GeminiClient{
				apiKey: "test-key",
				model:  "gemini-2.0-flash",
				client: &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						// Verify request
						if req.Method != "POST" {
							t.Errorf("expected POST request, got %s", req.Method)
						}
						if !strings.Contains(req.URL.String(), "gemini-2.0-flash") {
							t.Errorf("expected URL to contain model name, got %s", req.URL.String())
						}
						if !strings.Contains(req.URL.String(), "test-key") {
							t.Errorf("expected URL to contain API key")
						}

						// Return mock response
						if tt.mockErr != nil {
							return nil, tt.mockErr
						}
						return &http.Response{
							StatusCode: tt.mockStatusCode,
							Body:       &mockReadCloser{strings.NewReader(tt.mockResponse)},
						}, nil
					},
				},
			}

			got, err := client.Complete(context.Background(), tt.prompt)

			if (err != nil) != tt.wantErr {
				t.Errorf("GeminiClient.Complete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.wantResponse {
				t.Errorf("GeminiClient.Complete() = %v, want %v", got, tt.wantResponse)
			}
		})
	}
}

func TestGeminiClient_Complete_RequestValidation(t *testing.T) {
	client := &GeminiClient{
		apiKey: "test-key",
		model:  "gemini-2.0-flash",
		client: &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request body
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Errorf("failed to read request body: %v", err)
					return nil, err
				}

				// Verify request body contains prompt
				if !strings.Contains(string(body), "test prompt") {
					t.Errorf("request body does not contain prompt: %s", string(body))
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       &mockReadCloser{strings.NewReader(`{"candidates":[{"content":{"parts":[{"text":"test response"}]}}]}`)},
				}, nil
			},
		},
	}

	_, err := client.Complete(context.Background(), "test prompt")
	if err != nil {
		t.Errorf("GeminiClient.Complete() unexpected error: %v", err)
	}
}
