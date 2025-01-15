package mm

import (
	"context"

	"github.com/Armenian-Club/ak-onboarding/internal/config"
	"github.com/mattermost/mattermost/server/public/model"
)

// Client интерфейс для работы с mattermost
type Client interface {
	InviteToTeam(ctx context.Context, email string) error
	IsUserInTeam(ctx context.Context, email string) (bool, error)
	AddUserToChannels(ctx context.Context, email string) error
}

type http interface {
	InviteUsersToTeam(ctx context.Context, teamId string, userEmails []string) (*model.Response, error)
	GetPublicChannelsForTeam(ctx context.Context, teamId string, page int, perPage int, etag string) ([]*model.Channel, *model.Response, error)
	GetUserByEmail(ctx context.Context, email, etag string) (*model.User, *model.Response, error)
	AddChannelMember(ctx context.Context, channelId, userId string) (*model.ChannelMember, *model.Response, error)
	GetTeamMember(ctx context.Context, teamId string, userId string, etag string) (*model.TeamMember, *model.Response, error)
}

type client struct {
	modelClient http
	armClubID   string
}

// NewClient конструктор клиента для работы с mattermost
func NewClient() Client {
	c := model.NewAPIv4Client(config.MMBasicUrl)
	c.SetToken(config.MMBotAccessToken)
	myClient := client{
		modelClient: c,
		armClubID:   config.MMArmenianClubId,
	}
	c.InviteUsersToTeam()
	return &myClient
}
