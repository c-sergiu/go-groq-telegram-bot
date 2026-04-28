package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"telegram-ai-bot-client/internal/config"
)

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqClient struct {
	HTTPClient *http.Client
	History    map[int64][]GroqMessage
	mu         sync.RWMutex
	Config     *config.Config
}

func NewGroqClient(client *http.Client, cfg *config.Config) *GroqClient {
	return &GroqClient{
		HTTPClient: client,
		History:    make(map[int64][]GroqMessage),
		Config:     cfg,
	}
}

func (g *GroqClient) AddMessage(chatID int64, role, content string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.History[chatID]; !exists {
		g.History[chatID] = []GroqMessage{{Role: "system", Content: g.Config.SystemPrompt}}
	}
	
	g.History[chatID] = append(g.History[chatID], GroqMessage{Role: role, Content: content})

	if len(g.History[chatID]) > 11 {
		newHistory := []GroqMessage{g.History[chatID][0]}
		newHistory = append(newHistory, g.History[chatID][len(g.History[chatID])-10:]...)
		g.History[chatID] = newHistory
	}
}

func (g *GroqClient) Ask(chatID int64) (string, error) {
	g.mu.RLock()
	messages := g.History[chatID]
	g.mu.RUnlock()

	payload := map[string]interface{}{
		"model":    g.Config.GroqModel,
		"messages": messages,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", g.Config.GroqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.Config.GroqAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct{ Content string } `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	
	return "", fmt.Errorf("no content returned")
}

func (g *GroqClient) ClearHistory(chatID int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.History, chatID)
}
