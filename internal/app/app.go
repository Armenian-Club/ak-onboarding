package app

import (
	"context"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/calendar"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/drive"
	"github.com/Armenian-Club/ak-onboarding/internal/clients/mm"
)

// Onboarder интерфейс выполняющий онбординг новичка
type Onboarder interface {
	Onboard(ctx context.Context, email, gmail string) error
}

type app struct {
	mm  mm.Client
	cal calendar.Client
	dr  drive.Client
}

// New конструктор Onboarder-а
func New(mm mm.Client, cal calendar.Client, dr drive.Client) Onboarder {
	return &app{
		mm:  mm,
		cal: cal,
		dr:  dr,
	}
}
