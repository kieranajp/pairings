package validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// JSONValidator handles JSON validation and sanitization
type JSONValidator struct {
	schema string
}

// NewJSONValidator creates a new JSON validator with the given schema
func NewJSONValidator(schema string) *JSONValidator {
	return &JSONValidator{
		schema: schema,
	}
}

// ValidateAndSanitize extracts JSON from potentially markdown-wrapped text and validates it
func (v *JSONValidator) ValidateAndSanitize(input string) (string, error) {
	// First extract JSON from any markdown/text wrapping
	jsonStr, err := v.extractJSON(input)
	if err != nil {
		return "", fmt.Errorf("failed to extract JSON: %w", err)
	}

	// Then validate against schema
	if err := v.validate(jsonStr); err != nil {
		return "", fmt.Errorf("schema validation failed: %w", err)
	}

	return jsonStr, nil
}

// extractJSON finds and extracts the first valid JSON object or array from a string
func (v *JSONValidator) extractJSON(input string) (string, error) {
	// Find the first { or [ and the last } or ]
	start := -1
	end := -1

	// Look for object start/end
	objStart := strings.Index(input, "{")
	objEnd := strings.LastIndex(input, "}")

	// Look for array start/end
	arrStart := strings.Index(input, "[")
	arrEnd := strings.LastIndex(input, "]")

	// Determine which comes first and is a complete pair
	if objStart != -1 && objEnd != -1 && (arrStart == -1 || objStart < arrStart) {
		start = objStart
		end = objEnd
	} else if arrStart != -1 && arrEnd != -1 {
		start = arrStart
		end = arrEnd
	}

	if start == -1 || end == -1 || start > end {
		return "", fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := input[start : end+1]

	// Validate that it's actually valid JSON
	var js json.RawMessage
	if err := json.Unmarshal([]byte(jsonStr), &js); err != nil {
		return "", fmt.Errorf("extracted content is not valid JSON: %w", err)
	}

	return jsonStr, nil
}

// validate validates the JSON string against the schema
func (v *JSONValidator) validate(jsonStr string) error {
	schemaLoader := gojsonschema.NewStringLoader(v.schema)
	documentLoader := gojsonschema.NewStringLoader(jsonStr)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	if !result.Valid() {
		var errors []string
		for _, err := range result.Errors() {
			errors = append(errors, err.String())
		}
		return fmt.Errorf("invalid JSON response: %s", strings.Join(errors, "; "))
	}

	return nil
}
