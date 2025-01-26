package drive

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Client интерфейс для работы с гугл-диском
type Client interface {
	// AddUser добавляет пользователя в АК-ашный диск с возможностью внесения там изменений
	AddUser(ctx context.Context, gmail string) error
}
type client struct {
	srv *drive.Service
}

// NewClient конструктор клиента для работы с гугл-диском
func NewClient(ctx context.Context, jsonCreds []byte, opts ...option.ClientOption) (Client, error) {
	// Парсим сервис-аккаунт из JSON
	creds, err := google.CredentialsFromJSON(ctx, jsonCreds, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service account credentials: %w", err)
	}
	opts = append(opts, option.WithCredentials(creds))
	srv, err := drive.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create drive service: %w", err)
	}
	// Возвращаем структуру, реализующую интерфейс Client
	return &client{srv: srv}, nil
}