package drive

import "context"

// Client интерфейс для работы с гугл-диском
type Client interface {
	// AddUser добавляет пользователя в АК-ашный диск с возможностью внесения там изменений
	AddUser(ctx context.Context, gmail string) error
}

type client struct {
}

// NewClient конструктор клиента для работы с гугл-диском
func NewClient() Client {
	return &client{}
}
