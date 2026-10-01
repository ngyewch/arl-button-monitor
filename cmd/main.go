package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

var (
	version string

	configFileFlag = &cli.StringFlag{
		Name:    "config-file",
		Usage:   "config file",
		Value:   "config.yml",
		Sources: cli.EnvVars("CONFIG_FILE"),
	}

	app = &cli.Command{
		Name:    "arl-button-monitor",
		Usage:   "ARL Button Monitor",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:   "serve",
				Action: doServe,
				Flags: []cli.Flag{
					configFileFlag,
				},
			},
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		slog.Error("error",
			slog.Any("err", err),
		)
		os.Exit(1)
	}
}
