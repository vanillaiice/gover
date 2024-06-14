package cmd

import (
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/load"
)

var tagCmd = &cli.Command{
	Name:    "tag",
	Aliases: []string{"t"},
	Usage:   "tag git branch with the current version",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "load version from `FILE`",
			Value:   "gover.json",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		versionData, err := load.FromFile(ctx.Path("file"))
		if err != nil {
			return
		}

		cmd := exec.Command("git", "tag", versionData.Version)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err = cmd.Run(); err != nil {
			return
		}

		return nil
	},
}
