package main

import (
	"log"

	"github.com/thehxdev/telbot"
)

const BOT_TOKEN = "your_awesome_bot_token"

func main() {
	// The host argument is optional
	bot, err := telbot.New(BOT_TOKEN, "api.telegram.org")
	if err != nil {
		log.Fatal(err)
	}

	updatesChan, err := bot.StartPolling(telbot.UpdateParams{
		Offset:         0,
		Timeout:        30,
		Limit:          100,
		AllowedUpdates: []string{"message"},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("started polling updates")
	for update := range updatesChan {
		if update.Message == nil {
			continue
		}
		// Only handle private chats
		if update.Message.Chat.Type != telbot.ChatTypePrivate {
			continue
		}
		go func(update telbot.Update) {
			var err error
			switch update.Message.Text {
			case "/start":
				err = StartHandler(update)
			default:
				err = EchoHandler(update)
			}
			if err != nil {
				log.Println(err)
			}
		}(update)
	}
}

func StartHandler(update telbot.Update) error {
	_, err := update.SendMessage(telbot.TextMessageParams{
		ChatId: update.Message.Chat.Id,
		Text:   "Hello World!",
	})
	return err
}

func EchoHandler(update telbot.Update) error {
	_, err := update.SendMessage(telbot.TextMessageParams{
		ChatId: update.Message.Chat.Id,
		Text:   update.Message.Text,
	})
	return err
}
