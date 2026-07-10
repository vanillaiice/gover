package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/gover/v3/gen"
	"github.com/vanillaiice/gover/v3/lang"
	"github.com/vanillaiice/gover/v3/load"
)

type bumpTarget struct {
	file string
	lang lang.Lang
	root string
}

type bumpResult struct {
	File       string    `json:"file"`
	Lang       lang.Lang `json:"lang"`
	OldVersion string    `json:"old_version"`
	NewVersion string    `json:"new_version"`
	DryRun     bool      `json:"dry_run"`
	ConfigDir  string    `json:"-"`
}

// bumpCmd is the bump command.
var bumpCmd = &cli.Command{
	Name:    "bump",
	Usage:   "bump version",
	Aliases: []string{"b"},
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "load version from `FILE`",
			EnvVars: []string{"GOVER_VERSION_FILE"},
		},
		&cli.StringFlag{
			Name:    "package",
			Aliases: []string{"P"},
			Usage:   "set package name to `PACKAGE`",
			Value:   "version",
			EnvVars: []string{"GOVER_PACKAGE_NAME"},
		},
		&cli.BoolFlag{
			Name:    "local",
			Aliases: []string{"l"},
			Usage:   "make the version constant local",
			Value:   false,
			EnvVars: []string{"GOVER_LOCAL_VERSION"},
		},
		&cli.BoolFlag{
			Name:    "major",
			Aliases: []string{"m"},
			Usage:   "bump major version",
		},
		&cli.BoolFlag{
			Name:    "minor",
			Aliases: []string{"n"},
			Usage:   "bump minor version",
		},
		&cli.BoolFlag{
			Name:    "patch",
			Aliases: []string{"p"},
			Usage:   "bump patch version",
		},
		&cli.StringFlag{
			Name:  "set",
			Usage: "set exact version to `VERSION`",
		},
		&cli.BoolFlag{
			Name:    "recursive",
			Aliases: []string{"r"},
			Usage:   "find and bump supported version files below the current directory",
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "show the version change without writing files",
		},
		&cli.BoolFlag{
			Name:  "json",
			Usage: "print machine-readable JSON output",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		results, err := bumpSelection(ctx)
		if err != nil {
			return err
		}

		return outputBumpResults(ctx, results)
	},
}

func bumpSelection(ctx *cli.Context) ([]bumpResult, error) {
	if ctx.Bool("recursive") {
		return bumpRecursive(ctx)
	}
	if ctx.Args().Len() > 0 {
		return bumpPaths(ctx)
	}

	l := lang.Lang(ctx.String("lang"))

	file := ctx.Path("file")
	if file == "" {
		var err error
		file, err = lang.DefaultVersionFilePath(l)
		if err != nil {
			return nil, err
		}
	}

	result, err := bumpFile(ctx, bumpTarget{
		file: file,
		lang: l,
		root: targetConfigDirForFile(file),
	})
	if err != nil {
		return nil, err
	}
	return []bumpResult{result}, nil
}

func bumpRecursive(ctx *cli.Context) ([]bumpResult, error) {
	if ctx.Path("file") != "" {
		return nil, fmt.Errorf("--recursive cannot be used with --file")
	}

	roots := []string{"."}
	if ctx.Args().Len() > 0 {
		roots = ctx.Args().Slice()
	}

	filterLang, filter := explicitLang(ctx)
	var targets []bumpTarget
	for _, root := range roots {
		rootTargets, err := discoverBumpTargets(root, filterLang, filter)
		if err != nil {
			return nil, err
		}
		targets = append(targets, rootTargets...)
	}
	if len(targets) == 0 {
		return nil, fmt.Errorf("no supported version files found below current directory")
	}

	targets = uniqueBumpTargets(targets)
	results := make([]bumpResult, 0, len(targets))
	for _, target := range targets {
		result, err := bumpFile(ctx, target)
		if err != nil {
			return nil, fmt.Errorf("bump %s: %w", target.file, err)
		}
		results = append(results, result)
	}

	return results, nil
}

func bumpPaths(ctx *cli.Context) ([]bumpResult, error) {
	if ctx.Path("file") != "" {
		return nil, fmt.Errorf("path arguments cannot be used with --file")
	}

	filterLang, filter := explicitLang(ctx)
	var targets []bumpTarget
	for i := 0; i < ctx.Args().Len(); i++ {
		pathTargets, err := discoverBumpPathTargets(ctx.Args().Get(i), filterLang, filter)
		if err != nil {
			return nil, err
		}
		targets = append(targets, pathTargets...)
	}
	if len(targets) == 0 {
		return nil, fmt.Errorf("no supported version files found")
	}

	targets = uniqueBumpTargets(targets)
	results := make([]bumpResult, 0, len(targets))
	for _, target := range targets {
		result, err := bumpFile(ctx, target)
		if err != nil {
			return nil, fmt.Errorf("bump %s: %w", target.file, err)
		}
		results = append(results, result)
	}

	return results, nil
}

