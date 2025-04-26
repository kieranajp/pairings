package prompt

import (
	"fmt"

	"github.com/kieranajp/pairings/internal/domain/recipe"
	"gopkg.in/yaml.v3"
)

type Generator struct {
	schema  string
	prompts map[string]string
}

func NewGenerator(schema string, promptsYAML string) (*Generator, error) {
	var prompts map[string]string
	if err := yaml.Unmarshal([]byte(promptsYAML), &prompts); err != nil {
		return nil, fmt.Errorf("failed to parse prompts: %w", err)
	}

	if _, ok := prompts["wine_pairing"]; !ok {
		return nil, fmt.Errorf("missing required prompt template: wine_pairing")
	}

	return &Generator{
		schema:  schema,
		prompts: prompts,
	}, nil
}

func (g *Generator) GenerateWinePairingPrompt(r *recipe.Recipe) string {
	return fmt.Sprintf(g.prompts["wine_pairing"], r.Title, r.Ingredients, r.Instructions, r.Cuisine, g.schema)
}
