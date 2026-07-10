package cmd

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/v3/lang"
	"github.com/vanillaiice/gover/v3/load"
)

// commitCmdTemplateData is the template data for the commit command.
type commitCmdTemplateData struct {
	File    string
	Files   []string
	Version string
}

type commandResult struct {
	File    string    `json:"file,omitempty"`
	Lang    lang.Lang `json:"lang"`
	Version string    `json:"version"`
	Command string    `json:"command"`
	DryRun  bool      `json:"dry_run"`
}

// generateCommitCommand generates the commit command from the template.
func generateCommitCommand(tmpl, file, version string) (string, error) {
	template, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	if err = template.Execute(&b, commitCmdTemplateData{
		File:    file,
		Files:   []string{file},
		Version: version,
	}); err != nil {
		return "", err
	}

	return b.String(), nil
}

// commitCmd is the commit command.
var commitCmd = &cli.Command{
	Name:    "commit",
	Aliases: []string{"c"},
	Usage:   "commit version",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "load version from `FILE`",
			EnvVars: []string{"GOVER_VERSION_FILE"},
		},
		&cli.StringFlag{
			Name:    "command",
			Aliases: []string{"c"},
			Usage:   "template for commit `COMMAND`",
			Value:   "git commit {{ .File }} -m \"chore: bump version to {{ .Version }}\"",
			EnvVars: []string{"COMMIT_COMMAND"},
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "show the commit command without running it",
		},
		&cli.BoolFlag{
			Name:  "json",
			Usage: "print machine-readable JSON output",
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

		version, err := load.FromFile(file, l)
		if err != nil {
			return err
		}

		command, err := generateCommitCommand(ctx.String("command"), file, version)
		if err != nil {
			return err
		}

		if ctx.Bool("verbose") {
			log.Printf("running: %s", command)
		}

		result := commandResult{
			File:    file,
			Lang:    l,
			Version: version,
			Command: command,
			DryRun:  ctx.Bool("dry-run"),
		}
		if ctx.Bool("dry-run") {
			if ctx.Bool("json") {
				return printJSON(result)
			}
			fmt.Println(command)
			return nil
		}

		if err = runCommand(command); err != nil {
			return err
		}

		if ctx.Bool("json") {
			return printJSON(result)
		}

		return nil
	},
}
