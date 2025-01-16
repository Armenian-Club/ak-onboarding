package app_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Armenian-Club/ak-onboarding/internal/app"
	"github.com/Armenian-Club/ak-onboarding/internal/mocks/mock_calendar"
	"github.com/Armenian-Club/ak-onboarding/internal/mocks/mock_drive"
	"github.com/Armenian-Club/ak-onboarding/internal/mocks/mock_mm"
	"github.com/golang/mock/gomock"
)

// TODO: тут только положительный тест-кейс для примера, дополнить другими

func TestApp_Onboard(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		email, gmail string
		wantErr      error
		mockCal      func(c *mock_calendar.MockClient)
		mockDr       func(dr *mock_drive.MockClient)
		mockMm       func(mm *mock_mm.MockClient)
	}{
		{
			name:    "success",
			email:   "test@yandex.ru",
			gmail:   "test@gmail.com",
			wantErr: nil,
			mockCal: func(c *mock_calendar.MockClient) {
				c.EXPECT().InviteUser(gomock.Any(), "test@gmail.com").
					Return(nil)
			},
			mockDr: func(dr *mock_drive.MockClient) {
				dr.EXPECT().AddUser(gomock.Any(), "test@gmail.com").
					Return(nil)
			},
			mockMm: func(mm *mock_mm.MockClient) {
				mm.EXPECT().InviteToTeam(gomock.Any(), "test@yandex.ru").
					Return(nil)
				mm.EXPECT().AddUserToChannels(gomock.Any(), "test@yandex.ru").
					Return(nil)
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cal := mock_calendar.NewMockClient(ctrl)
			dr := mock_drive.NewMockClient(ctrl)
			mm := mock_mm.NewMockClient(ctrl)
			tt.mockCal(cal)
			tt.mockDr(dr)
			tt.mockMm(mm)
			a := app.New(mm, cal, dr)
			err := a.Onboard(ctx, tt.email, tt.gmail)
			if !strings.Contains(fmt.Sprint(err), fmt.Sprint(tt.wantErr)) {
				t.Errorf("Onboard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_AddMmUserAfterJoin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		email   string
		wantErr bool
		mockMm  func(mm *mock_mm.MockClient)
	}{
		{
			name:  "success",
			email: "test@test.com",
			mockMm: func(mm *mock_mm.MockClient) {
				mm.EXPECT().IsUserInTeam(gomock.Any(), "test@test.com").
					Return(true, nil)
				mm.EXPECT().AddUserToChannels(gomock.Any(), "test@test.com").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "fatal",
			email: "test@test.com",
			mockMm: func(mm *mock_mm.MockClient) {
				mm.EXPECT().IsUserInTeam(gomock.Any(), "test@test.com").
					Return(true, nil)
				mm.EXPECT().AddUserToChannels(gomock.Any(), "test@test.com").
					Return(fmt.Errorf("error LoL"))
			},
			wantErr: true,
		},
	}

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cal := mock_calendar.NewMockClient(ctrl)
			dr := mock_drive.NewMockClient(ctrl)
			mm := mock_mm.NewMockClient(ctrl)
			tt.mockMm(mm)
			a := app.New(mm, cal, dr)
			err := a.AddMmUserAfterJoin(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddMmUserAfterJoin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
