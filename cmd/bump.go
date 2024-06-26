package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/gen"
	"github.com/vanillaiice/gover/load"
)

// perm is the file permission.
const perm = 0644

// bumpCmd is the bump command.
var bumpCmd = &cli.Command{
	Name:    "bump",
	Usage:   "bump version",
	Aliases: []string{"b"},
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
		&cli.BoolFlag{
			Name:    "major",
			Aliases: []string{"m"},
			Usage:   "bump major version",
		},
		&cli.BoolFlag{
			Name:    "minor",
			Aliases: []string{"n"},
			Usage:   "bump minor version",
		},
		&cli.BoolFlag{
			Name:    "patch",
			Aliases: []string{"p"},
			Usage:   "bump patch version",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		versionData, err := load.FromFile(ctx.Path("file"))
		if err != nil {
			return
		}

		version, err := semver.NewVersion(versionData.Version)
		if err != nil {
			return
		}

		if ctx.Bool("major") {
			*version = version.IncMajor()
		} else if ctx.Bool("minor") {
			*version = version.IncMinor()
		} else if ctx.Bool("patch") {
			*version = version.IncPatch()
		} else {
			return fmt.Errorf("no version bump specified")
		}

		versionData.Version = "v" + version.String()
		data, err := json.MarshalIndent(versionData, "", "  ")
		if err != nil {
			return
		}

		if err = os.WriteFile(ctx.Path("file"), data, perm); err != nil {
			return
		}

		if err = gen.VersionFile(ctx.String("package"), version.String(), ctx.Bool("local"), ctx.Path("output")); err != nil {
			return
		}

		if ctx.Bool("verbose") {
			log.Printf("bumped version to %s", versionData.Version)
		}

		return
	},
}
