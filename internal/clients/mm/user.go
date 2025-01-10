package mm

import (
	"context"
	"fmt"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
)

func (c *client) InviteToTeam(ctx context.Context, email string) error {
	response, err := c.modelClient.InviteUsersToTeam(ctx, c.armClubID, []string{email})
	if response == nil {
		return fmt.Errorf("failed to make response: %w; response: %v", err, nil)
	}
	if err != nil || response.StatusCode/100 != 2 {
		return fmt.Errorf("failed to make response: %w; status code: %v", err, response.StatusCode)
	}
	return nil
}

func FilterChannels(allChannels []*model.Channel) []string {
	limit := time.Now().AddDate(0, 0, -180)
	var channelIds []string
	for _, channel := range allChannels {
		if !time.Unix(channel.LastPostAt, 0).Before(limit) {
			channelIds = append(channelIds, channel.Id)
		}
	}
	return channelIds
}

func (c *client) GetChannelsList(ctx context.Context) ([]string, error) {
	channels, response, err := c.modelClient.GetPublicChannelsForTeam(ctx, c.armClubID, 0, 1000, "")
	if response == nil {
		return nil, fmt.Errorf("failed to make response: %w; response: %v", err, nil)
	}
	if err != nil || response.StatusCode/100 != 2 {
		// TODO и тут
		return nil, fmt.Errorf("failed to make response: %w; status code: %v", err, response.StatusCode)
	}
	return FilterChannels(channels), nil
}

func (c *client) AddUserToChannels(ctx context.Context, email string) error {
	channelsList, err := c.GetChannelsList(ctx)
	if err != nil {
		return fmt.Errorf("failed to get channels list: %w", err)
	}
	user, response, err := c.modelClient.GetUserByEmail(ctx, email, "")
	if response == nil {
		return fmt.Errorf("failed to get user by id: %w; response: %v", err, nil)
	}
	if err != nil || response.StatusCode/100 != 2 {
		return fmt.Errorf("failed to get user by id: %w; status code: %v", err, response.StatusCode)
	}
	userId := user.Id

	for _, channelId := range channelsList {
		_, response, err = c.modelClient.AddChannelMember(ctx, channelId, userId)
		if response == nil {
			return fmt.Errorf("failed to invite user to the channel with id %v: %w; response: %v", channelId, err, nil)
		}
		if err != nil || response.StatusCode/100 != 2 {
			return fmt.Errorf("failed to invite user to the channel with id %v: %w; status code: %v", channelId, err, response.StatusCode)
		}
	}
	return nil
}
