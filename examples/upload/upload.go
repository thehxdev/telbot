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

	// upload source code
	file, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}

	fileInfo := &telbot.FileReader{
		Kind:     "document",
		FileName: filepath.Base(*path),
		Reader:   file,
	}
	params := telbot.UploadParams{
		ChatId: *chatid,
		Method: "sendDocument",
	}

	log.Println("uploading file", *path)
	msg, err := bot.UploadFile(context.Background(), params, []telbot.FileInfo{fileInfo})
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	fmt.Printf("\n%#v\n\n%#v\n", msg, msg.Document)
}
