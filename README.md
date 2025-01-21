# OnboardingBot

## How to run
- Install Go
- Create file `.env` in `secrets/` directory
- put `.json` files (with creds for connection with google or MM clients) into `secrets/`
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
- in `config/const.go` paste absolute path to `.json` files - it will ignore different exceptions. 
- name your files for connections in secrets with specific names (for instance: `googlecalendar-creds.json`)
- for example, `googlecalendar-creds.json` push into secrets and use creds.

