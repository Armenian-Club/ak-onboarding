package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

func (a *app) Onboard(ctx context.Context, email, gmail string) error {
	err := a.mm.InviteToTeam(ctx, email)
	if err != nil {
		// TODO надо бы как-то создать и ловить свои кастомные ошибки чтоб ретраить или не ретраить
		return errors.Wrap(err, "failed to invite to team")
	}

	go func(email string) {
		err := a.AddMmUserAfterJoin(email)
		if err != nil {
			log.Println(err)
		}
	}(email)

	if err = a.dr.AddUser(ctx, gmail); err != nil {
		return errors.Wrap(err, "failed to add user to google drive")
	}
	if err = a.cal.InviteUser(ctx, gmail); err != nil {
		return errors.Wrap(err, "failed to invite user to google calendar")
	}
	return nil
}

func (a *app) AddMmUserAfterJoin(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Hour)
	defer cancel()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			localInTeam, err := a.mm.IsUserInTeam(ctx, email)
			if err != nil {
				log.Println(fmt.Errorf("failed to check user in team %w", err))
			}
			if localInTeam {
				if err = a.mm.AddUserToChannels(ctx, email); err != nil {
					return errors.Wrap(err, "failed to add users to channel")
				}
				return nil
			}
			log.Println("User is not in the team.")
		case <-ctx.Done():
			return errors.New("time is out")
		}
		log.Printf("User %v is added to mattermost", email)
	}
}
