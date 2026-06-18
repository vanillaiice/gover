package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/v3/gen"
	"github.com/vanillaiice/gover/v3/lang"
	"github.com/vanillaiice/gover/v3/load"
)

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
			EnvVars: []string{"GOVER_VERSION_FILE"},
		},
		&cli.StringFlag{
			Name:    "package",
			Aliases: []string{"P"},
			Usage:   "set package name to `PACKAGE`",
			Value:   "version",
			EnvVars: []string{"GOVER_PACKAGE_NAME"},
		},
		&cli.BoolFlag{
			Name:    "local",
			Aliases: []string{"l"},
			Usage:   "make the version constant local",
			Value:   false,
			EnvVars: []string{"GOVER_LOCAL_VERSION"},
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
		l := lang.Lang(ctx.String("lang"))

		file := ctx.Path("file")
		if file == "" {
			file, err = lang.DefaultVersionFilePath(l)
			if err != nil {
				return err
			}
		}

		versionStr, err := load.FromFile(ctx.Path("file"), l)
		if err != nil {
			return
		}

		if ctx.Bool("verbose") {
			log.Printf("loaded version %s", versionStr)
		}

		version, err := semver.NewVersion(versionStr)
		if err != nil {
			return
		}

		var newVersion semver.Version
		switch {
		case ctx.Bool("major"):
			newVersion = version.IncMajor()
		case ctx.Bool("minor"):
			newVersion = version.IncMinor()
		case ctx.Bool("patch"):
			newVersion = version.IncPatch()
		default:
			return fmt.Errorf("no version bump specified")
		}

		var genOpts gen.Opts
		switch l {
		case lang.Go:
			genOpts = gen.Opts{
				PackageName: ctx.String("package"),
				Local:       ctx.Bool("local"),
				Version:     "v" + newVersion.String(),
			}
		case lang.JS, lang.TS:
			genOpts = gen.Opts{
				Version:    versionStr,
				OutputFile: file,
			}
		default:
			return fmt.Errorf("bump command not supported for lang %q", l)
		}

		out, err := gen.Version(l, &genOpts)
		if err != nil {
			return
		}

		if err := os.WriteFile(ctx.Path("file"), out, 0644); err != nil {
			return err
		}

		if ctx.Bool("verbose") {
			log.Printf("bumped version to %s & generated %s\n", genOpts.Version, ctx.Path("file"))
		}

		return
	},
}
