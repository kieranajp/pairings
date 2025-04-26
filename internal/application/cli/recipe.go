package cli

import (
	"context"
	"fmt"

	"github.com/kieranajp/pairings/internal/domain/recipe"
	"github.com/kieranajp/pairings/internal/infrastructure/client"
	"github.com/kieranajp/pairings/internal/infrastructure/logger"
	"github.com/kieranajp/pairings/internal/infrastructure/prompt"
)

type RecipeHandler struct {
	llm           client.LLMClient
	recipeService *recipe.Service
	promptGen     prompt.Generator
	logger        logger.Logger
}

func NewRecipeHandler(
	llm client.LLMClient,
	recipeService *recipe.Service,
	promptGen prompt.Generator,
	logger logger.Logger,
) *RecipeHandler {
	return &RecipeHandler{
		llm:           llm,
		recipeService: recipeService,
		promptGen:     promptGen,
		logger:        logger,
	}
}

func (h *RecipeHandler) Handle(ctx context.Context, url string) error {
	h.logger.Info().Str("url", url).Msg("Getting wine pairings")

	// Get recipe details
	r, err := h.recipeService.GetRecipe(ctx, url)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get recipe")
		return fmt.Errorf("failed to get recipe: %w", err)
	}
	h.logger.Info().Str("title", r.Title).Msg("Got recipe details")

	// Generate prompt
	prompt, err := h.promptGen.GenerateWinePairingPrompt(r)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to generate prompt")
		return fmt.Errorf("failed to generate prompt: %w", err)
	}
	h.logger.Debug().Str("prompt", prompt).Msg("Generated prompt")

	// Get wine pairings from LLM
	pairings, err := h.llm.Complete(ctx, prompt)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get pairings")
		return fmt.Errorf("failed to get pairings: %w", err)
	}

	// Display results
	fmt.Println("Wine Pairings for:", r.Title)
	fmt.Println(pairings)

	return nil
}
