package main

import (
	"context"
	"fmt"
	"log"

	"github.com/thehxdev/telbot"
	conv "github.com/thehxdev/telbot/ext/conversation"
)

const BOT_TOKEN string = "your_awesome_bot_token"

func main() {
	bot, err := telbot.New(BOT_TOKEN, "api.telegram.org")
	if err != nil {
		log.Fatal(err)
	}

	startConv, _ := conv.New([]telbot.UpdateHandlerFunc{startHandler, nameHandler})

	ctx := context.Background()
	updatesChan, err := bot.StartPolling(ctx, telbot.UpdateParams{
		Offset:         0,
		Limit:          100,
		Timeout:        30,
		AllowedUpdates: []string{"message"},
	})

	log.Println("started polling updates...")
	for update := range updatesChan {
		if update.Message == nil {
			continue
		}
		go func() {
			var err error
			switch update.Message.Text {
			case "/start":
				// If a message with text `/start` comes, start a conversation.
				startConv.Start(update)
			default:
				// Otherwise, there is no more routes. So check all other updates for conversation.
				// If they belong to a conversation, handle that.
				//
				// NOTE: Ordering of the handlers matter! `telbot` is a low-level library that
				// does not provide any routing. So the user must handle routing of the updates
				// and conversations.
				if conv.HasConversation(update.ChatId(), update.UserId()) {
					err = conv.HandleUpdate(update)
				}
			}
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func startHandler(update telbot.Update) error {
	params := telbot.TextMessageParams{
		ChatId: update.Message.Chat.Id,
		Text:   "Hey! This is a question bot. What is your name?",
	}
	_, err := update.Bot.SendMessage(context.Background(), params)
	return err
}

func nameHandler(update telbot.Update) error {
	params := telbot.TextMessageParams{
		ChatId: update.Message.Chat.Id,
		Text:   fmt.Sprintf("Nice to meet you %s!", update.Message.Text),
	}
	_, err := update.Bot.SendMessage(context.Background(), params)
	if err != nil {
		return err
	}
	return conv.EndConversation
}
