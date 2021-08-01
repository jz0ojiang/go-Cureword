package cli

import (
	"log"
	"os"

	// "github.com/0ojixueseno0/go-cureword/mods"
	"github.com/0ojixueseno0/go-cureword/mods"
	"github.com/urfave/cli/v2"
)

func Run() {
	app := cli.NewApp()
	// information
	app.Name = "CureWord"
	app.Usage = "API program for cureword.top"
	app.Version = "0.2.0"

	// Flags
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   256,
			Usage:   "listening port",
		},
		&cli.BoolFlag{
			Name:    "faststart",
			Aliases: []string{"fs"},
			Value:   false,
			Usage:   "Faststart Program",
		},
	}

	//commands
	app.Commands = []*cli.Command{
		{
			Name:  "account",
			Usage: "API account operations",
			// Category: "account",
			Subcommands: []*cli.Command{
				{
					Name:  "add",
					Usage: "add a new account",
					Action: func(c *cli.Context) error {
						return AddAccount(c)
					},
				},
				{
					Name:  "list",
					Usage: "list all accounts",
					Action: func(c *cli.Context) error {
						return ListAccount(c)
					},
				},
				{
					Name:  "delete",
					Usage: "delete an account",
					Action: func(c *cli.Context) error {
						return DeleteAccount(c)
					},
				},
				{
					Name:  "set",
					Usage: "set an account",
					Action: func(c *cli.Context) error {
						return SetAccount(c)
					},
				},
			},
		},
		{
			Name:  "run",
			Usage: "Run API serve",
			Action: func(c *cli.Context) error {
				return mods.Run(c)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
