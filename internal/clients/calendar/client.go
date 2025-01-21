package calendar

import (
	"context"
	"fmt"
	"github.com/Armenian-Club/ak-onboarding/internal/config"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Client интерфейс для работы с гугл-календарем
type Client interface {
	// InviteUser отправляет приглашение в АК-ашный календарь с возможностью внесения там изменений
	InviteUser(ctx context.Context, gmail string) error
}

type client struct {
	srv        *calendar.Service
	calendarID string
}

// NewClient конструктор клиента для работы с гугл-календарем
func NewClient(ctx context.Context, jsonCreds []byte, opts ...option.ClientOption) (Client, error) {
	creds, err := google.CredentialsFromJSON(ctx, jsonCreds, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service account credentials: %w", err)
	}

	// Собираем общий список опций: credentials + то, что пришло извне.
	baseOpts := []option.ClientOption{
		option.WithCredentials(creds),
	}
	baseOpts = append(baseOpts, opts...)

	srv, err := calendar.NewService(ctx, baseOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	return &client{
		srv:        srv,
		calendarID: config.CalendarID,
	}, nil
}
