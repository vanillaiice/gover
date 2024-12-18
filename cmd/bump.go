package cmd

import (
	"fmt"
	"log"

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
			Value:   "version/version.go",
			EnvVars: []string{"VERSION_FILE"},
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
			EnvVars: []string{"LOCAL_VERSION"},
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

		switch {
		case ctx.Bool("major"):
			*version = version.IncMajor()
		case ctx.Bool("minor"):
			*version = version.IncMinor()
		case ctx.Bool("patch"):
			*version = version.IncPatch()
		default:
			return fmt.Errorf("no version bump specified")
		}

		versionData.Version = "v" + version.String()
		if err = gen.VersionFile(ctx.String("package"), versionData.Version, ctx.Bool("local"), ctx.Path("file")); err != nil {
			return
		}

		if ctx.Bool("verbose") {
			log.Printf("bumped version to %s", versionData.Version)
		}

		return
	},
}
