package cmd

import (
	"fmt"
	"slices"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/v3/version"
)

// Exec starts the cli app.
func Exec(arguments []string) error {
	explicitCommandFlags = collectExplicitCommandFlags(arguments)

	app := &cli.App{
		Name:                   "gover",
		Usage:                  "package version management tool for Go, JS/TS, Rust, and PHP projects",
		Version:                version.Version,
		Suggest:                true,
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Authors:                []*cli.Author{{Name: "vanillaiice", Email: "vanillaiice1@proton.me"}},
		Commands: []*cli.Command{
			initCmd,
			bumpCmd,
			releaseCmd,
			commitCmd,
			tagCmd,
			getCmd,
			checkCmd,
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"V"},
				Usage:   "show verbose log",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "lang",
				Aliases: []string{"l"},
				Usage:   "use language `LANG`",
				Value:   "go",
				EnvVars: []string{"GOVER_LANG"},
			},
		},
		Before: func(ctx *cli.Context) error {
			if lang := ctx.String("lang"); !slices.Contains([]string{"go", "js", "ts", "rust", "php"}, lang) {
				return fmt.Errorf("unsupported lang %q", lang)
			}
			return nil
		},
	}

	_ = godotenv.Load(".env")
	_ = godotenv.Load(".gover")

	return app.Run(arguments)
}
