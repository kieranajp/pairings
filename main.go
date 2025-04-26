package main

import (
	_ "embed"
	"fmt"
	"os"

	recipeCLI "github.com/kieranajp/pairings/internal/application/cli"
	"github.com/kieranajp/pairings/internal/domain/recipe"
	"github.com/kieranajp/pairings/internal/infrastructure/client"
	"github.com/kieranajp/pairings/internal/infrastructure/logger"
	promptGenerator "github.com/kieranajp/pairings/internal/infrastructure/prompt"
	"github.com/urfave/cli/v2"
)

//go:embed config/schema.json
var schema string

//go:embed config/prompts.yaml
var prompts string

var (
	geminiClient  *client.GeminiClient
	recipeService *recipe.Service
	promptGen     *promptGenerator.Generator
	log           logger.Logger
	app           *cli.App
)

func setup(c *cli.Context) error {
	log = logger.New(c.String("log-level"))
	geminiClient = client.NewGeminiClient(
		c.String("gemini-api-key"),
		c.String("gemini-model"),
	)

	recipeService = recipe.NewService()

	var err error
	promptGen, err = promptGenerator.NewGenerator(schema, prompts)
	if err != nil {
		return fmt.Errorf("failed to initialize prompt generator: %w", err)
	}

	return nil
}
func main() {
	app = &cli.App{
		Name:  "pairings",
		Usage: "Find wine pairings for recipes",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "gemini-api-key",
				Usage:    "Gemini API key",
				EnvVars:  []string{"GEMINI_API_KEY"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "gemini-model",
				Usage:   "Gemini model to use",
				EnvVars: []string{"GEMINI_MODEL"},
				Value:   "gemini-2.0-flash",
			},
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "Log level (debug, info, warn, error)",
				EnvVars: []string{"LOG_LEVEL"},
				Value:   "info",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "pair",
				Usage: "Get wine pairings for a recipe URL",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "recipe",
						Usage:    "Recipe URL",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					handler := recipeCLI.NewRecipeHandler(
						geminiClient,
						recipeService,
						promptGen,
						log,
					)
					return handler.Handle(c.Context, c.String("recipe"))
				},
			},
		},
		Before: setup,
	}

	if err := app.Run(os.Args); err != nil {
		log.Error().Err(err).Msg("Application failed")
		os.Exit(1)
	}
}
