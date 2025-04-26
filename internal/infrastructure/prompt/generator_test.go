package prompt

import (
	"testing"

	"github.com/kieranajp/pairings/internal/domain/recipe"
	"github.com/stretchr/testify/assert"
)

func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name        string
		schema      string
		prompts     string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid schema and prompts",
			schema:  `{"type": "array"}`,
			prompts: `wine_pairing: "test prompt"`,
			wantErr: false,
		},
		{
			name:        "invalid YAML",
			schema:      `{"type": "array"}`,
			prompts:     `invalid: yaml: [`,
			wantErr:     true,
			errContains: "failed to parse prompts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := NewGenerator(tt.schema, tt.prompts)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, gen)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gen)
			}
		})
	}
}

func TestGeneratePrompt(t *testing.T) {
	gen, err := NewGenerator(
		`{"type": "array"}`,
		`wine_pairing: "Recipe: %s\nIngredients: %v\nMethod: %v\nCuisine: %s\nSchema: %s"`,
	)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		promptName  string
		args        []interface{}
		wantErr     bool
		errContains string
		expected    string
	}{
		{
			name:       "valid prompt",
			promptName: "wine_pairing",
			args: []interface{}{
				"Test Recipe",
				[]string{"ingredient1", "ingredient2"},
				[]string{"step1", "step2"},
				"Test Cuisine",
			},
			wantErr:  false,
			expected: "Recipe: Test Recipe\nIngredients: [ingredient1 ingredient2]\nMethod: [step1 step2]\nCuisine: Test Cuisine\nSchema: {\"type\": \"array\"}",
		},
		{
			name:        "missing prompt",
			promptName:  "nonexistent",
			args:        []interface{}{},
			wantErr:     true,
			errContains: "prompt template not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := gen.GeneratePrompt(tt.promptName, tt.args...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Empty(t, actual)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestGenerateWinePairingPrompt(t *testing.T) {
	gen, err := NewGenerator(
		`{"type": "array"}`,
		`wine_pairing: "Recipe: %s\nIngredients: %v\nMethod: %v\nCuisine: %s\nSchema: %s"`,
	)
	assert.NoError(t, err)

	recipe := &recipe.Recipe{
		Title:        "Test Recipe",
		Ingredients:  []string{"ingredient1", "ingredient2"},
		Instructions: []string{"step1", "step2"},
		Cuisine:      "Test Cuisine",
	}

	expected := "Recipe: Test Recipe\nIngredients: [ingredient1 ingredient2]\nMethod: [step1 step2]\nCuisine: Test Cuisine\nSchema: {\"type\": \"array\"}"
	actual, err := gen.GenerateWinePairingPrompt(recipe)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
