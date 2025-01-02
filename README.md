# OnboardingBot

## How to run
- Install Go
- Create file `.env` in `secrets/` directory
- Fill it with necessary vars. (check it [here](./internal/config/env.go))
```
MM_LOGIN=example
MM_PASSWORD=example
```
- build and run project
```bash
go build ./cmd/onboard/main.go
./main
```

