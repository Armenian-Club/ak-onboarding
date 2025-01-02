package drive_test

import (
	"context"
	"testing"
	"time"

	"github.com/Armenian-Club/ak-onboarding/internal/clients/drive"
)

// Тест для примера
func TestClient_AddUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		gmail   string
		wantErr bool
	}{
		{
			name:    "success",
			gmail:   "example@gmail.com",
			wantErr: false,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cl := drive.NewClient()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := cl.AddUser(ctx, tt.gmail)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.InviteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
