package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	baseURL      = "https://generativelanguage.googleapis.com/v1beta"
	defaultModel = "gemini-2.0-flash"
)

type GeminiClient struct {
	apiKey string
	model  string
	client *http.Client
}

type geminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func NewGeminiClient(apiKey string, model string) *GeminiClient {
	if model == "" {
		model = defaultModel
	}

	return &GeminiClient{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
	}
}

// extractJSON finds and extracts the first valid JSON object or array from a string
func extractJSON(input string) (string, error) {
	// Find the first { or [ and the last } or ]
	start := -1
	end := -1

	// Look for object start/end
	objStart := strings.Index(input, "{")
	objEnd := strings.LastIndex(input, "}")

	// Look for array start/end
	arrStart := strings.Index(input, "[")
	arrEnd := strings.LastIndex(input, "]")

	// Determine which comes first and is a complete pair
	if objStart != -1 && objEnd != -1 && (arrStart == -1 || objStart < arrStart) {
		start = objStart
		end = objEnd
	} else if arrStart != -1 && arrEnd != -1 {
		start = arrStart
		end = arrEnd
	}

	if start == -1 || end == -1 || start > end {
		return "", fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := input[start : end+1]

	// Validate that it's actually valid JSON
	var js json.RawMessage
	if err := json.Unmarshal([]byte(jsonStr), &js); err != nil {
		return "", fmt.Errorf("extracted content is not valid JSON: %w", err)
	}

	return jsonStr, nil
}

func (c *GeminiClient) GetPairings(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", baseURL, c.model, c.apiKey)

	reqBody := geminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from model")
	}

	rawText := response.Candidates[0].Content.Parts[0].Text

	// Extract JSON from the response
	jsonStr, err := extractJSON(rawText)
	if err != nil {
		return "", fmt.Errorf("failed to extract JSON from model response: %w", err)
	}

	return jsonStr, nil
}
