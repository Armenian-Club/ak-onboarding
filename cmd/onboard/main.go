package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Armenian-Club/ak-onboarding/internal/app"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/calendar"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/drive"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/mm"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	email := "example@gmail.ru"

	mmClient := mm.NewClient()
	calendarClient := calendar.NewClient()
	driveClient := drive.NewClient()

	onboarder := app.New(mmClient, calendarClient, driveClient)
	if err := onboarder.Onboard(ctx, email, email); err != nil {
		fmt.Printf("onboarding for %v -- finished with err: %v", email, err)
		return
	}
	fmt.Printf("onboarding for %v -- finished successfully\n", email)
	time.Sleep(time.Hour)
}
