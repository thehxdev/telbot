package main

import (
	"fmt"
	"log"

	"github.com/thehxdev/telbot"
)

const BOT_TOKEN string = "your_awesome_bot_token"

func main() {
	bot, err := telbot.New(BOT_TOKEN, "api.telegram.org")
	if err != nil {
		log.Fatal(err)
	}

	conv, err := telbot.NewConversation(telbot.ConversationConfig{
		Bot:    bot,
		Stages: []telbot.UpdateHandlerFunc{startHandler, nameHandler},
	})

	updatesChan, err := bot.StartPolling(telbot.UpdateParams{
		Offset:         0,
		Limit:          100,
		Timeout:        30,
		AllowedUpdates: []string{"message"},
	})

	log.Println("started polling updates...")
	for update := range updatesChan {
		if update.Message.Text == "/start" {
			conv.Start(update)
		}
	}
}

func startHandler(update telbot.Update) error {
	_, err := update.SendMessage(telbot.TextMessageParams{
		ChatId: update.Message.Chat.Id,
		Text:   "Hey! This is a question bot. What is your name?",
	})
	if err != nil {
		return err
	}
	return nil
}

func nameHandler(update telbot.Update) error {
	_, err := update.SendMessage(telbot.TextMessageParams{
		ChatId: update.Message.Chat.Id,
		Text:   fmt.Sprintf("Nice to meet you %s!", update.Message.Text),
	})
	if err != nil {
		return err
	}
	return &telbot.ErrConversationEnd{}
}
