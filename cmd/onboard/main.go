package main

import (
	"context"
	"fmt"
	"github.com/Armenian-Club/ak-onboarding/internal/bottg"
	"github.com/Armenian-Club/ak-onboarding/internal/config"
	"github.com/mymmrac/telego"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	botToken := config.BotToken
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bottg.Run(ctx, bot)
	/*email := "example@gmail.com"
	jsonCreds, err := os.ReadFile(config.GoogleCredsPath)
	if err != nil{
		log.Fatalf("Failed to read from json creds file drive client: %v", err)
	}
	mmClient := mm.NewClient()
	calendarClient := calendar.NewClient()
	driveClient, err := drive.NewClient(ctx, jsonCreds)

	if err != nil {
		log.Fatalf("Failed to create drive client: %v", err)
	}
	onboarder := app.New(mmClient, calendarClient, driveClient)
	if err := onboarder.Onboard(ctx, email, email); err != nil {
		fmt.Printf("onboarding for %v -- finished with err: %v", email, err)
		return
	}
	fmt.Printf("onboarding for %v -- finished successfully\n", email)
	*/
}
