package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/gen"
	"github.com/vanillaiice/gover/load"
)

// initCmd is the init command.
var initCmd = &cli.Command{
	Name:    "init",
	Usage:   "initialize a new version file",
	Aliases: []string{"i"},
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
			Usage:   "make the version constant local (version instead of Version)",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"F"},
			Usage:   "overwrite the json version file if it already exists",
			Value:   false,
		},
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "set version to `VERSION`",
			Value:   "v0.0.1",
		},
	},
	Action: func(ctx *cli.Context) error {
		if _, err := os.Stat(ctx.Path("file")); !errors.Is(err, os.ErrNotExist) {
			if err == nil {
				if !ctx.Bool("force") {
					return fmt.Errorf("file %s already exists", ctx.Path("file"))
				}
			} else {
				return err
			}
		}

		version, err := semver.NewVersion(ctx.String("version"))
		if err != nil {
			return err
		}

		versionData := load.VersionData{Version: "v" + version.String()}

		data, err := json.MarshalIndent(versionData, "", "  ")
		if err != nil {
			return err
		}

		if err = os.WriteFile(ctx.Path("file"), data, perm); err != nil {
			return err
		}

		path := filepath.Dir(ctx.Path("output"))
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}

		if err = gen.VersionFile(ctx.String("package"), version.String(), ctx.Bool("local"), ctx.Path("output")); err != nil {
			return err
		}

		if ctx.Bool("verbose") {
			fmt.Printf("wrote version to %s & generated %s\n", ctx.Path("file"), ctx.Path("output"))
		}

		return nil
	},
}
