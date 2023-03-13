package main

import (
	"log"
	"os"

	"github.com/SHA65536/wikilinkui"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func main() {
	var config_path string
	var (
		// UI_PORT
		port string = "3000"
		// API_ADDR
		apiAddr string = "localhost:2048"
		// REDIS_ADDR
		rAddr string = "localhost:6379"
		// VAULT_ADDR
		vAddr string = ""
		// VAULT_ROLE
		vRole string = ""
		// LOG_LEVEL
		logLevel string = "info"
		// LOG_PATH
		logPath string = "/var/wikilinkui/wikilinkui.log"
	)

	app := &cli.App{
		Name:  "wikilinkui",
		Usage: "Serves the WikiLink web ui",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config_path",
				Aliases:     []string{"c"},
				Value:       "/var/wikilinkui/.env",
				Usage:       "Path to config file",
				Destination: &config_path,
			},
		},
		Action: func(ctx *cli.Context) error {
			// Loading environment variables
			if err := godotenv.Load(config_path); err != nil {
				return err
			}
			if val, ok := os.LookupEnv("UI_PORT"); ok {
				port = val
			}
			if val, ok := os.LookupEnv("API_ADDR"); ok {
				apiAddr = val
			}
			if val, ok := os.LookupEnv("REDIS_ADDR"); ok {
				rAddr = val
			}
			if val, ok := os.LookupEnv("VAULT_ADDR"); ok {
				vAddr = val
			}
			if val, ok := os.LookupEnv("VAULT_ROLE"); ok {
				vRole = val
			}
			if val, ok := os.LookupEnv("LOG_LEVEL"); ok {
				logLevel = val
			}
			if val, ok := os.LookupEnv("LOG_PATH"); ok {
				logPath = val
			}
			var level, err = zerolog.ParseLevel(logLevel)
			if err != nil {
				return err
			}
			logf, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			api, err := wikilinkui.MakeUIHandler("heb", apiAddr, rAddr, vAddr, vRole, level, logf)
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
