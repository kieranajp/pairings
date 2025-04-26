package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/kieranajp/pairings/cmd"
	"github.com/kieranajp/pairings/internal/domain/recipe"
	"github.com/kieranajp/pairings/internal/infrastructure/client"
	"github.com/kieranajp/pairings/internal/infrastructure/logger"
	promptGenerator "github.com/kieranajp/pairings/internal/infrastructure/prompt"
	"github.com/urfave/cli/v2"
)

//go:embed config/pairings_schema.json
var pairingsSchema string

//go:embed config/preferences_schema.json
var preferencesSchema string

//go:embed config/prompts.yaml
var prompts string

var (
	geminiClient      *client.GeminiClient
	recipeService     *recipe.Service
	pairingsPromptGen *promptGenerator.Generator
	prefsPromptGen    *promptGenerator.Generator
	log               logger.Logger
)

func setup(c *cli.Context) error {
	log = logger.New(c.String("log-level"))
	geminiClient = client.NewGeminiClient(
		c.String("gemini-api-key"),
		c.String("gemini-model"),
	)

	recipeService = recipe.NewService()

	var err error
	pairingsPromptGen, err = promptGenerator.NewGenerator(pairingsSchema, prompts)
	if err != nil {
		return fmt.Errorf("failed to initialize pairings prompt generator: %w", err)
	}

	prefsPromptGen, err = promptGenerator.NewGenerator(preferencesSchema, prompts)
	if err != nil {
		return fmt.Errorf("failed to initialize preferences prompt generator: %w", err)
	}

	return nil
}

func newApp() *cli.App {
	preferences := cmd.NewPreferencesCommand()
	pair := cmd.NewPairCommand()

	return &cli.App{
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
				Name:  preferences.Name(),
				Usage: preferences.Usage(),
				Flags: preferences.Flags(),
				Action: func(c *cli.Context) error {
					if err := setup(c); err != nil {
						return err
					}
					return preferences.
						WithGeminiClient(geminiClient).
						WithPromptGen(prefsPromptGen).
						WithLog(log).
						WithSchema(preferencesSchema).
						Action(c)
				},
			},
			{
				Name:  pair.Name(),
				Usage: pair.Usage(),
				Flags: pair.Flags(),
				Action: func(c *cli.Context) error {
					if err := setup(c); err != nil {
						return err
					}
					return pair.
						WithGeminiClient(geminiClient).
						WithRecipeService(recipeService).
						WithPromptGen(pairingsPromptGen).
						WithLog(log).
						Action(c)
				},
			},
		},
	}
}

func main() {
	app := newApp()

	if err := app.Run(os.Args); err != nil {
		log.Error().Err(err).Msg("Application failed")
		os.Exit(1)
	}
}
