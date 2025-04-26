package client

import "context"

// LLMClient defines the interface for language model clients
type LLMClient interface {
	// Complete sends a prompt to the LLM and returns its response
	Complete(ctx context.Context, prompt string) (string, error)
}
