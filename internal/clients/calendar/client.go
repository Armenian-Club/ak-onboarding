package calendar

import "context"

// Client интерфейс для работы с гугл-календарем
type Client interface {
	// InviteUser отправляет приглашение в АК-ашный календарь с возможностью внесения там изменений
	InviteUser(ctx context.Context, gmail string) error
}

type client struct {
}

// NewClient конструктор клиента для работы с гугл-календарем
func NewClient() Client {
	return &client{}
}
