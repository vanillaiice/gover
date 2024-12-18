package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/load"
)

var getCmd = &cli.Command{
	Name:    "get",
	Aliases: []string{"e"},
	Usage:   "get the current version",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "load version from `FILE`",
			Value:   "version/version.go",
			EnvVars: []string{"VERSION_FILE"},
		},
	},
	Action: func(ctx *cli.Context) error {
		versionData, err := load.FromFile(ctx.Path("file"))
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", versionData.Version)

		return nil
	},
}
