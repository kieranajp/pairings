package wine

import (
	"context"
	"os"
	"testing"

	"github.com/kieranajp/pairings/internal/domain/recipe"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockLLMClient is a mock implementation of client.LLMClient
type mockLLMClient struct {
	mock.Mock
}

func (m *mockLLMClient) Complete(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

// mockPromptGenerator is a mock implementation of prompt.Generator
type mockPromptGenerator struct {
	mock.Mock
}

func (m *mockPromptGenerator) GenerateWineRecommendationPrompt(
	dish, budgetMin, budgetMax, currency,
	styleStr, preferencesStr, occasionStr string,
) (string, error) {
	args := m.Called(dish, budgetMin, budgetMax, currency, styleStr, preferencesStr, occasionStr)
	return args.String(0), args.Error(1)
}

func (m *mockPromptGenerator) GenerateWinePairingPrompt(r *recipe.Recipe) (string, error) {
	args := m.Called(r)
	return args.String(0), args.Error(1)
}

// mockLogger is a mock implementation of logger.Logger
type mockLogger struct {
	mock.Mock
	logger zerolog.Logger
}

func newMockLogger() *mockLogger {
	return &mockLogger{
		logger: zerolog.New(os.Stdout),
	}
}

func (m *mockLogger) Info() *zerolog.Event {
	args := m.Called()
	if args.Get(0) == nil {
		return m.logger.Info()
	}
	return args.Get(0).(*zerolog.Event)
}

func (m *mockLogger) Debug() *zerolog.Event {
	args := m.Called()
	if args.Get(0) == nil {
		return m.logger.Debug()
	}
	return args.Get(0).(*zerolog.Event)
}

func (m *mockLogger) Error() *zerolog.Event {
	args := m.Called()
	if args.Get(0) == nil {
		return m.logger.Error()
	}
	return args.Get(0).(*zerolog.Event)
}

func TestService_GetRecommendations(t *testing.T) {
	tests := []struct {
		name                string
		dish                string
		budgetMin           int64
		budgetMax           int64
		currency            string
		wineType            string
		body                string
		tastePreferences    []string
		occasion            string
		mockPrompt          string
		mockResponse        string
		mockPromptErr       error
		mockResponseErr     error
		wantRecommendations string
		wantErr             bool
	}{
		{
			name:                "successful recommendation",
			dish:                "steak",
			budgetMin:           20,
			budgetMax:           50,
			currency:            "USD",
			wineType:            "red",
			body:                "full",
			tastePreferences:    []string{"bold", "dry"},
			occasion:            "dinner",
			mockPrompt:          "test prompt",
			mockResponse:        "test recommendations",
			mockPromptErr:       nil,
			mockResponseErr:     nil,
			wantRecommendations: "test recommendations",
			wantErr:             false,
		},
		{
			name:                "prompt generation error",
			dish:                "steak",
			budgetMin:           20,
			budgetMax:           50,
			currency:            "USD",
			wineType:            "red",
			body:                "full",
			tastePreferences:    []string{"bold", "dry"},
			occasion:            "dinner",
			mockPrompt:          "",
			mockResponse:        "",
			mockPromptErr:       assert.AnError,
			mockResponseErr:     nil,
			wantRecommendations: "",
			wantErr:             true,
		},
		{
			name:                "LLM error",
			dish:                "steak",
			budgetMin:           20,
			budgetMax:           50,
			currency:            "USD",
			wineType:            "red",
			body:                "full",
			tastePreferences:    []string{"bold", "dry"},
			occasion:            "dinner",
			mockPrompt:          "test prompt",
			mockResponse:        "",
			mockPromptErr:       nil,
			mockResponseErr:     assert.AnError,
			wantRecommendations: "",
			wantErr:             true,
		},
		{
			name:                "no preferences",
			dish:                "steak",
			budgetMin:           20,
			budgetMax:           50,
			currency:            "USD",
			wineType:            "",
			body:                "",
			tastePreferences:    nil,
			occasion:            "",
			mockPrompt:          "test prompt",
			mockResponse:        "test recommendations",
			mockPromptErr:       nil,
			mockResponseErr:     nil,
			wantRecommendations: "test recommendations",
			wantErr:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			llm := new(mockLLMClient)
			promptGen := new(mockPromptGenerator)
			log := newMockLogger()

			// Create profile for formatting
			profile := &PreferenceProfile{
				Dish:             tt.dish,
				Budget:           *NewBudget(tt.budgetMin, tt.budgetMax, tt.currency),
				PreferredStyle:   &WineStyle{Type: WineType(tt.wineType), Body: BodyType(tt.body)},
				TastePreferences: tt.tastePreferences,
				Occasion:         tt.occasion,
			}

			// Set up mock expectations
			promptGen.On("GenerateWineRecommendationPrompt",
				profile.Dish,
				profile.Budget.Min.Display(),
				profile.Budget.Max.Display(),
				profile.Budget.Currency,
				profile.FormatStyle(),
				profile.FormatPreferences(),
				profile.FormatOccasion(),
			).Return(tt.mockPrompt, tt.mockPromptErr)

			if tt.mockPromptErr == nil {
				llm.On("Complete", mock.Anything, tt.mockPrompt).Return(tt.mockResponse, tt.mockResponseErr)
			}
			log.On("Info").Return(nil)
			log.On("Debug").Return(nil)
			log.On("Error").Return(nil)

			// Create service
			service := NewService(llm, promptGen, log)

			// Call service
			got, err := service.GetRecommendations(
				context.Background(),
				tt.dish,
				tt.budgetMin,
				tt.budgetMax,
				tt.currency,
				tt.wineType,
				tt.body,
				tt.tastePreferences,
				tt.occasion,
			)

			// Assert results
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRecommendations, got)
			}

			// Verify mock expectations
			promptGen.AssertExpectations(t)
			llm.AssertExpectations(t)
		})
	}
}
