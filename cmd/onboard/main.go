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

	// Создаём контекст для приложения
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Читаем JSON с учетными данными для Google Calendar
	data, err := os.ReadFile(config.GoogleCredsPath)
	if err != nil {
		log.Fatalf("Could not read credentials file: %v", err)
	}

	// Создаём клиентов
	email := "davidasl085@gmail.com"
	mmClient := mm.NewClient()
	calendarClient, err := calendar.NewClient(ctx, data)
	if err != nil {
		log.Fatalf("Failed to create calendar client: %v", err)
	}
	driveClient := drive.NewClient()

	// Создаём объект onboarder и запускаем процесс
	onboarder := app.New(mmClient, calendarClient, driveClient)
	if err := onboarder.Onboard(ctx, email, email); err != nil {
		fmt.Printf("onboarding for %v -- finished with err: %v\n", email, err)
		return
	}
	fmt.Printf("onboarding for %v -- finished successfully\n", email)
}
