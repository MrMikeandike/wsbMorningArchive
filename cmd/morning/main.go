package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MrMikeandike/wsbMorningArchive/pkg/report"

	"github.com/turnage/graw/reddit"
)

func main() {
	id, _ := os.LookupEnv("CLIENT_ID")
	secret, _ := os.LookupEnv("CLIENT_SECRET")
	username, _ := os.LookupEnv("USER_NAME")
	password, _ := os.LookupEnv("PASSWORD")
	agent, _ := os.LookupEnv("USER_AGENT")

	cfg := reddit.BotConfig{
		Agent: agent,
		App: reddit.App{
			ID:       id,
			Secret:   secret,
			Username: username,
			Password: password,
		},
	}
	fmt.Println(cfg)
	bot, err := reddit.NewBot(cfg)
	if err != nil {
		log.Fatal(err)
	}
	//page, err := GetPageBefore(bot, "t3_cmzm4z", "23")
	page, err := report.GetSpecific(bot, []string{
		"t3_d0gt23",
		"t3_d01c9u",
		"t3_czl07a",
	})
	if err != nil {
		log.Fatal(err)
	}
	if page.PostCount == 0 {
		log.Fatal("post count == 0")
	}
	// fmt.Printf("First: %s\n", page.FirstID)
	// fmt.Printf("Last: %s\n", page.LastID)
	// fmt.Printf("Count: %d\n", page.PostCount)
	//for i, v := range page.Reports {
	for i, v := range page.Reports {
		fmt.Println("----------------------")
		fmt.Printf("INDEX: %d\n", i)
		fmt.Printf("Title: %s\n", v.Title)
		fmt.Printf("ID   : %s\n", v.FullID)
		fmt.Printf("Time : %s\n", v.DateTime.String())
		fmt.Println("----------------------")
	}

}
