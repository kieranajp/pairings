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

	return &Generator{
		schema:  schema,
		prompts: prompts,
	}, nil
}

func (g *Generator) GeneratePrompt(promptName string, args ...interface{}) (string, error) {
	template, ok := g.prompts[promptName]
	if !ok {
		return "", fmt.Errorf("prompt template not found: %s", promptName)
	}

	// Add schema as the last argument
	args = append(args, g.schema)

	return fmt.Sprintf(template, args...), nil
}

func (g *Generator) GenerateWinePairingPrompt(r *recipe.Recipe) (string, error) {
	return g.GeneratePrompt("wine_pairing", r.Title, r.Ingredients, r.Instructions, r.Cuisine)
}
