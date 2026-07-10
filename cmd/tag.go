package cmd

import (
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/v3/lang"
	"github.com/vanillaiice/gover/v3/load"
)

// tagCmdTemplateData	is the template data for the tag command.
type tagCmdTemplateData struct {
	Version string
}

// generateTagCommand generates the tag command from the template.
func generateTagCommand(tmpl, version string) (string, error) {
	template, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	if err = template.Execute(&b, tagCmdTemplateData{
		Version: version,
	}); err != nil {
		return "", err
	}

	return b.String(), nil
}

// tagCmd is the tag command.
var tagCmd = &cli.Command{
	Name:    "tag",
	Aliases: []string{"t"},
	Usage:   "tag branch with the current version",
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
			Usage:   "template for tag `COMMAND`",
			Value:   "git tag {{ .Version }}",
			EnvVars: []string{"GOVER_TAG_COMMAND"},
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "show the tag command without running it",
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
			return
		}

		command, err := generateTagCommand(ctx.String("command"), version)
		if err != nil {
			return
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
			return
		}

		if ctx.Bool("json") {
			return printJSON(result)
		}

		return nil
	},
}
