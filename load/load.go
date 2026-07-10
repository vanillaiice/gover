package load

import (
	"fmt"

	"github.com/vanillaiice/gover/v3/lang"
	load "github.com/vanillaiice/gover/v3/load/go"
	loadJS "github.com/vanillaiice/gover/v3/load/js"
	loadPHP "github.com/vanillaiice/gover/v3/load/php"
	loadRust "github.com/vanillaiice/gover/v3/load/rust"
)

// FromFile loads the version from the specified file.
func FromFile(file string, l lang.Lang) (string, error) {
	var (
		version string
		err     error
	)

	switch l {
	case lang.Go:
		version, err = load.FromFile(file)
	case lang.JS, lang.TS:
		version, err = loadJS.FromFile(file)
	case lang.Rust:
		version, err = loadRust.FromFile(file)
	case lang.PHP:
		version, err = loadPHP.FromFile(file)
	default:
		return "", fmt.Errorf("invalid lang %q", l)
	}

	return version, err
}

// FromFilePanic is the same as FromFile but it panics on error.
func FromFilePanic(file string, l lang.Lang) string {
	var (
		version string
		err     error
	)

	switch l {
	case lang.Go:
		version, err = load.FromFile(file)
	case lang.JS, lang.TS:
		version, err = loadJS.FromFile(file)
	case lang.Rust:
		version, err = loadRust.FromFile(file)
	case lang.PHP:
		version, err = loadPHP.FromFile(file)
	}

	if err != nil {
		panic(err)
	}

	return version
}
