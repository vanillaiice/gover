package cmd

import (
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/v3/load"
)

// commitCmdTemplateData is the template data for the commit command.
type commitCmdTemplateData struct {
	File    string
	Version string
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
			Value:   "version/version.go",
			EnvVars: []string{"VERSION_FILE"},
		},
		&cli.StringFlag{
			Name:    "command",
			Aliases: []string{"c"},
			Usage:   "template for commit `COMMAND`",
			Value:   "git commit {{ .File }} -m \"chore: bump version to {{ .Version }}\"",
			EnvVars: []string{"COMMIT_COMMAND"},
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		versionData, err := load.FromFile(ctx.Path("file"))
		if err != nil {
			return err
		}

		command, err := generateCommitCommand(ctx.String("command"), ctx.Path("file"), versionData.Version)
		if err != nil {
			return err
		}

		if err = runCommand(command); err != nil {
			return err
		}

		return nil
	},
}
