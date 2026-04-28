package main

import (
	"log"
	"net/http"
	"telegram-ai-bot-client/internal/ai"
	"telegram-ai-bot-client/internal/config"
	"telegram-ai-bot-client/internal/telegram"
	"time"
)

func main() {
	cfg := config.Load()
	httpClient := &http.Client{Timeout: 45 * time.Second}
	
	tgClient := telegram.NewClient(httpClient, cfg)
	aiClient := ai.NewGroqClient(httpClient, cfg)

	offset := 0
	log.Printf("[INFO] Service started. Model: %s. Polling Telegram...", cfg.GroqModel)

	for {
		updates, err := tgClient.GetUpdates(offset)
		if err != nil {
			log.Printf("[ERROR] Polling: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates {
			handleUpdate(tgClient, aiClient, update)
			offset = update.UpdateID + 1
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func handleUpdate(tg *telegram.Client, aiClient *ai.GroqClient, update telegram.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[RECOVER] Critical error handling update %d: %v", update.UpdateID, r)
		}
	}()

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	text := update.Message.Text

	if text == "" {
		return
	}

	log.Printf("[RECV] [%d] From @%s", chatID, update.Message.From.Username)

	switch text {
	case "/start":
		tg.SendMessage(chatID, "Hello! I'm your AI assistant. Send me a message or use /clear to reset our chat.")
		return
	case "/clear":
		if !tg.IsAdmin(chatID, userID) {
			log.Printf("[AUTH] Denied /clear for user @%s in chat %d", update.Message.From.Username, chatID)
			tg.SendMessage(chatID, "🚫 Only group admins can clear my memory.")
			return
		}

		aiClient.ClearHistory(chatID)
		tg.SendMessage(chatID, "🧹 History cleared! What can I help you with now?")
		return
	}

	aiClient.AddMessage(chatID, "user", text)
	
	response, err := aiClient.Ask(chatID)
	if err != nil {
		log.Printf("[ERROR] AI: %v", err)
		tg.SendMessage(chatID, "Sorry, I'm having trouble thinking right now. Try again later.")
		return
	}

	aiClient.AddMessage(chatID, "assistant", response)

	if err := tg.SendMessage(chatID, response); err != nil {
		log.Printf("[ERROR] Sending: %v", err)
	}
}
