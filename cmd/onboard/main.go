package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Armenian-Club/ak-onboarding/internal/app"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/calendar"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/drive"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/mm"
	"github.com/Armenian-Club/ak-onboarding/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	email := "example@gmail.com"
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
}
