package main

import (
	"context"
	"github.com/Armenian-Club/ak-onboarding/internal/app"
	"github.com/Armenian-Club/ak-onboarding/internal/bottg"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/calendar"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/drive"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/mm"
	"github.com/Armenian-Club/ak-onboarding/internal/config"
	"github.com/mymmrac/telego"
	"log"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	jsonCreds, err := os.ReadFile(config.GoogleCredsPath)
	if err != nil {
		log.Fatalf("Failed to read from json creds file drive client: %v", err)
	}
	mmClient := mm.NewClient()
	calendarClient := calendar.NewClient()
	driveClient, err := drive.NewClient(ctx, jsonCreds)

	if err != nil {
		log.Fatalf("Failed to create drive client: %v", err)
	}
	onboarder := app.New(mmClient, calendarClient, driveClient)
	defer cancel()
	botToken := config.BotToken
	bot, err := telego.NewBot(botToken)
	if err != nil {
		log.Fatal(err)
	}
	appBot := bottg.NewBotApp(bot, onboarder)
	err = appBot.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
