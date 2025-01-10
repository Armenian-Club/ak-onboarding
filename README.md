# OnboardingBot

## How to run
- Install Go
- Create file `.env` in `secrets/` directory
- Fill it with necessary vars. (check it [here](./internal/config/env.go))
```
MM_ARMENIAN_CLUB_ID=example
MM_BASIC_URL=example
MM_BOT_ACCESS_TOKEN=example
```
- build and run project
```bash
go build ./cmd/onboard/main.go
./main
```

