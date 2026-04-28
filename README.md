# Telegram AI Bot (Go + Groq)

A high-performance, concurrent Telegram bot powered by Groq's Llama 3. Designed for speed, thread-safety, and easy deployment.

## 🚀 How to Run

You can choose one of the two methods below to get started:

### Option A: Download the Executable (No Coding Required)
If you just want to run the bot without installing Go:
1. **Download**: Go to the **Releases** section of this repository and download the zip file for your OS (Windows or Linux).
2. **Extract**: Unzip the folder.
3. **Configure**: Open the .env file (included in the zip) and add your API keys.
4. **Run**: Double-click the bot-client executable.

### Option B: Build from Source (For Developers)
If you have Go installed and want to modify or compile the code yourself:
1. **Clone** this repository.
2. **Build**: Run 'make build' (or 'make windows' / 'make linux') to create the binary.
3. **Run**: Execute './bin/bot-client'.

---

## ⚙️ Configuration

Regardless of how you run it, the bot requires a .env file in the same directory.
Open .env and fill in the following:

TELEGRAM_TOKEN=your_telegram_bot_token
GROQ_API_KEY=your_groq_api_key
GROQ_MODEL=llama3-8b-8192
GROQ_URL=https://api.groq.com/openai/v1/chat/completions
TELEGRAM_URL=https://api.telegram.org/bot
SYSTEM_PROMPT="You are a helpful AI assistant."

---

## 🎮 Features & Commands

- **Context Memory**: The bot intelligently tracks the last 10 messages of your conversation to provide context-aware answers.
- **Concurrent & Fast**: Built with Go routines and RWMutex, handling multiple users and groups simultaneously without lag.
- **Commands**:
  - /start: Initializes the bot and sends a greeting.
  - /clear: Wipes the conversation history for the current chat. 
    - Note: In group chats, this command is strictly restricted to Administrators.

## 📁 Project Structure

- cmd/bot/: Application entry point and main polling loop.
- internal/ai/: Groq API client with sliding-window history and thread-safe logic.
- internal/telegram/: Telegram API client with dynamic admin permission checking.
- internal/config/: Environment variable management.
