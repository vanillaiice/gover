# gover [![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/github.com/vanillaiice/gover) [![Go Report Card](https://goreportcard.com/badge/github.com/vanillaiice/gover)](https://goreportcard.com/report/github.com/vanillaiice/gover)

gover is package version management tool for Go projects.

Instead of manually incrementing the version number in your code like a - 🗿,
you can simply use `gover` to automatically do it.
Also, you can use `gover` to commit changes to your go version file,
and tag branches.

Under the hood, `gover` will read a Go version file (e.g. `version.go`) in your project
and update the version number accordingly.

# Installation

```sh
$ go install github.com/vanillaiice/gover/v3@latest
```

# Usage

In your Go project, use `gover init` to initialize the version file:

```sh
# create a Go version file with a default version of v0.0.1
$ gover init
# create a Go version file with a version of v1.0.0
$ gover init -v v1.0.0
# create a custom Go version file with a version of v1.0.0,
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

## `bump`

You can increment the verison on the Go version file using the `bump` command:

```sh
# bump to major version (e.g. v1.0.0 -> v2.0.0)
$ gover bump --major
# bump to minor version (e.g. v1.0.0 -> v1.1.0) with verbose log and custom Go version file
$ gover -V bump --minor -f ver.go
# bump to patch version (e.g. v1.0.0 -> v1.0.1) with custom package name
$ gover bump --patch -P pkg
```

## `commit`

You can commit the Go version file using the `commit` command:

```sh
# commit using default git template
$ gover commit
# commit with custom template and Go version file
$ gover commit -f cmd/ver.go --command "git commit {{ .File }} -m 'bump to {{ .Version }}'"
```

## `tag`

You can tag the current branch using the `tag` command:

```sh
# tag using default git command ("git tag {{ .Version }}")
$ gover tag
```

## `get`

The `get` commands returns the current version of the package:

```sh
$ gover get
# with custom file
$ gover get -f cmd/ver.go
```

> You can also set some arguments with environment variables:

> - version file: VERSION_FILE
> - package name: PACKAGE_NAME
> - commit command: COMMIT_COMMAND
> - tag command: TAG_COMMAND

> these values can be defined in a file named `.env` or `.gover`.

# Help

```sh
NAME:
   gover - package version management tool for Go projects

USAGE:
   gover [global options] command [command options]

VERSION:
   v3.0.0

AUTHOR:
   vanillaiice <vanillaiice1@proton.me>

COMMANDS:
   init, i    initialize a new version file
   bump, b    bump version
   commit, c  commit version
   tag, t     tag branch with the current version
   get, e     get the current version
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, -V  show verbose log (default: false)
   --help, -h     show help
   --version, -v  print the version
```

# Related Projects

- [gover-js](https://github.com/vanillaiice/gover-js), a package version management tool for JavaScript projects.

# License

GPLv3

# Author

[vanillaiice](https://github.com/vanillaiice)
