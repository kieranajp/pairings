package wine

import "fmt"

// PreferenceProfile represents a user's wine preferences
type PreferenceProfile struct {
	Dish             string
	Budget           Budget
	PreferredStyle   *WineStyle
	TastePreferences []string
	Occasion         string
}

// FormatStyle formats the wine style preferences for the prompt
func (p *PreferenceProfile) FormatStyle() string {
	if p.PreferredStyle == nil || (p.PreferredStyle.Type == "" && p.PreferredStyle.Body == "") {
		return "No specific style preferences"
	}
	if p.PreferredStyle.Type == "" {
		return fmt.Sprintf("Preferred Style: %s body", p.PreferredStyle.Body)
	}
	if p.PreferredStyle.Body == "" {
		return fmt.Sprintf("Preferred Style: %s", p.PreferredStyle.Type)
	}
	return fmt.Sprintf("Preferred Style: %s, %s body", p.PreferredStyle.Type, p.PreferredStyle.Body)
}

// FormatPreferences formats the taste preferences for the prompt
func (p *PreferenceProfile) FormatPreferences() string {
	if len(p.TastePreferences) == 0 {
		return "No specific taste preferences"
	}
	return fmt.Sprintf("Taste Preferences: %v", p.TastePreferences)
}

// FormatOccasion formats the occasion for the prompt
func (p *PreferenceProfile) FormatOccasion() string {
	if p.Occasion == "" {
		return "No specific occasion"
	}
	return fmt.Sprintf("Occasion: %s", p.Occasion)
}
