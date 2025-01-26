package drive_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"testing"
	"time"

	"github.com/Armenian-Club/ak-onboarding/internal/clients/drive"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
)

func TestClient_AddUser(t *testing.T) {
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
				_, _ = w.Write([]byte(`{ "kind": "drive#aclRule", "id": "user:test@example.com" }`))
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:  "server error",
			gmail: "test500@example.com",
			httpHandler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`Something went wrong on the server`))
			},
			wantErr:    true,
			wantErrMsg: "failed to add drive permission",
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
				tt.httpHandler(t, w, r)
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Фиктивные креды для теста
			fakeCreds := []byte(`{
			  "type": "service_account",
			  "project_id": "test-448519",
			  "private_key_id": "c297bc0fee46f13e7fd7476065f3cd7b4a060b24",
			  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDW2rHK6+1eM2G9\nlPXTQHMEYlyvR1MKOp1kOsUL9EFiIFAJsTRT7dcKx3tnnIumGhN6L6rw5h6dhwLe\n9SLVdt4ioe089hdp00ayLFSzF3tT6Qq73oAsEoS00EZqgQPmmGL8+3nsMVwdro50\nMR1bb1ai5MdcT7bikh3h/aRgA8A2lnJSeHNt09i1hdUgAHrf8ZlhfZQ0DOqZeT8w\njRcsP4zQt4IO0wRMSdSIuc4yD/9wE1smCaqctEwrHv4PCHp8F0WprsmuR8UFUuNL\nFq/pfIJNj1/JBH63YtAQ5x+pJw2lmogD98w7rzWxQqbnQ5Z2PsRY0D+GSo6uoZj6\nJKqc7XUNAgMBAAECggEACRI8tXOWlwaWVtnGM0AiWwoIHcJmKCVnZcbxcNrEM+9n\nUbFwoyaEkMjxVeOPJdkt/1ep4PfmTQJZRa6V5Ota3510lcFSJb6s2nLytIkGRPmu\n4VW4laPGhJfSkUaXMpI2g7XeZPGEkSBAlXlJYwXlY4VDQYuADjrbFiKOYRtnbyaB\nKUzEgCAWNFIXn5SSKf81K+kBfOtIYhf6qWn1jxe381njzE8iw9Qlmq9mJWlyk1Cj\nPACJVecZ3j6hkViqxNvMJ8V3/VbIO0isBSO7xBCmdJ9lP9rQn9A/W6u9PH8ChKev\nsVquupoR2dEN7rYqpK1Zy6/TsDtqt/BVDPUB3ueWRQKBgQDrmpP7H319TmeG1WTv\ncYQRMLM6D7jNsTM/kIT9+ozH5GxGhccVun26bv57O1zN8cSJ8Pa+yy9n6eQvtoub\nZKYTslknCZdPoUQKceRMn4qvKgei3wCdpXG1qXplcea0/HjbOTB6KFuFiUoltBsl\nbSwnxJYGOX/Z8oPPx0lS+4jcewKBgQDpdESF6owsLjEXpJxO9sQ06aB/KAzwF+Me\nThhKr95Xd3UJjoLnipbFXmFpRzqkbIZHHB1v5Re9jMWCY8Dcut7r4Gyx1Y94EVre\nmh8nE+ADUODYdizpzR6bSl1gumuk1bTu2M6HhnBrbNBIOtPwWyM4p7D7xPbXlOI0\nZsEs5RMSFwKBgDr5vbMtxcbZGncY8aQyYSHAdAzDpLnwcmil73R4BEeBOU1J7XTV\n8uT5JcCJMojmzRDOfaVyzRIQ7Sq4YifqwNvLWB+6eeLX9mU67y/y+88vESxG8CuG\nH3mey+Ga6mpBjKsrnKPneElr/WCEvgrXUic+QWObfxJ6b15Sf1tDVZYDAoGBAK3p\nlZqFrkLDboMEeAVDQ8t/N7dCaND7mpBa8THCbkqOjTu/VLmUvtjthHffPkp7JlUX\nFr7i1Zq5ofGOyoAlHihuGcspIyX5F8641fhQkBMoTzgyYScTTXe2IHYMqmAzbAR6\nsvC3MEx21XrZiEWIP2bXVbtZceIL4a/T1JjTVi+lAoGBAIsHlLAooBVlIfcPV5S0\nh8rv8Do5chvJKqecbX+zmYao0DtvgikPOa4x474I1OoSaad/uXer1XvEiuY/Apnc\ncpJDvNUrwqh5Ixc//yo+4lM5MZJpO+Vgn6ifVEXJjVDMJxGM7WtCB15xYYDtt9n5\nBxZK0L51Prn18R+DKOcuRES9\n-----END PRIVATE KEY-----\n",
			  "client_email": "testpr@test-448519.iam.gserviceaccount.com",
			  "client_id": "104437331581182177415",
			  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
			  "token_uri": "https://oauth2.googleapis.com/token",
			  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/testpr%40test-448519.iam.gserviceaccount.com",
			  "universe_domain": "googleapis.com"
			}`)

			// Создаём клиент с фиктивными кредами
			c, err := drive.NewClient(
				ctx,
				fakeCreds,
				option.WithEndpoint(server.URL),
			)
			require.NoError(t, err, "failed to create calendar client")

			err = c.AddUser(ctx, tt.gmail)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrMsg,
					"error message mismatch, want substring: %q, got: %q",
					tt.wantErrMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}