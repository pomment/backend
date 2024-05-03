package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var commands []*cli.Command

func init() {
	commands = []*cli.Command{
		{
			Name:     "start",
			Usage:    "Run the Pomment RESTful API service",
			Category: "server",
			Action: func(c *cli.Context) error {
				err := StartStandaloneServer(c.String(""))
				return err
			},
		},
	}
}

func main() {
	app := &cli.App{
		Name:     "Pomment",
		Usage:    "Pomment backend service",
		Commands: commands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
