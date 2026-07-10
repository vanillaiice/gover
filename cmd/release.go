package cmd

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

type releaseResult struct {
	Bump   bumpResult     `json:"bump"`
	Commit commandResult  `json:"commit"`
	Tag    commandResult  `json:"tag"`
	Push   *commandResult `json:"push,omitempty"`
}

var releaseCmd = &cli.Command{
	Name:    "release",
	Aliases: []string{"r"},
	Usage:   "bump, commit, and tag a release",
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
		&cli.StringFlag{
			Name:  "set",
			Usage: "set exact version to `VERSION`",
		},
		&cli.BoolFlag{
			Name:    "recursive",
			Aliases: []string{"R"},
			Usage:   "find and release supported version files below the current directory",
		},
		&cli.StringFlag{
			Name:    "commit-command",
			Usage:   "template for commit `COMMAND`",
			Value:   "git commit {{ .File }} -m \"chore: bump version to {{ .Version }}\"",
			EnvVars: []string{"COMMIT_COMMAND"},
		},
		&cli.StringFlag{
			Name:    "tag-command",
			Usage:   "template for tag `COMMAND`",
			Value:   "git tag {{ .Version }}",
			EnvVars: []string{"GOVER_TAG_COMMAND"},
		},
		&cli.BoolFlag{
			Name:  "push",
			Usage: "run the push command after committing and tagging",
		},
		&cli.StringFlag{
			Name:    "push-command",
			Usage:   "push `COMMAND`",
			Value:   "git push --tags",
			EnvVars: []string{"GOVER_PUSH_COMMAND"},
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "show release actions without writing files or running commands",
		},
		&cli.BoolFlag{
			Name:  "json",
			Usage: "print machine-readable JSON output",
		},
	},
	Action: func(ctx *cli.Context) error {
		results, err := bumpSelection(ctx)
		if err != nil {
			return err
		}

		releaseResults := make([]releaseResult, 0, len(results))
		for _, bump := range results {
			commitTemplate, err := targetStringValue(ctx, bump.ConfigDir, "commit-command", []string{"COMMIT_COMMAND"})
			if err != nil {
				return err
			}
			tagTemplate, err := targetStringValue(ctx, bump.ConfigDir, "tag-command", []string{"GOVER_TAG_COMMAND"})
			if err != nil {
				return err
			}

			commitCommand, err := generateCommitCommand(commitTemplate, bump.File, bump.NewVersion)
			if err != nil {
				return err
			}
			tagCommand, err := generateTagCommand(tagTemplate, bump.NewVersion)
			if err != nil {
				return err
			}

			commit := commandResult{
				File:    bump.File,
				Lang:    bump.Lang,
				Version: bump.NewVersion,
				Command: commitCommand,
				DryRun:  ctx.Bool("dry-run"),
			}
			tag := commandResult{
				File:    bump.File,
				Lang:    bump.Lang,
				Version: bump.NewVersion,
				Command: tagCommand,
				DryRun:  ctx.Bool("dry-run"),
			}

			releaseResults = append(releaseResults, releaseResult{
				Bump:   bump,
				Commit: commit,
				Tag:    tag,
			})
		}

		if ctx.Bool("dry-run") {
			var push *commandResult
			if ctx.Bool("push") {
				push = &commandResult{
					Command: ctx.String("push-command"),
					DryRun:  true,
				}
			}
			return outputReleaseResults(ctx, releaseResults, push)
		}

		for _, result := range releaseResults {
			if ctx.Bool("verbose") {
				log.Printf("running: %s", result.Commit.Command)
			}
			if err := runCommand(result.Commit.Command); err != nil {
				return err
			}

			if ctx.Bool("verbose") {
				log.Printf("running: %s", result.Tag.Command)
			}
			if err := runCommand(result.Tag.Command); err != nil {
				return err
			}
		}

		var push *commandResult
		if ctx.Bool("push") {
			push = &commandResult{
				Command: ctx.String("push-command"),
				DryRun:  false,
			}
			if ctx.Bool("verbose") {
				log.Printf("running: %s", push.Command)
			}
			if err := runCommand(push.Command); err != nil {
				return err
			}
		}

		return outputReleaseResults(ctx, releaseResults, push)
	},
}

func outputReleaseResults(ctx *cli.Context, results []releaseResult, push *commandResult) error {
	if push != nil {
		for i := range results {
			results[i].Push = push
		}
	}

	if ctx.Bool("json") {
		return printJSON(results)
	}

	if ctx.Bool("dry-run") {
		for _, result := range results {
			fmt.Printf("%s (%s): %s -> %s\n", result.Bump.File, result.Bump.Lang, result.Bump.OldVersion, result.Bump.NewVersion)
			fmt.Printf("commit: %s\n", result.Commit.Command)
			fmt.Printf("tag: %s\n", result.Tag.Command)
		}
		if ctx.Bool("push") {
			fmt.Printf("push: %s\n", ctx.String("push-command"))
		}
	}

	return nil
}
