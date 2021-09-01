package main

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// other imports

type Credentials struct {
	AccessToken       string
	AccessTokenSecret string
	ConsumerKey       string
	ConsumerSecret    string
}

func getClient(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return client, nil
}

func laodAWSEnv() (env Credentials) {
	env = Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}
	return env
}
func filter(unfilterdTweets []twitter.Tweet) (filteredTweets []twitter.Tweet) {
	keyWords := []string{"commission", "commissions", "commissions open"}

	for _, tweet := range unfilterdTweets {
		for _, key := range keyWords {
			if strings.Contains(strings.ToLower(tweet.Text), key) {
				filteredTweets = append(filteredTweets, tweet)
			}
		}
	}
	return
}

func sendDM(client *twitter.Client, tweet twitter.Tweet) {
	client.DirectMessages.EventsNew(&twitter.DirectMessageEventsNewParams{
		Event: &twitter.DirectMessageEvent{
			Type: "message_create",
			Message: &twitter.DirectMessageEventMessage{
				SenderID: "1015612268384587776",
				Target: &twitter.DirectMessageTarget{
					RecipientID: "1015612268384587776",
				},
				Data: &twitter.DirectMessageData{
					Text: "A new opening? in my timeline? what are the odds! \n https://twitter.com/twitter/status/" + tweet.IDStr,
				},
			},
		},
	})

	client.DirectMessages.EventsNew(&twitter.DirectMessageEventsNewParams{
		Event: &twitter.DirectMessageEvent{
			Type: "message_create",
			Message: &twitter.DirectMessageEventMessage{
				SenderID: "1015612268384587776",
				Target: &twitter.DirectMessageTarget{
					RecipientID: "747198468465242112",
				},
				Data: &twitter.DirectMessageData{
					Text: "A new opening? in my timeline? what are the odds! \n https://twitter.com/twitter/status/" + tweet.IDStr,
				},
			},
		},
	})
}

func lookupTweets() {
	creds := laodAWSEnv()

	client, err := getClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
	}

	tweets, _, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: 500,
	})

	if err != nil {
		log.Println(err)
	}

	for _, tweet := range filter(tweets) {
		sendDM(client, tweet)
	}
}

func main() {
	lambda.Start(lookupTweets)
}
