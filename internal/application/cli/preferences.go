package cli

import (
	"context"
	"fmt"

	"github.com/kieranajp/pairings/internal/domain/wine"
	"github.com/kieranajp/pairings/internal/infrastructure/client"
	"github.com/kieranajp/pairings/internal/infrastructure/logger"
	promptGenerator "github.com/kieranajp/pairings/internal/infrastructure/prompt"
	"github.com/xeipuuv/gojsonschema"
)

// PreferencesHandler handles the preferences command
type PreferencesHandler struct {
	geminiClient *client.GeminiClient
	promptGen    *promptGenerator.Generator
	log          logger.Logger
	schema       string
}

// NewPreferencesHandler creates a new preferences handler
func NewPreferencesHandler(
	geminiClient *client.GeminiClient,
	promptGen *promptGenerator.Generator,
	log logger.Logger,
	schema string,
) *PreferencesHandler {
	return &PreferencesHandler{
		geminiClient: geminiClient,
		promptGen:    promptGen,
		log:          log,
		schema:       schema,
	}
}

// Handle processes the preferences command
func (h *PreferencesHandler) Handle(
	ctx context.Context,
	dish string,
	budgetMin, budgetMax int64,
	currency string,
	wineType, body string,
	tastePreferences []string,
	occasion string,
) error {
	h.log.Info().Str("dish", dish).Msg("Getting wine recommendations")

	// Create budget
	budget := wine.NewBudget(budgetMin, budgetMax, currency)

	// Create wine style if preferences are provided
	var style *wine.WineStyle
	if wineType != "" || body != "" {
		style = &wine.WineStyle{
			Type: wine.WineType(wineType),
			Body: wine.BodyType(body),
		}
	}

	// Create preference profile
	profile := &wine.PreferenceProfile{
		Dish:             dish,
		Budget:           *budget,
		PreferredStyle:   style,
		TastePreferences: tastePreferences,
		Occasion:         occasion,
	}

	// Generate the prompt
	prompt, err := h.promptGen.GeneratePrompt(
		"wine_preferences",
		profile.Dish,
		profile.Budget.Min.Display(),
		profile.Budget.Currency,
		profile.Budget.Max.Display(),
		profile.Budget.Currency,
		formatStyle(profile.PreferredStyle),
		formatPreferences(profile.TastePreferences),
		formatOccasion(profile.Occasion),
	)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to generate prompt")
		return fmt.Errorf("failed to generate prompt: %w", err)
	}
	h.log.Debug().Str("prompt", prompt).Msg("Generated prompt")

	// Get wine recommendations
	rawResponse, err := h.geminiClient.GetPairings(ctx, prompt)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get recommendations")
		return fmt.Errorf("failed to get recommendations: %w", err)
	}

	// Validate against schema
	schemaLoader := gojsonschema.NewStringLoader(h.schema)
	documentLoader := gojsonschema.NewStringLoader(rawResponse)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to validate response")
		return fmt.Errorf("failed to validate response: %w", err)
	}

	if !result.Valid() {
		h.log.Error().Interface("errors", result.Errors()).Msg("Invalid response")
		return fmt.Errorf("invalid response: %v", result.Errors())
	}

	// Display results
	fmt.Println("Wine Recommendations for:", profile.Dish)
	fmt.Println(rawResponse)

	return nil
}

// formatStyle formats the wine style preferences for the prompt
func formatStyle(style *wine.WineStyle) string {
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
