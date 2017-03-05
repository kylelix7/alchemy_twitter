package main

import (
	"fmt"
	"time"
	"stocker/tweets"
	"stocker/analysis"
  "os"
)

var consumerKey string = os.Getenv("CONSUMER_KEY")
var consumerSecret string = os.Getenv("CONSUMER_SECRET")

func main() {
	key := make([]string, 1)
	key = append(key, "merger acquisition")
	fmt.Print("Keywords: ")
	for _, keyword := range key {
		fmt.Printf("%s ", keyword)
	}
	fmt.Println("")
	tweets, _ := tweets.GetTweets(key, "recent", 100)

	var score float32
	urlSet := make(map[string]bool)

	for count := 0; ;  {
		for _, tweet := range tweets {
			if count >= 50 {
				fmt.Println("Max is reached. exit.")
				return
			}
			if len(tweet.Entities.Urls) != 0 {
				link := tweet.Entities.Urls[0].URL
				if _, ok := urlSet[link]; ok {
					fmt.Printf("Link %s already appeared before, skipping. \n", link)
					continue;
				}
				urlSet[link] = true
				fmt.Println("=========================")
				fmt.Printf("%d: %s \n", count, link)
				count++
				result, _ := analysis.AnalyzeSentimentText(tweet.Entities.Urls[0].URL, "url")
				score = score + result.DocSentiment.Score
				fmt.Printf("this link: %v, current avg: %v\n", result.DocSentiment.Score, score/float32(count))
				fmt.Println("=========================")
			}
		}
		fmt.Println("Sleeping for 60 seconds")
		time.Sleep(60 * time.Second)
	}
	fmt.Println(score)
}
