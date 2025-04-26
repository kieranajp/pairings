package client

import (
	"context"
	"fmt"

	"github.com/kieranajp/pairings/internal/infrastructure/validator"
)

// ValidatorDecorator wraps an LLMClient and adds JSON validation
type ValidatorDecorator struct {
	client    LLMClient
	validator *validator.JSONValidator
}

// NewValidatorDecorator creates a new validator decorator
func NewValidatorDecorator(client LLMClient, schema string) *ValidatorDecorator {
	return &ValidatorDecorator{
		client:    client,
		validator: validator.NewJSONValidator(schema),
	}
}

// Complete wraps the underlying client's Complete method with JSON validation
func (d *ValidatorDecorator) Complete(ctx context.Context, prompt string) (string, error) {
	// Get response from underlying client
	response, err := d.client.Complete(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("client error: %w", err)
	}

	// Validate and sanitize the response
	validJSON, err := d.validator.ValidateAndSanitize(response)
	if err != nil {
		return "", fmt.Errorf("validation error: %w", err)
	}

	return validJSON, nil
}
