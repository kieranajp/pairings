package wine

import (
	"github.com/Rhymond/go-money"
)

// Level represents a qualitative measure (Low, Medium, High)
type Level string

const (
	LowLevel    Level = "low"
	MediumLevel Level = "medium"
	HighLevel   Level = "high"
)

// WineType represents the type of wine
type WineType string

const (
	Red       WineType = "red"
	White     WineType = "white"
	Rose      WineType = "rose"
	Sparkling WineType = "sparkling"
)

// BodyType represents the body of the wine
type BodyType string

const (
	LightBody  BodyType = "light"
	MediumBody BodyType = "medium"
	FullBody   BodyType = "full"
)

// Budget represents a price range using go-money for precise monetary values
type Budget struct {
	Min      *money.Money
	Max      *money.Money
	Currency string // ISO 4217 currency code (e.g., "EUR", "USD")
}

// NewBudget creates a new Budget with the given min and max values in the specified currency
func NewBudget(min, max int64, currency string) *Budget {
	return &Budget{
		Min:      money.New(min, currency),
		Max:      money.New(max, currency),
		Currency: currency,
	}
}

// WineStyle represents the characteristics of a wine
type WineStyle struct {
	Type      WineType
	Body      BodyType
	Sweetness Level
	Acidity   Level
	Tannin    Level // Primarily for reds
}

// Wine represents a specific wine with its characteristics
type Wine struct {
	Name            string
	Grape           string
	Region          string
	Style           WineStyle
	TastingNotes    []string
	PriceRange      Budget
	AgeingPotential string // Optional
}

// WineRecommendation represents a wine suggestion with pairing information
type WineRecommendation struct {
	Wine               Wine
	PairingExplanation string
	ConfidenceScore    float64
	UpgradeSuggestion  *Wine // Optional premium suggestion
}

// PreferenceProfile represents user preferences for wine recommendations
type PreferenceProfile struct {
	Dish             string
	Budget           Budget
	PreferredStyle   *WineStyle // Optional - nil means no preference
	TastePreferences []string   // "fruity", "dry", etc.
	Occasion         string     // Optional context
}
