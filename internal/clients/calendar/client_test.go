package calendar_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	clientcalendar "github.com/Armenian-Club/ak-onboarding/internal/clients/calendar"
	"github.com/stretchr/testify/require"
	gcalendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func TestClient_InviteUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		gmail       string
		httpHandler func(t *testing.T, w http.ResponseWriter, r *http.Request)
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name:  "success",
			gmail: "test@example.com",
			httpHandler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodPost, r.Method)
				require.Contains(t, r.URL.Path, "/acl")

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"kind":"calendar#aclRule",
					"id":"user:test@example.com"
				}`))
			},
		},
		{
			name:  "server error",
			gmail: "test500@example.com",
			httpHandler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Something went wrong"))
			},
			wantErr:    true,
			wantErrMsg: "failed to insert calendar rule",
		},
		{
			name:  "invalid json response",
			gmail: "testinvalid@example.com",
			httpHandler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("not-a-json"))
			},
			wantErr:    true,
			wantErrMsg: "invalid character 'o' in literal null",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				switch {

				case r.Method == http.MethodGet &&
					strings.Contains(r.URL.Path, "/calendars/") &&
					!strings.Contains(r.URL.Path, "/acl"):

					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{
						"id":"test-calendar",
						"summary":"Test calendar"
					}`))

				case r.Method == http.MethodPost &&
					strings.Contains(r.URL.Path, "/acl"):

					tt.httpHandler(t, w, r)

				default:
					t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
				}
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			srv, err := gcalendar.NewService(
				ctx,
				option.WithEndpoint(server.URL),
				option.WithHTTPClient(server.Client()),
				option.WithoutAuthentication(),
			)
			require.NoError(t, err)

			c := clientcalendar.NewClientWithService(
				srv,
				"test-calendar",
			)

			err = c.InviteUser(ctx, tt.gmail)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
