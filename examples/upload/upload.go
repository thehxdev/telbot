package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/thehxdev/telbot"
)

const BOT_TOKEN = "your_awesome_bot_token"

func main() {
	path := flag.String("path", "", "path to file")
	chatid := flag.Int("chatid", -1, "target chat id")
	flag.Parse()

	if *path == "" {
		log.Fatal("file path is not specified")
	}

	// The host argument is optional
	bot, err := telbot.New(BOT_TOKEN, "api.telegram.org")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// upload source code
	file, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}

	fi := &telbot.FileReader{
		Kind:     "document",
		FileName: filepath.Base(*path),
		Reader:   file,
	}

	log.Println("uploading file", *path)
	msg, err := bot.UploadFile(ctx, telbot.UploadParams{
		// send file to bot chat
		ChatId: *chatid,
		Method: "sendDocument",
	}, fi)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%#v\n\n%#v\n", msg, msg.Document)
}
