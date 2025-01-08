package mm

import (
	"context"
	"fmt"
	"github.com/Armenian-Club/ak-onboarding/internal/config"
	"github.com/mattermost/mattermost/server/public/model"
	"time"
)

func (c *client) InviteToTeam(ctx context.Context, email string) error {
	//TODO implement me
	//mmLogin := config.MMLogin
	//mmPassword := config.MMPassword
	response, err := c.modelClient.InviteUsersToTeam(ctx, config.MMArmenianClubId, []string{email})
	if err != nil || response.StatusCode/100 != 2 {
		return fmt.Errorf("failed to make response: %w; status code: %v", err, response.StatusCode)
	}
	return nil
}

func filterChannels(allChannels []*model.Channel) []string {
	limit := time.Now().AddDate(0, 0, -180)
	var channelIds []string
	for _, channel := range allChannels {
		if !time.Unix(int64(channel.LastPostAt), 0).Before(limit) {
			channelIds = append(channelIds, channel.Id)
		}
	}
	return channelIds
}

func (c *client) GetChannelsList(ctx context.Context) ([]string, error) {
	channels, response, err := c.modelClient.GetPublicChannelsForTeam(ctx, config.MMArmenianClubId, 0, 1000, "")
	if err != nil || response.StatusCode/100 != 2 {
		return nil, fmt.Errorf("failed to make response: %w; status code: %v", err, response.StatusCode)
	}
	return filterChannels(channels), nil
}

func (c *client) AddUserToChannels(ctx context.Context, email string) error {
	//TODO implement me
	channelsList, err := c.GetChannelsList(ctx)
	if err != nil {
		return fmt.Errorf("failed to get channels list: %w", err)
	}
	user, response, err := c.modelClient.GetUserByEmail(ctx, email, "")
	if err != nil || response.StatusCode/100 != 2 {
		return fmt.Errorf("failed to get user by id: %w; status code: %v", err, response.StatusCode)
	}
	userId := user.Id

	for _, channelId := range channelsList {
		_, response, err = c.modelClient.AddChannelMember(ctx, channelId, userId)
		if err != nil || response.StatusCode/100 != 2 {
			return fmt.Errorf("failed to invite user to the channel with id %v: %w; status code: %v", channelId, err, response.StatusCode)
		}
	}
	//c.modelClient.AddChannelMember()
	//_, response, err := c.modelClient.InviteUsersToTeamAndChannelsGracefully(ctx, config.MMArmenianClubId, []string{email}, channelsList, "Message")
	//if err != nil || response.StatusCode/100 != 2 {
	//	return fmt.Errorf("failed to invite users: %w; status code: %v", err, response.StatusCode)
	//}
	return nil
}
