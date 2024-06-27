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
			EnvVars: []string{"VERSION_FILE"},
		},
		&cli.PathFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "write version to `FILE`",
			Value:   "version/version.go",
			EnvVars: []string{"OUTPUT_FILE"},
		},
		&cli.StringFlag{
			Name:    "package",
			Aliases: []string{"P"},
			Usage:   "set package name to `PACKAGE`",
			Value:   "version",
			EnvVars: []string{"PACKAGE_NAME"},
		},
		&cli.BoolFlag{
			Name:    "local",
			Aliases: []string{"l"},
			Usage:   "make the version constant local",
			Value:   false,
		},
	},
	Action: func(ctx *cli.Context) error {
		versionData, err := load.FromFile(ctx.Path("file"))
		if err != nil {
			return err
		}

		if err = gen.VersionFile(ctx.String("package"), "v"+versionData.Version, ctx.Bool("local"), ctx.Path("output")); err != nil {
			return err
		}

		if ctx.Bool("verbose") {
			log.Printf("generated %s", ctx.Path("output"))
		}

		return nil
	},
}