func bumpFile(ctx *cli.Context, target bumpTarget) (bumpResult, error) {
	file := target.file
	l := target.lang

	versionStr, err := load.FromFile(file, l)
	if err != nil {
		return bumpResult{}, err
	}

	if ctx.Bool("verbose") {
		log.Printf("loaded version %s from %s", versionStr, file)
	}

	version, err := semver.NewVersion(versionStr)
	if err != nil {
		return bumpResult{}, err
	}

	bumpCount := 0
	for _, flagSet := range []bool{ctx.Bool("major"), ctx.Bool("minor"), ctx.Bool("patch"), ctx.String("set") != ""} {
		if flagSet {
			bumpCount++
		}
	}
	if bumpCount != 1 {
		return bumpResult{}, fmt.Errorf("specify exactly one of --major, --minor, --patch, or --set")
	}

	var newVersion semver.Version
	switch {
	case ctx.Bool("major"):
		newVersion = version.IncMajor()
	case ctx.Bool("minor"):
		newVersion = version.IncMinor()
	case ctx.Bool("patch"):
		newVersion = version.IncPatch()
	case ctx.String("set") != "":
		setVersion, err := semver.NewVersion(ctx.String("set"))
		if err != nil {
			return bumpResult{}, err
		}
		newVersion = *setVersion
	default:
		return bumpResult{}, fmt.Errorf("no version bump specified")
	}

	newVersionStr := formatVersionForLang(l, newVersion)
	result := bumpResult{
		File:       file,
		Lang:       l,
		OldVersion: versionStr,
		NewVersion: newVersionStr,
		DryRun:     ctx.Bool("dry-run"),
		ConfigDir:  targetConfigDir(target),
	}
	if ctx.Bool("dry-run") {
		if ctx.Bool("verbose") {
			log.Printf("would bump %s from %s to %s", file, versionStr, newVersionStr)
		}
		return result, nil
	}

	var genOpts gen.Opts
	switch l {
	case lang.Go:
		packageName, err := targetStringValue(ctx, result.ConfigDir, "package", []string{"GOVER_PACKAGE_NAME"})
		if err != nil {
			return bumpResult{}, err
		}
		local, err := targetBoolValue(ctx, result.ConfigDir, "local", []string{"GOVER_LOCAL_VERSION"})
		if err != nil {
			return bumpResult{}, err
		}
		genOpts = gen.Opts{
			PackageName: packageName,
			Local:       local,
			Version:     newVersionStr,
		}
	case lang.JS, lang.TS, lang.Rust, lang.PHP:
		genOpts = gen.Opts{
			Version:    newVersionStr,
			OutputFile: file,
		}
	default:
		return bumpResult{}, fmt.Errorf("bump command not supported for lang %q", l)
	}

	out, err := gen.Version(l, &genOpts)
	if err != nil {
		return bumpResult{}, err
	}

	if err := os.WriteFile(file, out, 0644); err != nil {
		return bumpResult{}, err
	}

	if ctx.Bool("verbose") {
		log.Printf("bumped version to %s & generated %s\n", genOpts.Version, file)
	}

	return result, nil
}

func outputBumpResults(ctx *cli.Context, results []bumpResult) error {
	if ctx.Bool("json") {
		return printJSON(results)
	}

	if ctx.Bool("dry-run") {
		for _, result := range results {
			fmt.Printf("%s (%s): %s -> %s\n", result.File, result.Lang, result.OldVersion, result.NewVersion)
		}
	}

	return nil
}

func formatVersionForLang(l lang.Lang, version semver.Version) string {
	if l == lang.Go {
		return "v" + version.String()
	}
	return version.String()
}

