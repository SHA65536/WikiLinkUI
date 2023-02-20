# WikiLinkUI
This is a repository for the web UI for the [WikiLinkAPI](https://github.com/SHA65536/WikiLinkApi)

## Running
To run the UI run `go run ./cmd/wikilinkui`:
```
NAME:
   wikilinkui - Serves the WikiLink web ui

USAGE:
   wikilinkui [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config_path value, -c value  Path to config file (default: "/var/wikilinkui/.env")
   --help, -h                     show help
```

## Config
The .env configuration file should have the following format:
```
UI_PORT="3000"
API_ADDR="localhost:2048"
REDIS_ADDR="localhost:6379"
LOG_LEVEL="info"
LOG_PATH="wikilinkui.log"
```