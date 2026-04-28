package config

import (
	"log"
	"os"
	"strings"
	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	TelegramURL   string
	GroqAPIKey    string
	GroqModel     string
	GroqURL       string
	SystemPrompt  string
}

func Load() *Config {
	_ = godotenv.Load()

	conf := &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		GroqAPIKey:    os.Getenv("GROQ_API_KEY"),
		TelegramURL:   getEnv("TELEGRAM_URL", "https://api.telegram.org/bot"),
		GroqURL:       getEnv("GROQ_URL", "https://api.groq.com/openai/v1/chat/completions"),
		GroqModel:     getEnv("GROQ_MODEL", "llama-3.3-70b-versatile"),
	}

	prompt, source := loadSystemPromptWithSource("system-prompt.txt")
	conf.SystemPrompt = prompt

	if conf.TelegramToken == "" || conf.GroqAPIKey == "" {
		log.Fatal("[CRITICAL] Missing required environment variables: TELEGRAM_TOKEN and GROQ_API_KEY must be set.")
	}

	log.Println("--------------------------------------------------")
	log.Println("[CONFIG] Settings loaded successfully:")
	log.Printf("[CONFIG] Telegram URL:  %s", conf.TelegramURL)
	log.Printf("[CONFIG] Groq URL:      %s", conf.GroqURL)
	log.Printf("[CONFIG] Groq Model:    %s", conf.GroqModel)
	log.Printf("[CONFIG] Prompt Source: %s", source)
	log.Printf("[CONFIG] Prompt (val):  %s", truncatePrompt(conf.SystemPrompt, 50))
	log.Println("[CONFIG] Auth Status:   Keys loaded (not displayed for security)")
	log.Println("--------------------------------------------------")

	return conf
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func loadSystemPromptWithSource(filename string) (string, string) {
	data, err := os.ReadFile(filename)
	if err == nil {
		content := strings.TrimSpace(string(data))
		if content != "" {
			return content, "External File (" + filename + ")"
		}
	} else {
		log.Printf("[INFO] Could not find or read %s: %v. Falling back to other sources.", filename, err)
	}

	envPrompt := os.Getenv("SYSTEM_PROMPT")
	if envPrompt != "" {
		return envPrompt, "Environment Variable (SYSTEM_PROMPT)"
	}

	return "Be a helpful agent! Keep your answers short and to the point.", "Hardcoded Default"
}

func truncatePrompt(p string, max int) string {
	if len(p) <= max {
		return p
	}
	return p[:max] + "..."
}
