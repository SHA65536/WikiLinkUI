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
   --link value, -l value   Address for LinkAPI (default: "localhost:2048")
   --redis value, -r value  Address for Redis (default: "localhost:6379")
   --help, -h               show help
```