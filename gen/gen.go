package gen

import (
	"fmt"

	gen "github.com/vanillaiice/gover/v3/gen/go"
	genJS "github.com/vanillaiice/gover/v3/gen/js"
	genPHP "github.com/vanillaiice/gover/v3/gen/php"
	genRust "github.com/vanillaiice/gover/v3/gen/rust"
	"github.com/vanillaiice/gover/v3/lang"
)

type Opts struct {
	OutputFile  string
	PackageName string
	Version     string
	Local       bool
}

// Version generates a version file.
func Version(l lang.Lang, opts *Opts) ([]byte, error) {
	var (
		generated []byte
		err       error
	)

	switch l {
	case lang.Go:
		generated, err = gen.VersionFile(
			opts.PackageName,
			opts.Version,
			opts.Local,
		)
	case lang.JS, lang.TS:
		generated, err = genJS.UpdatePackageVersion(
			opts.OutputFile,
			opts.Version,
		)
	case lang.Rust:
		generated, err = genRust.UpdateCargoVersion(
			opts.OutputFile,
			opts.Version,
		)
	case lang.PHP:
		generated, err = genPHP.UpdateComposerVersion(
			opts.OutputFile,
			opts.Version,
		)
	default:
		return []byte{}, fmt.Errorf("unsupported lang %q", l)
	}

	return generated, err
}
