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

	updatesChan, _ := bot.StartPolling(context.Background(), telbot.UpdateParams{
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
				// Start a new conversation once a "/start" command received
				conv.Start(startHandler, update)
			default:
				// Otherwise, there is no more routes. So check the userId and
				// chatId for a conversatoin. If they belong to a conversation,
				// call next handler registered with that.
				//
				// NOTE: Ordering of the handlers matter! `telbot` is a low
				// level library that does not provide any routing. So routing
				// the updates and conversatoins must be handled by the user.
				if conv.HasConversation(update.ChatId(), update.UserId()) {
					err = conv.CallNext(update)
				}
			}
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func startHandler(c *conv.Conversation, update telbot.Update) error {
	params := telbot.TextMessageParams{
		ChatId: update.Message.Chat.Id,
		Text:   "Hey! This is a question bot. What is your name?",
	}
	_, err := update.Bot.SendMessage(context.Background(), params)
	c.Next = nameHandler
	return err
}

func nameHandler(c *conv.Conversation, update telbot.Update) error {
	params := telbot.TextMessageParams{
		ChatId: update.Message.Chat.Id,
		Text:   fmt.Sprintf("Nice to meet you %s!", update.Message.Text),
	}
	update.Bot.SendMessage(context.Background(), params)
	return &conv.EndConversation{}
}
