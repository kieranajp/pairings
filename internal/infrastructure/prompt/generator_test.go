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

func TestGenerateWineRecommendationPrompt(t *testing.T) {
	gen, err := NewGenerator(
		`{"type": "array"}`,
		`wine_preferences: "Dish: %s\nBudget: %s %s - %s %s\n%s\n%s\n%s\nSchema: %s"`,
	)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		dish           string
		budgetMin      string
		budgetMax      string
		currency       string
		styleStr       string
		preferencesStr string
		occasionStr    string
		wantErr        bool
		errContains    string
		expected       string
	}{
		{
			name:           "valid prompt",
			dish:           "steak",
			budgetMin:      "20",
			budgetMax:      "50",
			currency:       "USD",
			styleStr:       "Preferred Style: red, full body",
			preferencesStr: "Taste Preferences: [bold dry]",
			occasionStr:    "Occasion: dinner",
			wantErr:        false,
			expected:       "Dish: steak\nBudget: 20 USD - 50 USD\nPreferred Style: red, full body\nTaste Preferences: [bold dry]\nOccasion: dinner\nSchema: {\"type\": \"array\"}",
		},
		{
			name:           "missing prompt",
			dish:           "steak",
			budgetMin:      "20",
			budgetMax:      "50",
			currency:       "USD",
			styleStr:       "No specific style preferences",
			preferencesStr: "No specific taste preferences",
			occasionStr:    "No specific occasion",
			wantErr:        false,
			expected:       "Dish: steak\nBudget: 20 USD - 50 USD\nNo specific style preferences\nNo specific taste preferences\nNo specific occasion\nSchema: {\"type\": \"array\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := gen.GenerateWineRecommendationPrompt(
				tt.dish,
				tt.budgetMin,
				tt.budgetMax,
				tt.currency,
				tt.styleStr,
				tt.preferencesStr,
				tt.occasionStr,
			)
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
