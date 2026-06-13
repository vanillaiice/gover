package cmd

import (
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/v3/version"
)

// Exec starts the cli app.
func Exec(arguments []string) error {
	app := &cli.App{
		Name:                   "gover",
		Usage:                  "package version management tool for Go projects",
		Version:                version.Version,
		Suggest:                true,
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Authors:                []*cli.Author{{Name: "vanillaiice", Email: "vanillaiice1@proton.me"}},
		Commands: []*cli.Command{
			initCmd,
			bumpCmd,
			commitCmd,
			tagCmd,
			getCmd,
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"V"},
				Usage:   "show verbose log",
				Value:   false,
			},
			/*
				&cli.PathFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "load version from `FILE`",
					Value:   "version/version.go",
					EnvVars: []string{"VERSION_FILE"},
				},
			*/
		},
	}

	_ = godotenv.Load(".env")
	_ = godotenv.Load(".gover")

	return app.Run(arguments)
}
