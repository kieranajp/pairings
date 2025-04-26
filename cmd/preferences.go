package cmd

import (
	appCLI "github.com/kieranajp/pairings/internal/application/cli"
	"github.com/kieranajp/pairings/internal/infrastructure/client"
	"github.com/kieranajp/pairings/internal/infrastructure/logger"
	promptGenerator "github.com/kieranajp/pairings/internal/infrastructure/prompt"
	"github.com/urfave/cli/v2"
)

// PreferencesCommand implements the Command interface for wine preferences
type PreferencesCommand struct {
	llm       client.LLMClient
	promptGen *promptGenerator.Generator
	log       logger.Logger
	schema    string
}

// NewPreferencesCommand creates a new preferences command
func NewPreferencesCommand() *PreferencesCommand {
	return &PreferencesCommand{}
}

func (c *PreferencesCommand) WithLLMClient(llm client.LLMClient) *PreferencesCommand {
	c.llm = llm
	return c
}

func (c *PreferencesCommand) WithPromptGen(promptGen *promptGenerator.Generator) *PreferencesCommand {
	c.promptGen = promptGen
	return c
}

func (c *PreferencesCommand) WithLog(log logger.Logger) *PreferencesCommand {
	c.log = log
	return c
}

func (c *PreferencesCommand) WithSchema(schema string) *PreferencesCommand {
	c.schema = schema
	return c
}

// Name returns the name of the command
func (c *PreferencesCommand) Name() string {
	return "preferences"
}

// Usage returns the usage description of the command
func (c *PreferencesCommand) Usage() string {
	return "Create a wine preference profile"
}

// Flags returns the command's flags
func (c *PreferencesCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "dish",
			Usage:    "Name of the dish to pair with",
			Required: true,
		},
		&cli.Int64Flag{
			Name:     "budget-min",
			Usage:    "Minimum budget in cents (e.g., 2000 for 20.00)",
			Required: true,
		},
		&cli.Int64Flag{
			Name:     "budget-max",
			Usage:    "Maximum budget in cents (e.g., 5000 for 50.00)",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "currency",
			Usage: "Currency code (e.g., EUR, USD)",
			Value: "EUR",
		},
		&cli.StringFlag{
			Name:  "wine-type",
			Usage: "Preferred wine type (red, white, rose, sparkling) (optional)",
		},
		&cli.StringFlag{
			Name:  "body",
			Usage: "Preferred wine body (light, medium, full) (optional)",
		},
		&cli.StringSliceFlag{
			Name:  "taste-preferences",
			Usage: "Taste preferences (e.g., fruity, dry, oaky) (optional)",
		},
		&cli.StringFlag{
			Name:  "occasion",
			Usage: "Occasion context (e.g., dinner party, casual meal) (optional)",
		},
	}
}

// Action returns a function that will be executed when the command is run
func (c *PreferencesCommand) Action(ctx *cli.Context) error {
	handler := appCLI.NewPreferencesHandler(
		c.llm,
		c.promptGen,
		c.log,
	)
	return handler.Handle(
		ctx.Context,
		ctx.String("dish"),
		ctx.Int64("budget-min"),
		ctx.Int64("budget-max"),
		ctx.String("currency"),
		ctx.String("wine-type"),
		ctx.String("body"),
		ctx.StringSlice("taste-preferences"),
		ctx.String("occasion"),
	)
}
