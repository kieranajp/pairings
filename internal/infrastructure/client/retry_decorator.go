package client

import (
	"context"
	"fmt"
	"math"
	"time"
)

// RetryDecorator wraps an LLMClient and adds exponential backoff retry functionality
type RetryDecorator struct {
	client         LLMClient
	maxRetries     int
	initialBackoff time.Duration
	maxBackoff     time.Duration
}

// NewRetryDecorator creates a new retry decorator with the given configuration
func NewRetryDecorator(client LLMClient, maxRetries int, initialBackoff, maxBackoff time.Duration) *RetryDecorator {
	return &RetryDecorator{
		client:         client,
		maxRetries:     maxRetries,
		initialBackoff: initialBackoff,
		maxBackoff:     maxBackoff,
	}
}

// Complete implements the LLMClient interface with exponential backoff retry
func (d *RetryDecorator) Complete(ctx context.Context, prompt string) (string, error) {
	var lastErr error
	backoff := d.initialBackoff

	for attempt := 0; attempt <= d.maxRetries; attempt++ {
		// Try to get response from underlying client
		response, err := d.client.Complete(ctx, prompt)
		if err == nil {
			return response, nil
		}

		lastErr = err

		// If this was the last attempt, don't wait
		if attempt == d.maxRetries {
			break
		}

		// Calculate next backoff duration with exponential increase
		backoff = time.Duration(math.Min(float64(backoff*2), float64(d.maxBackoff)))

		// Create a timer for the backoff
		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return "", fmt.Errorf("context cancelled during retry: %w", ctx.Err())
		case <-timer.C:
			// Continue to next attempt
		}
	}

	return "", fmt.Errorf("failed after %d retries: %w", d.maxRetries, lastErr)
}