func discoverBumpPathTargets(path string, filterLang lang.Lang, filter bool) ([]bumpTarget, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		target, ok := bumpTargetForFile(path)
		if !ok {
			return nil, fmt.Errorf("unsupported version file %q", path)
		}
		if filter && target.lang != filterLang {
			return nil, fmt.Errorf("version file %q is %s, not %s", path, target.lang, filterLang)
		}
		return []bumpTarget{target}, nil
	}

	var targets []bumpTarget
	configTarget, ok, err := configuredTargetForDir(path)
	if err != nil {
		return nil, err
	}
	if ok && (!filter || configTarget.lang == filterLang) {
		targets = append(targets, configTarget)
	}
	for _, candidate := range []string{
		filepath.Join(path, "version", "version.go"),
		filepath.Join(path, "package.json"),
		filepath.Join(path, "Cargo.toml"),
		filepath.Join(path, "composer.json"),
	} {
		target, ok := bumpTargetForFile(candidate)
		if !ok {
			continue
		}
		if _, err := os.Stat(candidate); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		if filter && target.lang != filterLang {
			continue
		}
		targets = append(targets, target)
	}

	if len(targets) == 0 {
		if filter {
			return nil, fmt.Errorf("no %s version file found in %s", filterLang, path)
		}
		return nil, fmt.Errorf("no supported version file found in %s", path)
	}

	return targets, nil
}

func discoverBumpTargets(root string, filterLang lang.Lang, filter bool) ([]bumpTarget, error) {
	var targets []bumpTarget
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != root && shouldSkipBumpDir(d.Name()) {
				return filepath.SkipDir
			}
			target, ok, err := configuredTargetForDir(path)
			if err != nil {
				return err
			}
			if ok && (!filter || target.lang == filterLang) {
				targets = append(targets, target)
			}
			return nil
		}

		target, ok := bumpTargetForFile(path)
		if ok && (!filter || target.lang == filterLang) {
			targets = append(targets, target)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	sortBumpTargets(targets)
	return targets, nil
}

func bumpTargetForFile(path string) (bumpTarget, bool) {
	if filepath.Base(path) == "package.json" {
		return bumpTarget{file: path, lang: lang.JS, root: targetConfigDirForFile(path)}, true
	}

	if filepath.Base(path) == "Cargo.toml" {
		return bumpTarget{file: path, lang: lang.Rust, root: targetConfigDirForFile(path)}, true
	}

	if filepath.Base(path) == "composer.json" {
		return bumpTarget{file: path, lang: lang.PHP, root: targetConfigDirForFile(path)}, true
	}

	if filepath.Base(path) == "version.go" && filepath.Base(filepath.Dir(path)) == "version" {
		return bumpTarget{file: path, lang: lang.Go, root: targetConfigDirForFile(path)}, true
	}

	return bumpTarget{}, false
}

func configuredTargetForDir(dir string) (bumpTarget, bool, error) {
	values, err := readPackageEnv(dir)
	if err != nil {
		return bumpTarget{}, false, err
	}

	versionFile, ok := values["GOVER_VERSION_FILE"]
	if !ok || versionFile == "" {
		return bumpTarget{}, false, nil
	}

	file := versionFile
	if !filepath.IsAbs(file) {
		file = filepath.Join(dir, file)
	}
	if _, err := os.Stat(file); err != nil {
		return bumpTarget{}, false, err
	}

	if langName, ok := values["GOVER_LANG"]; ok && langName != "" {
		l, err := lang.ParseLang(langName)
		if err != nil {
			return bumpTarget{}, false, err
		}
		return bumpTarget{file: file, lang: l, root: dir}, true, nil
	}

	target, ok := bumpTargetForFile(file)
	if !ok {
		return bumpTarget{}, false, fmt.Errorf("%s configures %s, but GOVER_LANG is required for custom version files", dir, versionFile)
	}
	target.root = dir
	return target, true, nil
}

func targetConfigDir(target bumpTarget) string {
	if target.root != "" {
		return target.root
	}
	return targetConfigDirForFile(target.file)
}

func targetConfigDirForFile(file string) string {
	if filepath.Base(file) == "version.go" && filepath.Base(filepath.Dir(file)) == "version" {
		return filepath.Dir(filepath.Dir(file))
	}
	return filepath.Dir(file)
}

func sortBumpTargets(targets []bumpTarget) {
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].file < targets[j].file
	})
}

func uniqueBumpTargets(targets []bumpTarget) []bumpTarget {
	sortBumpTargets(targets)

	unique := targets[:0]
	var previous string
	for i, target := range targets {
		if i > 0 && target.file == previous {
			continue
		}
		unique = append(unique, target)
		previous = target.file
	}

	return unique
}

func explicitLang(ctx *cli.Context) (lang.Lang, bool) {
	if !ctx.IsSet("lang") {
		return "", false
	}

	return lang.Lang(ctx.String("lang")), true
}

func shouldSkipBumpDir(name string) bool {
	switch name {
	case ".git", "node_modules", "vendor", "dist", "build", "coverage", "target":
		return true
	default:
		return false
	}
}
