package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/Masterminds/semver/v3"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/v3/gen"
	"github.com/vanillaiice/gover/v3/lang"
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
			Usage:   "write version to `FILE`",
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
			Usage:   "make the version constant local (version instead of Version)",
			Value:   false,
			EnvVars: []string{"GOVER_LOCAL_VERSION"},
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"F"},
			Usage:   "overwrite the version file if it already exists",
			Value:   false,
		},
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "set version to `VERSION`",
			Value:   "0.0.1",
		},
	},
	Before: func(ctx *cli.Context) error {
		if lang := ctx.String("lang"); slices.Contains([]string{"js", "ts"}, lang) {
			return errors.New("init should not be called for js projects")
		}
		return nil
	},
	Action: func(ctx *cli.Context) (err error) {
		file := ctx.Path("file")

		l := lang.Lang(ctx.String("lang"))
		if file == "" {
			file, err = lang.DefaultVersionFilePath(l)
			if err != nil {
				return err
			}
		}

		if _, err := os.Stat(file); !errors.Is(err, os.ErrNotExist) {
			if err == nil {
				if !ctx.Bool("force") {
					return fmt.Errorf("file %s already exists", file)
				}
			} else {
				return err
			}
		}

		path := filepath.Dir(file)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}

		version, err := semver.NewVersion(ctx.String("version"))
		if err != nil {
			return err
		}

		var genOpts gen.Opts
		switch l {
		case lang.Go:
			versionStr := "v" + version.String()
			genOpts = gen.Opts{
				PackageName: ctx.String("package"),
				Local:       ctx.Bool("local"),
				Version:     versionStr,
			}
		default:
			return fmt.Errorf("init command not supported for lang %q", l)
		}

		out, err := gen.Version(l, &genOpts)
		if err != nil {
			return err
		}

		if err := os.WriteFile(file, out, 0644); err != nil {
			return err
		}

		if ctx.Bool("verbose") {
			fmt.Printf("created %s with version %s\n", file, genOpts.Version)
		}

		return nil
	},
}
