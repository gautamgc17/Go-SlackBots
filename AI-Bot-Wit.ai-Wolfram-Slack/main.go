package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Krognol/go-wolfram"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
)

// var wolframClient *wolfram.Client

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	// CommandEvent is an event to capture executed commands
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}


func main() {
	// Load will read your env file(s) and load them into ENV for this process
	e := godotenv.Load(".env")
	if e != nil {
		log.Fatal("Error Loading .env File", e)
	}


	// NewClient creates a new client using the Slack API
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	fmt.Printf("Type of bot is: %T \n", bot)


	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}


	// Prints the events the bot is (given command) subscribed to
	// CommandEvents returns read only command events channel
	go printCommandEvents(bot.CommandEvents())


	// Command define a new command and append it to the list of existing commands
	bot.Command("query for bot - <message>", &slacker.CommandDefinition{
		Description: "send any question to wolfram",
		Example:     "who is president of india",
		// A BotContext interface is used to respond to an event
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("message")
			fmt.Println("Query is: ", query)

			// Parse - parses query text to witai and returns entities
			msg, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})
			fmt.Println("Message: ", msg)

			data, _ := json.MarshalIndent(msg, "", "    ")
			rough := string(data)
			fmt.Println("Data Received: ", rough)

			value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			fmt.Println("Value: ", value)

			question := value.String()
			res, error := wolframClient.GetSpokentAnswerQuery(question, wolfram.Metric, 1000)
			if error!= nil{
				log.Fatal("Error in getting Response: ", error)
			}
			fmt.Println("Response is: ", res)
			response.Reply(res)
		},
	})

	
	// WithCancel returns a copy of parent with a new Done channel.
	// The returned context's Done channel is closed when the returned cancel function is called or when the parent context's Done channel is closed, whichever happens first.
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Printf("Context Type: %T and Cancel Type: %T \n", ctx, cancel)
	fmt.Println(ctx, "\t", cancel)
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
