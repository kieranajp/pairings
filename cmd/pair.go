package cmd

import (
	recipeCLI "github.com/kieranajp/pairings/internal/application/cli"
	"github.com/kieranajp/pairings/internal/domain/recipe"
	"github.com/kieranajp/pairings/internal/infrastructure/client"
	"github.com/kieranajp/pairings/internal/infrastructure/logger"
	"github.com/kieranajp/pairings/internal/infrastructure/prompt"
	"github.com/urfave/cli/v2"
)

// PairCommand implements the Command interface for wine pairings
type PairCommand struct {
	llm           client.LLMClient
	recipeService *recipe.Service
	promptGen     prompt.Generator
	log           logger.Logger
}

// NewPairCommand creates a new pair command
func NewPairCommand() *PairCommand {
	return &PairCommand{}
}

func (c *PairCommand) WithLLMClient(llm client.LLMClient) *PairCommand {
	c.llm = llm
	return c
}

func (c *PairCommand) WithRecipeService(recipeService *recipe.Service) *PairCommand {
	c.recipeService = recipeService
	return c
}

func (c *PairCommand) WithPromptGen(promptGen prompt.Generator) *PairCommand {
	c.promptGen = promptGen
	return c
}

func (c *PairCommand) WithLog(log logger.Logger) *PairCommand {
	c.log = log
	return c
}

// Name returns the name of the command
func (c *PairCommand) Name() string {
	return "pair"
}

// Usage returns the usage description of the command
func (c *PairCommand) Usage() string {
	return "Get wine pairings for a recipe URL"
}

// Flags returns the command's flags
func (c *PairCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "recipe",
			Usage:    "Recipe URL",
			Required: true,
		},
	}
}

// Action returns a function that will be executed when the command is run
func (c *PairCommand) Action(ctx *cli.Context) error {
	handler := recipeCLI.NewRecipeHandler(
		c.llm,
		c.recipeService,
		c.promptGen,
		c.log,
	)
	return handler.Handle(ctx.Context, ctx.String("recipe"))
}
