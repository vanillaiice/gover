# gover [![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/github.com/vanillaiice/gover) [![Go Report Card](https://goreportcard.com/badge/github.com/vanillaiice/gover)](https://goreportcard.com/report/github.com/vanillaiice/gover)

gover is package version management tool for Go projects.

Instead of manually incrementing the version number in your code (🗿),
you can simply use `gover` to automatically do it.
Also, you can use `gover` to tag the git branch to the current version of your project (using `git tag`).

Under the hood, `gover` will read a `gover.json` file in your project
and update the version number accordingly.
It will also generate a `version.go` file that contains the current version number,
which you can import in your project.

# Installation

```sh
$ go install github.com/vanillaiice/gover@latest
```

# Usage

In your Go project, use `gover init` to initialize the version file:

```sh
# create a gover.json & version.go with a default version of v0.0.1
$ gover init
# create a version.json with a version of v1.0.0
$ gover init -v v1.0.0 -f version.json
# create a gover.json with a version of v1.0.0, custom path for the Go output file,
# and a custom package name
$ gover init -v v1.0.0 -o cmd/version.go -P cmd
```

Then, you can use the `version` constant in your project:

```go
package cmd

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// Exec starts the cli app.
func Exec() {
	app := &cli.App{
		Name:                   "gover",
		Usage:                  "package version management tool for Go",
		Version:                version,
```

> here, the `version.go` file is in the `cmd` directory.

Finally, you can use `gover bump` to increment the version number:

```
# bump to major version (e.g. v1.0.0 -> v2.0.0)
$ gover bump --major
# bump to minor version (e.g. v1.0.0 -> v1.1.0) with verbose log and custom output file
$ gover -V bump --minor -o ver.go
# bump to patch version (e.g. v1.0.0 -> v1.0.1) with custom package name
$ gover bump --patch -P pkg
```

On a side note, you can also use environment variables to define the package name, version file, and output file:

```sh
VERSION_FILE=gover.json
OUTPUT_FILE=cmd/version.go
PACKAGE_NAME=cmd
```

> these values can be defined in a file named `.env` or `.gover`.

# Help

```sh
NAME:
   gover - package version management tool for Go

USAGE:
   gover [global options] command [command options]

VERSION:
   1.0.0

AUTHOR:
   vanillaiice <vanillaiice1@proton.me>

COMMANDS:
   init, i  initialize a new version file
   gen, g   generate go version file from json version file
   bump, b  bump version
   tag, t   tag git branch with the current version
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, -V  show verbose log (default: false)
   --help, -h     show help
   --version, -v  print the version
```

# License

GPLv3

# Author

[vanillaiice](https://github.com/vanillaiice)
