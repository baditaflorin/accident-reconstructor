package reconstruction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/baditaflorin/accident-reconstructor/pkg/reconstruct"
)

// LLMClient calls a local Ollama-compatible model for optional narrative text.
type LLMClient struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

// Summary asks the local model for a short cautious case summary.
func (c LLMClient) Summary(ctx context.Context, artifact reconstruct.Artifact) (string, error) {
	if c.BaseURL == "" {
		return "", nil
	}
	client := c.Client
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	prompt := fmt.Sprintf(
		"Summarize this accident reconstruction in three cautious bullet points. Speed %.2f km/h, confidence %.2f, warnings: %s",
		artifact.Speed.MeanKPH,
		artifact.Speed.Confidence,
		strings.Join(artifact.Quality.Warnings, "; "),
	)
	body, _ := json.Marshal(map[string]any{
		"model":  c.Model,
		"prompt": prompt,
		"stream": false,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create ollama request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("call ollama: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("ollama returned %s", resp.Status)
	}
	var parsed struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("decode ollama response: %w", err)
	}
	return strings.TrimSpace(parsed.Response), nil
}
