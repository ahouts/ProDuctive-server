package main

import (
	"os"
	"time"

	"github.com/urfave/cli"
)

func runCli() {
	app := cli.NewApp()
	t := time.Now().Format(time.RFC3339)
	app.Version = "pre-release-" + t
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Andrew Houts",
			Email: "ahouts@scu.edu",
		},
		{
			Name:  "Joe Dion",
			Email: "jdion@scu.edu",
		},
	}
	app.Usage = "ProDuctive Server"

	var configFile string
	var webPort int

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "c, config",
			Value:       "./config.json",
			Usage:       "configuration `file` to load",
			Destination: &configFile,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "serve connections",
			Action: func(c *cli.Context) error {
				serve(configFile, webPort)
				return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "p, port",
					Usage:       "`port` to serve on",
					Destination: &webPort,
					Value:       444,
				},
			},
		},
		{
			Name:  "drop",
			Usage: "drop the database",
			Action: func(c *cli.Context) error {
				dropDb(configFile)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
