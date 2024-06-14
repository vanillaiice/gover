package cmd

import (
	"log"

	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/gen"
	"github.com/vanillaiice/gover/load"
)

// genCmd is the gen command.
var genCmd = &cli.Command{
	Name:    "gen",
	Aliases: []string{"g"},
	Usage:   "generate go version file from json version file",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "load version from `FILE`",
			Value:   "gover.json",
		},
		&cli.PathFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "write version to `FILE`",
			Value:   "version.go",
		},
		&cli.StringFlag{
			Name:    "package",
			Aliases: []string{"P"},
			Usage:   "set package name to `PACKAGE`",
			Value:   "main",
		},
	},
	Action: func(ctx *cli.Context) error {
		versionData, err := load.FromFile(ctx.Path("file"))
		if err != nil {
			return err
		}

		if err = gen.VersionFile(ctx.String("package"), versionData.Version, ctx.Path("output")); err != nil {
			return err
		}

		if ctx.Bool("verbose") {
			log.Printf("generated %s", ctx.Path("output"))
		}

		return nil
	},
}
