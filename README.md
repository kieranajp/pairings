# Pairings

A CLI tool that suggests wine pairings for recipes. Given a recipe URL, it analyzes the dish and suggests three wines that would pair well, including detailed reasoning for each pairing.

## Features

- Recipe analysis from URLs
- Structured wine pairing suggestions
- Detailed reasoning for each pairing
- Configurable logging levels
- Support for different Gemini models

## Installation

```bash
go install github.com/kieranajp/pairings@latest
```

## Usage

```bash
pairings pair --recipe "https://example.com/recipe"
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
