# gover [![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/github.com/vanillaiice/gover) [![Go Report Card](https://goreportcard.com/badge/github.com/vanillaiice/gover)](https://goreportcard.com/report/github.com/vanillaiice/gover)

gover is package version management tool for Go, JS/TS, Rust, and PHP projects.

Instead of manually incrementing the version number in your code like a - 🗿,
you can simply use `gover` to automatically do it.
Also, you can use `gover` to commit changes to your go version file,
and tag branches.

Under the hood, `gover` will read a version file (e.g. `version.go`, `package.json`, `Cargo.toml`, `composer.json`) in your project
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
	"fmt"
	"slices"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/v3/version"
)

// Exec starts the cli app.
func Exec() {
	app := &cli.App{
		Name:                   "gover",
		Usage:                  "package version management tool for Go and JS projects",
		Version:                version.Version,
```

> here, the `version.go` file is in the `version` directory.

## `bump`

You can increment the package version using the `bump` command:

```sh
# bump to major version (e.g. v1.0.0 -> v2.0.0)
$ gover bump --major
# bump to minor version (e.g. v1.0.0 -> v1.1.0) with verbose log and custom Go version file
$ gover -V bump --minor -f ver.go
# bump to patch version (e.g. v1.0.0 -> v1.0.1) with custom package name
$ gover bump --patch -P pkg
# bump to patch version (e.g. v1.0.0 -> v1.0.1) with custom package name for js
$ gover --lang js bump --patch
# set an exact version
$ gover bump --set v2.0.0
# preview a bump without writing files
$ gover bump --dry-run --patch
# emit machine-readable bump results
$ gover bump --json --patch
```

### Monorepos

From a monorepo root, pass one or more package paths to bump only those packages:

```sh
# bump one subrepo from a monorepo root
$ gover bump --patch apps/web
# bump multiple selected subrepos from a monorepo root
$ gover bump --patch services/api apps/web
```

Use `--recursive` when you want `gover` to search below the current directory or a selected subtree:

```sh
# bump all supported version files in a monorepo from the root
$ gover bump --recursive --patch
# bump all supported version files below one subtree
$ gover bump --recursive --patch apps
```

Use `--lang` to filter recursive or path-targeted bumps:

```sh
# bump only Go version files in a monorepo from the root
$ gover --lang go bump --recursive --patch
# bump only JS package files below one subtree
$ gover --lang js bump --recursive --patch apps
```

When a package path is passed without `--recursive`, `gover` checks that directory for:

- Go: `version/version.go`
- JS/TS: `package.json`
- Rust: `Cargo.toml`
- PHP: `composer.json`

When `--recursive` is used, `gover` searches below the selected root and skips `.git`, `node_modules`, `vendor`, `dist`, `build`, `coverage`, and `target`.

For recursive and path-targeted runs, `gover` also reads `.env` and `.gover` in each package directory. Command-line flags take precedence, then package-local config, then root environment/default values.

Package-local config can customize generated Go files:

```sh
# services/api/.gover
GOVER_PACKAGE_NAME=api
GOVER_LOCAL_VERSION=true
```

It can also point `gover` at a custom version file:

```sh
# services/api/.gover
GOVER_LANG=go
GOVER_VERSION_FILE=internal/version.go
GOVER_PACKAGE_NAME=internal
```

For `release --recursive`, package-local release templates are supported:

```sh
# services/api/.gover
COMMIT_COMMAND=git commit {{ .File }} -m "release {{ .Version }}"
GOVER_TAG_COMMAND=git tag api/{{ .Version }}
```

## `release`

You can run the common bump, commit, and tag workflow with `release`:

```sh
# bump patch, commit the version file, and tag the new version
$ gover release --patch
# preview the release without writing files or running git commands
$ gover release --dry-run --patch
# customize commit/tag commands
$ gover release --minor \
  --commit-command "git commit {{ .File }} -m 'release {{ .Version }}'" \
  --tag-command "git tag {{ .Version }}"
```

Use `--push` to run the configured push command after committing and tagging.

## `commit`

You can commit the Go version file using the `commit` command:

```sh
# commit using default git template
$ gover commit
# commit with custom template and Go version file
$ gover commit -f cmd/ver.go --command "git commit {{ .File }} -m 'bump to {{ .Version }}'"
# preview the rendered command
$ gover commit --dry-run
```

## `tag`

You can tag the current branch using the `tag` command:

```sh
# tag using default git command ("git tag {{ .Version }}")
$ gover tag
# preview the rendered command
$ gover tag --dry-run
```

## `get`

The `get` commands returns the current version of the package:

```sh
$ gover get
# with custom file
$ gover get -f cmd/ver.go
# JSON output
$ gover get --json
```

## `check`

The `check` command validates that version files exist and contain valid semantic versions:

```sh
$ gover check
$ gover check --recursive
$ gover check --json --recursive
```

> You can also set some arguments with environment variables:

> - lang: GOVER_LANG
> - version file: GOVER_VERSION_FILE
> - package name: GOVER_PACKAGE_NAME
> - commit command: COMMIT_COMMAND
> - tag command: GOVER_TAG_COMMAND

> these values can be defined in a file named `.env` or `.gover`.

# Help

```sh
NAME:
   gover - package version management tool for Go, JS/TS, Rust, and PHP projects

USAGE:
   gover [global options] command [command options]

VERSION:
   v3.4.0

AUTHOR:
   vanillaiice <vanillaiice1@proton.me>

COMMANDS:
   init, i    initialize a new version file
   bump, b    bump version
   release, r bump, commit, and tag a release
   commit, c  commit version
   tag, t     tag branch with the current version
   get, e     get the current version
   check, k   validate version files
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, -V         show verbose log (default: false)
   --lang LANG, -l LANG  use language LANG (default: "go") [$GOVER_LANG]
   --help, -h            show help
   --version, -v         print the version
```

# License

GPLv3

# Author

[vanillaiice](https://github.com/vanillaiice)
