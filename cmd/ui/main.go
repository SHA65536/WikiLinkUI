package main

import (
	"log"
	"os"

	"github.com/SHA65536/wikilinkui"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func main() {
	var port string
	var linkapi string
	var redis string
	var loglevel string

	app := &cli.App{
		Name:  "ui",
		Usage: "Serves the WikiLink web ui",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       "3000",
				Usage:       "Port to listen to",
				Destination: &port,
			},
			&cli.StringFlag{
				Name:        "api",
				Aliases:     []string{"a"},
				Value:       "localhost:2048",
				Usage:       "Address for LinkAPI",
				Destination: &linkapi,
			},
			&cli.StringFlag{
				Name:        "redis",
				Aliases:     []string{"r"},
				Value:       "localhost:6379",
				Usage:       "Address for Redis",
				Destination: &redis,
			},
			&cli.StringFlag{
				Name:        "log",
				Aliases:     []string{"l"},
				Value:       "info",
				Usage:       `Level of log to be shown ("trace", "debug", "info", "warn", "error", "fatal", "panic")`,
				Destination: &loglevel,
			},
		},
		Action: func(ctx *cli.Context) error {
			var level, err = zerolog.ParseLevel(loglevel)
			if err != nil {
				return err
			}
			logf, err := os.OpenFile("linkapi.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			api, err := wikilinkui.MakeUIHandler("heb", linkapi, redis, level, logf)
			if err != nil {
				return err
			}
			err = api.Serve(":" + port)
			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
