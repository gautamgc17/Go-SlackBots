package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
	}
}

func main() {
	e := godotenv.Load(".env")
	if e != nil {
		log.Fatal("Error Loading .env File", e)
	}
	// os.Setenv("SLACK_APP_TOKEN", "")
	// os.Setenv("SLACK_BOT_TOKEN", "")

	// NewClient creates a new client using the Slack API
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	// CommandEvents returns read only command events channel
	go printCommandEvents(bot.CommandEvents())

	// Command define a new command and append it to the list of existing commands
	bot.Command("My yob is <year>", &slacker.CommandDefinition{
		Description: "YOB Calculator",
		Example:     "My YOB is 2020",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				fmt.Println(err)
			}
			curr := time.Now().Year()
			age := curr - yob
			// formats according to a format specifier and returns the resulting string.
			r := fmt.Sprintf("Age is %d", age)
			response.Reply(r)
		},
	})

	// The returned context's Done channel is closed when the returned cancel function is called or when the parent context's Done channel is closed, whichever happens first.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
