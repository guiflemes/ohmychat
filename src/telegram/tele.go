package telegram

import (
	"log"
	"notion-agenda/settings"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Tele() {

	bot, err := tgbotapi.NewBotAPI(settings.GETENV("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Replace USER_ID with the actual user ID you want to initiate the chat with
	userID := 6870062760

	// Send a message to initiate the chat
	msg := tgbotapi.NewMessage(int64(userID), "Hello! This is your bot initiating the chat.")
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

	// Set up an update configuration
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Get updates from the bot
	updates := bot.GetUpdatesChan(u)

	// Handle incoming updates
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Log the received message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Reply to the received message
		replyMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! I received your message.")
		replyMsg.ReplyToMessageID = update.Message.MessageID

		// Send the reply message
		_, err := bot.Send(replyMsg)
		if err != nil {
			log.Println(err)
		}
	}

	// Handle graceful shutdown on interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down the bot...")

}
