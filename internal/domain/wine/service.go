package wine

import (
	"context"
	"fmt"

	"github.com/kieranajp/pairings/internal/infrastructure/client"
	"github.com/kieranajp/pairings/internal/infrastructure/logger"
	"github.com/kieranajp/pairings/internal/infrastructure/prompt"
)

// PromptGenerator is an interface for generating prompts
type PromptGenerator interface {
	GeneratePrompt(promptName string, args ...interface{}) (string, error)
}

// Service handles wine recommendation business logic
type Service struct {
	llm       client.LLMClient
	promptGen prompt.Generator
	log       logger.Logger
}

// NewService creates a new wine service
func NewService(
	llm client.LLMClient,
	promptGen prompt.Generator,
	log logger.Logger,
) *Service {
	return &Service{
		llm:       llm,
		promptGen: promptGen,
		log:       log,
	}
}

// GetRecommendations gets wine recommendations based on preferences
func (s *Service) GetRecommendations(
	ctx context.Context,
	dish string,
	budgetMin, budgetMax int64,
	currency string,
	wineType, body string,
	tastePreferences []string,
	occasion string,
) (string, error) {
	s.log.Info().Str("dish", dish).Msg("Getting wine recommendations")

	// Create budget
	budget := NewBudget(budgetMin, budgetMax, currency)

	// Create wine style if preferences are provided
	var style *WineStyle
	if wineType != "" || body != "" {
		style = &WineStyle{
			Type: WineType(wineType),
			Body: BodyType(body),
		}
	}

	// Create preference profile
	profile := &PreferenceProfile{
		Dish:             dish,
		Budget:           *budget,
		PreferredStyle:   style,
		TastePreferences: tastePreferences,
		Occasion:         occasion,
	}

	// Generate the prompt
	prompt, err := s.promptGen.GenerateWineRecommendationPrompt(
		profile.Dish,
		profile.Budget.Min.Display(),
		profile.Budget.Max.Display(),
		profile.Budget.Currency,
		profile.FormatStyle(),
		profile.FormatPreferences(),
		profile.FormatOccasion(),
	)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to generate prompt")
		return "", fmt.Errorf("failed to generate prompt: %w", err)
	}
	s.log.Debug().Str("prompt", prompt).Msg("Generated prompt")

	// Get recommendations from LLM
	recommendations, err := s.llm.Complete(ctx, prompt)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get recommendations")
		return "", fmt.Errorf("failed to get recommendations: %w", err)
	}

	return recommendations, nil
}

// formatStyle formats the wine style preferences for the prompt
func formatStyle(style *WineStyle) string {
	if style == nil {
		return "No specific style preferences"
	}
	return fmt.Sprintf("Preferred Style: %s, %s body", style.Type, style.Body)
}

// formatPreferences formats the taste preferences for the prompt
func formatPreferences(preferences []string) string {
	if len(preferences) == 0 {
		return "No specific taste preferences"
	}
	return fmt.Sprintf("Taste Preferences: %v", preferences)
}

// formatOccasion formats the occasion for the prompt
func formatOccasion(occasion string) string {
	if occasion == "" {
		return "No specific occasion"
	}
	return fmt.Sprintf("Occasion: %s", occasion)
}
