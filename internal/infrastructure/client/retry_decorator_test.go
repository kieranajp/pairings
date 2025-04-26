package client

import (
	"context"
	"errors"
	"testing"
	"time"
)

// mockLLMClient is a mock implementation of LLMClient for testing
type mockLLMClient struct {
	responses []string
	errors    []error
	callCount int
}

func (m *mockLLMClient) Complete(ctx context.Context, prompt string) (string, error) {
	if m.callCount >= len(m.responses) {
		return "", errors.New("mock client: no more responses")
	}
	response := m.responses[m.callCount]
	err := m.errors[m.callCount]
	m.callCount++
	return response, err
}

func TestRetryDecorator(t *testing.T) {
	tests := []struct {
		name           string
		responses      []string
		errors         []error
		maxRetries     int
		initialBackoff time.Duration
		maxBackoff     time.Duration
		wantResponse   string
		wantErr        bool
	}{
		{
			name:           "success on first try",
			responses:      []string{"success"},
			errors:         []error{nil},
			maxRetries:     3,
			initialBackoff: 10 * time.Millisecond,
			maxBackoff:     100 * time.Millisecond,
			wantResponse:   "success",
			wantErr:        false,
		},
		{
			name:           "success after retry",
			responses:      []string{"", "success"},
			errors:         []error{errors.New("temporary error"), nil},
			maxRetries:     3,
			initialBackoff: 10 * time.Millisecond,
			maxBackoff:     100 * time.Millisecond,
			wantResponse:   "success",
			wantErr:        false,
		},
		{
			name:           "max retries exceeded",
			responses:      []string{"", "", "", ""},
			errors:         []error{errors.New("error 1"), errors.New("error 2"), errors.New("error 3"), errors.New("error 4")},
			maxRetries:     3,
			initialBackoff: 10 * time.Millisecond,
			maxBackoff:     100 * time.Millisecond,
			wantResponse:   "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockLLMClient{
				responses: tt.responses,
				errors:    tt.errors,
			}

			client := NewRetryDecorator(mock, tt.maxRetries, tt.initialBackoff, tt.maxBackoff)
			ctx := context.Background()

			got, err := client.Complete(ctx, "test prompt")

			if (err != nil) != tt.wantErr {
				t.Errorf("RetryDecorator.Complete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.wantResponse {
				t.Errorf("RetryDecorator.Complete() = %v, want %v", got, tt.wantResponse)
			}
		})
	}
}

func TestRetryDecoratorContextCancellation(t *testing.T) {
	mock := &mockLLMClient{
		responses: []string{"", "", "success"},
		errors:    []error{errors.New("error 1"), errors.New("error 2"), nil},
	}

	client := NewRetryDecorator(mock, 3, 100*time.Millisecond, 1*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.Complete(ctx, "test prompt")
	if err == nil {
		t.Error("RetryDecorator.Complete() expected error due to context cancellation")
	}
}
