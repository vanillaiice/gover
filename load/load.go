package load

import (
	"fmt"

	"github.com/vanillaiice/gover/v3/lang"
	load "github.com/vanillaiice/gover/v3/load/go"
	loadJS "github.com/vanillaiice/gover/v3/load/js"
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
	}

	if err != nil {
		panic(err)
	}

	return version
}
