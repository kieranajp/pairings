package cli

import (
	"context"
	"fmt"

	"github.com/kieranajp/pairings/internal/domain/wine"
)

// PreferencesHandler handles the preferences command
type PreferencesHandler struct {
	service *wine.Service
}

// NewPreferencesHandler creates a new preferences handler
func NewPreferencesHandler(service *wine.Service) *PreferencesHandler {
	return &PreferencesHandler{
		service: service,
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
	// Get recommendations from service
	recommendations, err := h.service.GetRecommendations(
		ctx,
		dish,
		budgetMin,
		budgetMax,
		currency,
		wineType,
		body,
		tastePreferences,
		occasion,
	)
	if err != nil {
		return err
	}

	// Display results
	fmt.Println("Wine Recommendations for:", dish)
	fmt.Println(recommendations)

	return nil
}
