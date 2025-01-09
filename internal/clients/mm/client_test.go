package mm

import (
	"context"
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
		// TODO тест кейс на ошибку
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
