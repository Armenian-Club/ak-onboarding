package mm

import (
	"context"
	
	"github.com/Armenian-Club/ak-onboarding/internal/config"
	"github.com/mattermost/mattermost/server/public/model"
)

// Client интерфейс для работы с mattermost
type Client interface {
	InviteToTeam(ctx context.Context, email string) error
	AddUserToChannels(ctx context.Context, email string) error
}

type client struct {
	modelClient *model.Client4
}

// NewClient конструктор клиента для работы с mattermost
func NewClient() Client {
	myClient := client{modelClient: model.NewAPIv4Client(config.MMBasicUrl)}
	myClient.modelClient.SetToken(config.MMBotAccessToken)
	return &myClient
}
