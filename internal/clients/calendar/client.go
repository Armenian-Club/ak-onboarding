package calendar

import (
	"context"
	"fmt"

	"github.com/Armenian-Club/ak-onboarding/internal/config"
	"golang.org/x/oauth2/google"
	gcalendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Client интерфейс для работы с гугл-календарем.
type Client interface {
	// InviteUser отправляет приглашение в календарь.
	InviteUser(ctx context.Context, gmail string) error
}

type client struct {
	srv        *gcalendar.Service
	calendarID string
}

// NewClient создает обычный клиент для продакшена.
func NewClient(ctx context.Context, jsonCreds []byte, opts ...option.ClientOption) (Client, error) {
	creds, err := google.CredentialsFromJSON(ctx, jsonCreds, gcalendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service account credentials: %w", err)
	}

	opts = append([]option.ClientOption{
		option.WithCredentials(creds),
	}, opts...)

	srv, err := gcalendar.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	return &client{
		srv:        srv,
		calendarID: config.CalendarID,
	}, nil
}

// NewClientWithService используется в unit-тестах.
func NewClientWithService(srv *gcalendar.Service, calendarID string) Client {
	return &client{
		srv:        srv,
		calendarID: calendarID,
	}
}
