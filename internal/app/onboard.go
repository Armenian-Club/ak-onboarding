package app

import (
	"context"

	"github.com/pkg/errors"
)

func (a *app) Onboard(ctx context.Context, email, gmail string) error {
	err := a.mm.InviteToTeam(ctx, email)
	if err != nil {
		// TODO надо бы как-то создать и ловить свои кастомные ошибки чтоб ретраить или не ретраить
		return errors.Wrap(err, "failed to invite to team")
	}
	if err = a.mm.AddUserToChannels(ctx, email); err != nil {
		return errors.Wrap(err, "failed to add users to channel")
	}
	if err = a.dr.AddUser(ctx, gmail); err != nil {
		return errors.Wrap(err, "failed to add user to google drive")
	}
	if err = a.cal.InviteUser(ctx, gmail); err != nil {
		return errors.Wrap(err, "failed to invite user to google calendar")
	}
	return nil
}
