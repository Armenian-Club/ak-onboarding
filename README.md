# OnboardingBot

## How to run
- Install Go
- Create file `.env` in `secrets/` directory
- Put google-drive service account creds json-file into secrets/ as ak-drive-creds.json
- Fill it with necessary vars. (check it [here](./internal/config/env.go))
```
MM_ARMENIAN_CLUB_ID=example
MM_BASIC_URL=example
MM_BOT_ACCESS_TOKEN=example
FOLDER_ID=your_gdrive_folder_id
BOT_TOKEN=example
ADMIN_ID=admin_telegram_chat_id
SYSADMIN_TAG=sysadmin_telegram_tag
CALENDAR_ID=google_calendar_id
START_INFO_PAGE_URL=https://outline.armenianclub.org/doc/start-armyanskij-klub-gtHebxatkR
MATTERMOST_INSTRUCTIONS_URL=https://outline.armenianclub.org/s/9814ee83-3a0e-4e7d-872f-c767d2216558
GOOGLE_DRIVE_INSTRUCTIONS_URL=https://outline.armenianclub.org/s/30b3026a-b656-4b1f-9415-d775effdcf22


```
- build and run project
```bash
go build ./cmd/onboard/main.go
./main
```
