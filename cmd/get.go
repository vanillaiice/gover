package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/v3/lang"
	"github.com/vanillaiice/gover/v3/load"
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
			EnvVars: []string{"GOVER_VERSION_FILE"},
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

		fmt.Printf("%s\n", version)

		return nil
	},
}
