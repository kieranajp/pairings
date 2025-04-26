package prompt

import (
	"fmt"

	"github.com/kieranajp/pairings/internal/domain/recipe"
	"gopkg.in/yaml.v3"
)

// Generator is an interface for generating prompts
type Generator interface {
	GenerateWineRecommendationPrompt(
		dish, budgetMin, budgetMax, currency,
		styleStr, preferencesStr, occasionStr string,
	) (string, error)
	GenerateWinePairingPrompt(r *recipe.Recipe) (string, error)
}

type generator struct {
	schema  string
	prompts map[string]string
}

// NewGenerator creates a new prompt generator
func NewGenerator(schema string, promptsYAML string) (Generator, error) {
	var prompts map[string]string
	if err := yaml.Unmarshal([]byte(promptsYAML), &prompts); err != nil {
		return nil, fmt.Errorf("failed to parse prompts: %w", err)
	}

	return &generator{
		schema:  schema,
		prompts: prompts,
	}, nil
}

// generatePrompt is a private method for generating prompts from templates
func (g *generator) generatePrompt(promptName string, args ...interface{}) (string, error) {
	template, ok := g.prompts[promptName]
	if !ok {
		return "", fmt.Errorf("prompt template not found: %s", promptName)
	}

	// Add schema as the last argument
	args = append(args, g.schema)

	return fmt.Sprintf(template, args...), nil
}

// GenerateWineRecommendationPrompt generates a prompt for wine recommendations
func (g *generator) GenerateWineRecommendationPrompt(
	dish, budgetMin, budgetMax, currency,
	styleStr, preferencesStr, occasionStr string,
) (string, error) {
	return g.generatePrompt(
		"wine_preferences",
		dish,
		budgetMin,
		currency,
		budgetMax,
		currency,
		styleStr,
		preferencesStr,
		occasionStr,
	)
}

// GenerateWinePairingPrompt generates a prompt for wine pairing
func (g *generator) GenerateWinePairingPrompt(r *recipe.Recipe) (string, error) {
	return g.generatePrompt("wine_pairing", r.Title, r.Ingredients, r.Instructions, r.Cuisine)
}
