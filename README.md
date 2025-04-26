# Pairings

[![Main](https://github.com/kieranajp/pairings/actions/workflows/main.yml/badge.svg)](https://github.com/kieranajp/pairings/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/kieranajp/pairings/branch/main/graph/badge.svg)](https://codecov.io/gh/kieranajp/pairings)

A CLI tool that suggests wine pairings for recipes. Given a recipe URL, it analyzes the dish and suggests three wines that would pair well, including detailed reasoning for each pairing.

## Features

- Recipe analysis from URLs
- Structured wine pairing suggestions
- Detailed reasoning for each pairing
- Configurable logging levels
- Support for different Gemini models
- Wine preference profile creation

## Installation

```bash
go install github.com/kieranajp/pairings@latest
```

## Usage

### Pair Command
```bash
pairings pair --recipe "https://example.com/recipe"
```

### Preferences Command
```bash
pairings preferences \
  --dish "Beef Bourguignon" \
  --budget-min 2000 \
  --budget-max 5000 \
  --currency EUR \
  --wine-type red \
  --body full \
  --taste-preferences "fruity" "oaky" \
  --occasion "dinner party"
```

### Required Environment Variables

- `GEMINI_API_KEY`: Your Google Gemini API key
  - Get one from [Google AI Studio](https://makersuite.google.com/app/apikey)

### Optional Environment Variables

- `GEMINI_MODEL`: The Gemini model to use (default: "gemini-2.0-flash")
- `LOG_LEVEL`: Logging level (default: "info")
  - Options: debug, info, warn, error

### Command Line Flags

```bash
# Global flags
--gemini-api-key string    Gemini API key
--gemini-model string      Gemini model to use (default: "gemini-2.0-flash")
--log-level string         Log level (debug, info, warn, error) (default: "info")

# Pair command flags
--recipe string           Recipe URL to analyze

# Preferences command flags
--dish string            Name of the dish to pair with
--budget-min int64       Minimum budget in cents (e.g., 2000 for 20.00)
--budget-max int64       Maximum budget in cents (e.g., 5000 for 50.00)
--currency string        Currency code (e.g., EUR, USD) (default: "EUR")
--wine-type string       Preferred wine type (red, white, rose, sparkling)
--body string           Preferred wine body (light, medium, full)
--taste-preferences     Taste preferences (e.g., fruity, dry, oaky)
--occasion string       Occasion context (e.g., dinner party, casual meal)
```

## Development

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run tests:
   ```bash
   go test ./...
   ```

## Configuration

The application uses two configuration files in the `config` directory:
- `schema.json`: Defines the structure of wine pairing responses
- `prompts.yaml`: Contains the prompt templates for the AI

## License

MIT
