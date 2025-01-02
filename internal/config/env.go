package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var (
	// MMLogin login for admin account mattermost. MM_LOGIN
	MMLogin string
	// MMPassword password for mattermost. MM_PASSWORD
	MMPassword string
)

func init() {
	err := godotenv.Load("secrets/.env")
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}
	MMLogin = os.Getenv("MM_LOGIN")
	MMPassword = os.Getenv("MM_PASSWORD")
}
