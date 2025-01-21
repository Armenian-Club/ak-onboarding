package calendar

import (
	"context"
	"fmt"

	"google.golang.org/api/calendar/v3"
)

func (c *client) InviteUser(ctx context.Context, gmail string) error {
	rule := &calendar.AclRule{
		Role: "writer",
		Scope: &calendar.AclRuleScope{
			Type:  "user",
			Value: gmail,
		},
	}

	_, err := c.srv.Acl.Insert(c.calendarID, rule).
		SendNotifications(true).
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("failed to insert calendar rule: %w", err)
	}

	return nil
}
