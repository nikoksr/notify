package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/twitter"
)

type TwitterConf struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	TwitterId         string
	MessageTitle      string
	Message           string
}

func sendTwitterDM(conf TwitterConf) {
	twitterService, _ := twitter.New(twitter.Credentials{
		ConsumerKey:       conf.ConsumerKey,
		ConsumerSecret:    conf.ConsumerSecret,
		AccessToken:       conf.AccessToken,
		AccessTokenSecret: conf.AccessTokenSecret,
	})

	twitterService.AddReceivers(conf.TwitterId)
	notifier := notify.New()
	notifier.UseServices(twitterService)

	_ = notifier.Send(context.Background(), conf.MessageTitle, conf.Message)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	conf := TwitterConf{
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("ACCESS_TOKEN"),
		os.Getenv("ACCESS_TOKEN_SECRET"),
		os.Getenv("TWITTER_ID"),
		"YOUR MESSAGE TITLE HERE",
		"YOUR MESSAGE BODY HERE",
	}
	sendTwitterDM(conf)
}
