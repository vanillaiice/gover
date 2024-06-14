package cmd

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// Exec starts the cli app.
func Exec() {
	app := &cli.App{
		Name:                   "gover",
		Usage:                  "package version management tool for Go projects",
		Version:                version,
		Suggest:                true,
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Authors:                []*cli.Author{{Name: "vanillaiice", Email: "vanillaiice1@proton.me"}},
		Commands: []*cli.Command{
			initCmd,
			genCmd,
			bumpCmd,
			tagCmd,
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"V"},
				Usage:   "show verbose log",
				Value:   false,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
