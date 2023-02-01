# WikiLinkUI
This is a repository for the web UI for the [WikiLinkAPI](https://github.com/SHA65536/WikiLinkApi)

## Running
To run the UI run `go run ./cmd.ui`:
```
NAME:
   ui - Serves the WikiLink web ui

USAGE:
   ui [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value   Port to listen to (default: "3000")
   --api value, -a value    Address for LinkAPI (default: "localhost:2048")
   --redis value, -r value  Address for Redis (default: "localhost:6379")
   --log value, -l value    Level of log to be shown ("trace", "debug", "info", "warn", "error", "fatal", "panic") (default: "info")
   --help, -h               show help
```