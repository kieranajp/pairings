package client

import (
	"context"
	"errors"
	"testing"
)

// mockValidatorClient is a mock implementation of LLMClient for testing
type mockValidatorClient struct {
	response string
	err      error
}

func (m *mockValidatorClient) Complete(ctx context.Context, prompt string) (string, error) {
	return m.response, m.err
}

func TestValidatorDecorator(t *testing.T) {
	// Define a simple JSON schema for testing
	schema := `{
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"age": { "type": "number" }
		},
		"required": ["name", "age"]
	}`

	tests := []struct {
		name         string
		response     string
		err          error
		wantResponse string
		wantErr      bool
	}{
		{
			name:         "valid JSON response",
			response:     `{"name": "John", "age": 30}`,
			err:          nil,
			wantResponse: `{"name": "John", "age": 30}`,
			wantErr:      false,
		},
		{
			name:         "invalid JSON response",
			response:     `{"name": "John", "age": "30"}`, // age should be a number
			err:          nil,
			wantResponse: "",
			wantErr:      true,
		},
		{
			name:         "missing required field",
			response:     `{"name": "John"}`, // missing age
			err:          nil,
			wantResponse: "",
			wantErr:      true,
		},
		{
			name:         "client error",
			response:     "",
			err:          errors.New("client error"),
			wantResponse: "",
			wantErr:      true,
		},
		{
			name:         "malformed JSON",
			response:     `{"name": "John", "age": 30`, // missing closing brace
			err:          nil,
			wantResponse: "",
			wantErr:      true,
		},
		{
			name:         "JSON in markdown code block",
			response:     "```json\n{\"name\": \"John\", \"age\": 30}\n```",
			err:          nil,
			wantResponse: `{"name": "John", "age": 30}`,
			wantErr:      false,
		},
		{
			name:         "JSON with explanatory text",
			response:     "Sure! Here's the JSON you asked for! {\"name\": \"John\", \"age\": 30}",
			err:          nil,
			wantResponse: `{"name": "John", "age": 30}`,
			wantErr:      false,
		},
		{
			name:         "JSON with markdown and explanatory text",
			response:     "Here's the JSON response:\n```json\n{\"name\": \"John\", \"age\": 30}\n```\nLet me know if you need anything else!",
			err:          nil,
			wantResponse: `{"name": "John", "age": 30}`,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockValidatorClient{
				response: tt.response,
				err:      tt.err,
			}

			client := NewValidatorDecorator(mock, schema)
			ctx := context.Background()

			got, err := client.Complete(ctx, "test prompt")

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatorDecorator.Complete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.wantResponse {
				t.Errorf("ValidatorDecorator.Complete() = %v, want %v", got, tt.wantResponse)
			}
		})
	}
}

func TestValidatorDecoratorWithComplexSchema(t *testing.T) {
	// Define a more complex JSON schema for testing
	schema := `{
		"type": "object",
		"properties": {
			"user": {
				"type": "object",
				"properties": {
					"name": { "type": "string" },
					"email": { "type": "string", "format": "email" },
					"roles": {
						"type": "array",
						"items": { "type": "string" }
					}
				},
				"required": ["name", "email", "roles"]
			}
		},
		"required": ["user"]
	}`

	tests := []struct {
		name         string
		response     string
		wantResponse string
		wantErr      bool
	}{
		{
			name:         "valid complex JSON",
			response:     `{"user": {"name": "John", "email": "john@example.com", "roles": ["admin", "user"]}}`,
			wantResponse: `{"user": {"name": "John", "email": "john@example.com", "roles": ["admin", "user"]}}`,
			wantErr:      false,
		},
		{
			name:         "invalid email format",
			response:     `{"user": {"name": "John", "email": "invalid-email", "roles": ["admin"]}}`,
			wantResponse: "",
			wantErr:      true,
		},
		{
			name:         "missing required array",
			response:     `{"user": {"name": "John", "email": "john@example.com"}}`,
			wantResponse: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockValidatorClient{
				response: tt.response,
				err:      nil,
			}

			client := NewValidatorDecorator(mock, schema)
			ctx := context.Background()

			got, err := client.Complete(ctx, "test prompt")

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatorDecorator.Complete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.wantResponse {
				t.Errorf("ValidatorDecorator.Complete() = %v, want %v", got, tt.wantResponse)
			}
		})
	}
}
