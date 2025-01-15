package app

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

func (a *app) Onboard(ctx context.Context, email, gmail string) error {
	err := a.mm.InviteToTeam(ctx, email)
	if err != nil {
		// TODO надо бы как-то создать и ловить свои кастомные ошибки чтоб ретраить или не ретраить
		return errors.Wrap(err, "failed to invite to team")
	}
	inTeam := make(chan bool)
	go func() {
		for {
			localInTeam, err := a.mm.IsUserInTeam(ctx, email)
			if err != nil {
				fmt.Println(fmt.Errorf("failed to check user in team %w", err))
			}
			if localInTeam {
				inTeam <- true
				return
			}
			time.Sleep(time.Minute)
		}
	}()
	select {
	case <-inTeam:
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
}
