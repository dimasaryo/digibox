package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "digibox"
	app.Usage = "Digitalocean development box cli"
	app.Version = "0.0.1"

	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			err := errors.New("Command is empty. Please see --help for available commands.")
			return err
		}

		if c.Args().Get(0) == "start" {
			log.Printf("Start development box")
			return nil
		} else if c.Args().Get(0) == "stop" {
			log.Printf("Stop development box")
			return nil
		} else {
			err := errors.New("Unknown command.")
			return err
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
