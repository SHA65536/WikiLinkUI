package main

import (
	"log"
	"os"

	"github.com/SHA65536/linkapiui"
	"github.com/urfave/cli/v2"
)

func main() {
	var port string

	app := &cli.App{
		Name:        "ui",
		Description: "Serves the linkapi ui",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       "3000",
				Usage:       "Port to listen to",
				Destination: &port,
			},
		},
		Action: func(ctx *cli.Context) error {
			api, err := linkapiui.MakeUIHandler("heb", "localhost:2048")
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
