package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"telegram-ai-bot-client/internal/config"
)

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		From struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
		} `json:"from"`
	} `json:"message"`
}

type Client struct {
	HTTPClient *http.Client
	Config     *config.Config
}

func NewClient(httpClient *http.Client, cfg *config.Config) *Client {
	return &Client{
		HTTPClient: httpClient,
		Config:     cfg,
	}
}

func (c *Client) GetUpdates(offset int) ([]Update, error) {
	url := fmt.Sprintf("%s%s/getUpdates?offset=%d&timeout=20", c.Config.TelegramURL, c.Config.TelegramToken, offset)
	
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Result []Update `json:"result"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data.Result, nil
}

func (c *Client) SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("%s%s/sendMessage", c.Config.TelegramURL, c.Config.TelegramToken)
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api error: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) IsAdmin(chatID int64, userID int64) bool {
	if chatID == userID {
		return true
	}

	url := fmt.Sprintf("%s%s/getChatMember?chat_id=%d&user_id=%d", c.Config.TelegramURL, c.Config.TelegramToken, chatID, userID)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var data struct {
		Result struct {
			Status string `json:"status"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false
	}

	s := data.Result.Status
	return s == "creator" || s == "administrator"
}
