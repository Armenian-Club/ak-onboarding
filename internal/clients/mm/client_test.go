package mm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Armenian-Club/ak-onboarding/internal/mocks/mock_mm"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/mattermost/mattermost/server/public/model"
)

// Делай тесты по примеру с календарем и гугл-диском
func TestClient_InviteToTeam(t *testing.T) {
	t.Parallel()
	armClubId := uuid.NewString()
	cases := []struct {
		name    string
		email   string
		mockMm  func(m *mock_mm.Mockhttp)
		wantErr bool
	}{
		{
			name: "success",
			mockMm: func(m *mock_mm.Mockhttp) {
				m.EXPECT().InviteUsersToTeam(gomock.Any(), armClubId, []string{"test@test.com"}).
					Return(
						&model.Response{
							StatusCode: 201,
						}, nil)
			},
			wantErr: false,
			email:   "test@test.com",
		},
		{
			name: "fatal",
			mockMm: func(m *mock_mm.Mockhttp) {
				m.EXPECT().InviteUsersToTeam(gomock.Any(), armClubId, []string{"fatal_test@test.com"}).
					Return(
						&model.Response{
							StatusCode: 401,
						},
						fmt.Errorf("test_error"))
			},
			wantErr: true,
			email:   "fatal_test@test.com",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockHttp := mock_mm.NewMockhttp(ctrl)
			tc.mockMm(mockHttp)
			cl := client{
				modelClient: mockHttp,
				armClubID:   armClubId,
			}
			err := cl.InviteToTeam(ctx, tc.email)
			if (err != nil) != tc.wantErr {
				t.Errorf("client.InviteToTeam() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestFilterChannels(t *testing.T) {
	t.Parallel()
	limit := time.Now().AddDate(0, 0, -180)
	dur, _ := time.ParseDuration("2h30m15s")
	earlyChannel, lateChannel := model.Channel{
		Id:         "7fd9bdhfgjh8",
		LastPostAt: limit.Add(-dur).Unix(),
	}, model.Channel{
		Id:         "57689dsfafs",
		LastPostAt: limit.Add(dur).Unix(),
	}
	cases := []struct {
		name     string
		channels []*model.Channel
		//mockMm		func(m *mock_mm.Mockhttp)
		result []string
	}{
		{
			name: "different channels",
			channels: []*model.Channel{
				&lateChannel, &earlyChannel,
			},
			result: []string{
				lateChannel.Id,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := FilterChannels(tc.channels)
			if len(res) != len(tc.result) {
				t.Errorf("client.FilterChannels() size of result = %v, wanted size %v", len(res), len(tc.result))
			}
			for i := range res {
				if res[i] != tc.result[i] {
					t.Errorf("client.FilterChannels() %v-th result = %v, wanted = %v", i, res[i], tc.result[i])
				}
			}
		})
	}
}

func TestClient_GetChannelsList(t *testing.T) {
	t.Parallel()
	armClubId := uuid.NewString()
	limit := time.Now().AddDate(0, 0, -180)
	dur, _ := time.ParseDuration("2h30m15s")
	successfulResponse, fatalResponse := model.Response{
		StatusCode: 200,
	}, model.Response{
		StatusCode: 400,
	}
	earlyChannel, lateChannel := model.Channel{
		Id:         "7fd9bdhfgjh8",
		LastPostAt: limit.Add(-dur).Unix(),
	}, model.Channel{
		Id:         "57689dsfafs",
		LastPostAt: limit.Add(dur).Unix(),
	}
	cases := []struct {
		name               string
		mockMm             func(m *mock_mm.Mockhttp)
		wantedChannelsList []string
		wantedErr          bool
	}{
		{
			name: "success",
			mockMm: func(m *mock_mm.Mockhttp) {
				m.EXPECT().GetPublicChannelsForTeam(gomock.Any(), armClubId, 0, 1000, "").
					Return(
						[]*model.Channel{
							&lateChannel,
							&earlyChannel,
						},
						&successfulResponse,
						nil,
					)
			},
			wantedChannelsList: []string{
				lateChannel.Id,
			},
			wantedErr: false,
		},
		{
			name: "empty",
			mockMm: func(m *mock_mm.Mockhttp) {
				m.EXPECT().GetPublicChannelsForTeam(gomock.Any(), armClubId, 0, 1000, "").
					Return([]*model.Channel{},
						&successfulResponse,
						nil,
					)
			},
			wantedChannelsList: []string{},
			wantedErr:          false,
		},
		{
			name: "fatal",
			mockMm: func(m *mock_mm.Mockhttp) {
				m.EXPECT().GetPublicChannelsForTeam(gomock.Any(), armClubId, 0, 1000, "").
					Return([]*model.Channel{},
						&fatalResponse,
						fmt.Errorf("failed response"),
					)
			},
			wantedChannelsList: []string{},
			wantedErr:          true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockHttp := mock_mm.NewMockhttp(ctrl)
			tc.mockMm(mockHttp)
			cl := client{
				modelClient: mockHttp,
				armClubID:   armClubId,
			}
			res, err := cl.GetChannelsList(ctx)
			if (err != nil) != tc.wantedErr {
				t.Errorf("client.InviteToTeam() error = %v, wantErr %v", err, tc.wantedErr)
			}
			if len(res) != len(tc.wantedChannelsList) {
				t.Errorf("client.FilterChannels() size of result = %v, wanted size %v", len(res), len(tc.wantedChannelsList))
			}
			for i := range res {
				if res[i] != tc.wantedChannelsList[i] {
					t.Errorf("client.FilterChannels() %v-th result = %v, wanted = %v", i, res[i], tc.wantedChannelsList[i])
				}
			}
		})
	}
}

func TestClient_AddUserToChannels(t *testing.T) {
	t.Parallel()
	armClubId := uuid.NewString()
	limit := time.Now().AddDate(0, 0, -180)
	dur, _ := time.ParseDuration("2h30m15s")
	channelMember := model.ChannelMember{}
	successfulResponse2, fatalResponse := model.Response{
		StatusCode: 201,
	}, model.Response{
		StatusCode: 400,
	}
	earlyChannel, lateChannel := model.Channel{
		Id:         "7fd9bdhfgjh8",
		LastPostAt: limit.Add(-dur).Unix(),
	}, model.Channel{
		Id:         "57689dsfafs",
		LastPostAt: limit.Add(dur).Unix(),
	}
	user := model.User{Id: "jklsdgghkjk89"}
	cases := []struct {
		name    string
		email   string
		mockMm  func(m *mock_mm.Mockhttp)
		wantErr bool
	}{
		{
			name:  "success",
			email: "test@test.com",
			mockMm: func(m *mock_mm.Mockhttp) {
				m.EXPECT().GetPublicChannelsForTeam(gomock.Any(), armClubId, 0, 1000, "").
					Return(
						[]*model.Channel{
							&lateChannel,
							&earlyChannel,
						},
						&successfulResponse2,
						nil,
					)
				m.EXPECT().GetUserByEmail(gomock.Any(), "test@test.com", "").
					Return(
						&user,
						&successfulResponse2,
						nil,
					)
				m.EXPECT().AddChannelMember(gomock.Any(), lateChannel.Id, user.Id).
					Return(
						&channelMember,
						&successfulResponse2,
						nil,
					)
			},
			wantErr: false,
		},
		{
			name:  "fatal adding",
			email: "test@test.com",
			mockMm: func(m *mock_mm.Mockhttp) {
				m.EXPECT().GetPublicChannelsForTeam(gomock.Any(), armClubId, 0, 1000, "").
					Return(
						[]*model.Channel{
							&lateChannel,
							&earlyChannel,
						},
						&successfulResponse2,
						nil,
					)
				m.EXPECT().GetUserByEmail(gomock.Any(), "test@test.com", "").
					Return(
						&user,
						&successfulResponse2,
						nil,
					)
				m.EXPECT().AddChannelMember(gomock.Any(), lateChannel.Id, user.Id).
					Return(
						&channelMember,
						&fatalResponse,
						fmt.Errorf("wrong add channel"),
					)
			},
			wantErr: true,
		},
		{
			name:  "wrong user",
			email: "test@test.com",
			mockMm: func(m *mock_mm.Mockhttp) {
				m.EXPECT().GetPublicChannelsForTeam(gomock.Any(), armClubId, 0, 1000, "").
					Return(
						[]*model.Channel{
							&lateChannel,
							&earlyChannel,
						},
						&successfulResponse2,
						nil,
					)
				m.EXPECT().GetUserByEmail(gomock.Any(), "test@test.com", "").
					Return(
						&user,
						&fatalResponse,
						fmt.Errorf("wrong user"),
					)
			},
			wantErr: true,
		},
		{
			name:  "wrong channels",
			email: "test@test.com",
			mockMm: func(m *mock_mm.Mockhttp) {
				m.EXPECT().GetPublicChannelsForTeam(gomock.Any(), armClubId, 0, 1000, "").
					Return(
						[]*model.Channel{},
						&fatalResponse,
						fmt.Errorf("wrong user"),
					)
				/*m.EXPECT().GetUserByEmail(gomock.Any(), "test@test.com", "").
				Return(
					&user ,
					&successfulResponse2,
					nil,
				)

				*/
			},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockHttp := mock_mm.NewMockhttp(ctrl)
			tc.mockMm(mockHttp)
			cl := client{
				modelClient: mockHttp,
				armClubID:   armClubId,
			}
			err := cl.AddUserToChannels(ctx, tc.email)
			if (err != nil) != tc.wantErr {
				t.Errorf("client.InviteToTeam() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
