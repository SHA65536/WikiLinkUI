package main

import (
	"log"
	"os"

	"github.com/SHA65536/wikilinkui"
	"github.com/urfave/cli/v2"
)

func main() {
	var port string
	var linkapi string
	var redis string

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
				Name:        "link",
				Aliases:     []string{"l"},
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
		},
		Action: func(ctx *cli.Context) error {
			api, err := wikilinkui.MakeUIHandler("heb", linkapi, redis)
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
