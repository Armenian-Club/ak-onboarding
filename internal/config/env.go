package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var (
	// MMArmenianClubId The id of our team
	MMArmenianClubId string
	// MMBasicUrl Basic url of queries
	MMBasicUrl       string
	MMBotAccessToken string
)

func init() {
	err := godotenv.Load("secrets/.env")
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}
	MMArmenianClubId = os.Getenv("MM_ARMENIAN_CLUB_ID")
	MMBasicUrl = os.Getenv("MM_BASIC_URL")
	MMBotAccessToken = os.Getenv("MM_BOT_ACCESS_TOKEN")
}
