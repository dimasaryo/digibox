package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	tokenSource := &TokenSource{
		AccessToken: os.Getenv("DIGITALOCEAN_TOKEN"),
	}

	doClient := NewDigitalOceanClient(tokenSource)
	app := cli.NewApp()
	app.Name = "digibox"
	app.Usage = "Digitalocean remote development server cli"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start a remote development server `NAME`",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					err := errors.New("Missing devbox name")
					return err
				}

				log.Printf("Start remote development server")
				err := doClient.Start(c.Args().Get(0))
				if err != nil {
					log.Fatal(err)
					return err
				}
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "stop remote development server `NAME`",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					err := errors.New("Missing name")
					return err
				}

				log.Printf("Stop remote development server")
				err := doClient.Stop()
				if err != nil {
					log.Fatal(err)
					return err
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
