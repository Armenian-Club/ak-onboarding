package mm

import (
	"context"

	"github.com/Armenian-Club/ak-onboarding/internal/config"
)

func (c *client) InviteToTeam(ctx context.Context, email string) error {
	//TODO implement me
	_ = config.MMLogin
	_ = config.MMPassword
	return nil
}

func (c *client) AddUserToChannels(ctx context.Context, email string) error {
	//TODO implement me
	return nil
}
