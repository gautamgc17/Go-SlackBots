package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	e := godotenv.Load(".env")
	if e != nil {
		log.Fatal("Error Loading .env File", e)
	}

	// New builds a slack client or api connection to slack API from the provided token and options.
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	channelArr := []string{os.Getenv("CHANNEL_ID")}
	fileArr := []string{"sample.pdf" , "ZIPL.pdf"}

	// Loop through the fles and upload using given params 
	for i := 0; i<len(fileArr); i++{
		params := slack.FileUploadParameters{
			Channels: channelArr,
			File: fileArr[i],
		}

		file, err := api.UploadFile(params)
		if err != nil{
			fmt.Println("Error in uploading File: ", err)
			return 
		}
		fmt.Printf("File-name: %s, Size: %d, Type: %s, URL: %s \n", file.Name, file.Size, file.Mimetype, file.URLPrivateDownload)
	}
}
