package cmd

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/v3/lang"
	"github.com/vanillaiice/gover/v3/load"
)

type checkResult struct {
	File    string    `json:"file"`
	Lang    lang.Lang `json:"lang"`
	Version string    `json:"version,omitempty"`
	OK      bool      `json:"ok"`
	Error   string    `json:"error,omitempty"`
}

var checkCmd = &cli.Command{
	Name:    "check",
	Aliases: []string{"k"},
	Usage:   "validate version files",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "load version from `FILE`",
			EnvVars: []string{"GOVER_VERSION_FILE"},
		},
		&cli.BoolFlag{
			Name:    "recursive",
			Aliases: []string{"r"},
			Usage:   "find and validate supported version files below the current directory",
		},
		&cli.BoolFlag{
			Name:  "json",
			Usage: "print machine-readable JSON output",
		},
	},
	Action: func(ctx *cli.Context) error {
		targets, err := checkTargets(ctx)
		if err != nil {
			return err
		}

		results := make([]checkResult, 0, len(targets))
		var failed bool
		for _, target := range targets {
			result := validateTarget(target)
			if !result.OK {
				failed = true
			}
			results = append(results, result)
		}

		if ctx.Bool("json") {
			if err := printJSON(results); err != nil {
				return err
			}
		} else {
			for _, result := range results {
				if result.OK {
					fmt.Printf("%s (%s): %s\n", result.File, result.Lang, result.Version)
				} else {
					fmt.Printf("%s (%s): %s\n", result.File, result.Lang, result.Error)
				}
			}
		}

		if failed {
			return fmt.Errorf("one or more version files are invalid")
		}
		return nil
	},
}

func checkTargets(ctx *cli.Context) ([]bumpTarget, error) {
	if ctx.Path("file") != "" {
		if ctx.Bool("recursive") || ctx.Args().Len() > 0 {
			return nil, fmt.Errorf("--file cannot be used with --recursive or path arguments")
		}
		return []bumpTarget{{
			file: ctx.Path("file"),
			lang: lang.Lang(ctx.String("lang")),
			root: targetConfigDirForFile(ctx.Path("file")),
		}}, nil
	}

	filterLang, filter := explicitLang(ctx)
	if ctx.Bool("recursive") {
		roots := []string{"."}
		if ctx.Args().Len() > 0 {
			roots = ctx.Args().Slice()
		}

		var targets []bumpTarget
		for _, root := range roots {
			rootTargets, err := discoverBumpTargets(root, filterLang, filter)
			if err != nil {
				return nil, err
			}
			targets = append(targets, rootTargets...)
		}
		targets = uniqueBumpTargets(targets)
		if len(targets) == 0 {
			return nil, fmt.Errorf("no supported version files found")
		}
		return targets, nil
	}

	if ctx.Args().Len() > 0 {
		var targets []bumpTarget
		for i := 0; i < ctx.Args().Len(); i++ {
			pathTargets, err := discoverBumpPathTargets(ctx.Args().Get(i), filterLang, filter)
			if err != nil {
				return nil, err
			}
			targets = append(targets, pathTargets...)
		}
		return uniqueBumpTargets(targets), nil
	}

	l := lang.Lang(ctx.String("lang"))
	file, err := lang.DefaultVersionFilePath(l)
	if err != nil {
		return nil, err
	}
	return []bumpTarget{{file: file, lang: l, root: targetConfigDirForFile(file)}}, nil
}

func validateTarget(target bumpTarget) checkResult {
	result := checkResult{
		File: target.file,
		Lang: target.lang,
	}

	version, err := load.FromFile(target.file, target.lang)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	result.Version = version

	if _, err := semver.NewVersion(version); err != nil {
		result.Error = err.Error()
		return result
	}

	result.OK = true
	return result
}
