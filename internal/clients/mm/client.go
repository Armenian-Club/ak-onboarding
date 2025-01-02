package mm

import (
	"context"
)

// Client интерфейс для работы с mattermost
type Client interface {
	InviteToTeam(ctx context.Context, email string) error
	AddUserToChannels(ctx context.Context, email string) error
}

type client struct {
}

// NewClient конструктор клиента для работы с mattermost
func NewClient() Client {
	return &client{}
}
