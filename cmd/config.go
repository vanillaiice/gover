package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var explicitCommandFlags map[string]bool

func collectExplicitCommandFlags(arguments []string) map[string]bool {
	flags := make(map[string]bool)
	commandIndex := commandArgIndex(arguments)
	if commandIndex == -1 {
		return flags
	}

	shortFlags := map[string]string{
		"P": "package",
		"l": "local",
	}
	longFlags := map[string]string{
		"package":        "package",
		"local":          "local",
		"commit-command": "commit-command",
		"tag-command":    "tag-command",
		"push-command":   "push-command",
	}

	for _, arg := range arguments[commandIndex+1:] {
		if arg == "--" {
			break
		}
		if strings.HasPrefix(arg, "--") {
			name := strings.TrimPrefix(arg, "--")
			if i := strings.IndexByte(name, '='); i >= 0 {
				name = name[:i]
			}
			if canonical, ok := longFlags[name]; ok {
				flags[canonical] = true
			}
			continue
		}
		if strings.HasPrefix(arg, "-") && len(arg) == 2 {
			if canonical, ok := shortFlags[strings.TrimPrefix(arg, "-")]; ok {
				flags[canonical] = true
			}
		}
	}

	return flags
}

func commandArgIndex(arguments []string) int {
	for i := 1; i < len(arguments); i++ {
		arg := arguments[i]
		switch {
		case arg == "--":
			return -1
		case arg == "--verbose" || arg == "-V" || arg == "--help" || arg == "-h" || arg == "--version" || arg == "-v":
			continue
		case arg == "--lang" || arg == "-l":
			i++
			continue
		case strings.HasPrefix(arg, "--lang="):
			continue
		case strings.HasPrefix(arg, "-"):
			continue
		default:
			return i
		}
	}
	return -1
}

func cliFlagExplicit(name string) bool {
	return explicitCommandFlags[name]
}

func readPackageEnv(dir string) (map[string]string, error) {
	values := make(map[string]string)
	for _, name := range []string{".env", ".gover"} {
		file := filepath.Join(dir, name)
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		fileValues, err := godotenv.Read(file)
		if err != nil {
			return nil, err
		}
		for key, value := range fileValues {
			if _, exists := values[key]; !exists {
				values[key] = value
			}
		}
	}
	return values, nil
}

func targetStringValue(ctx *cli.Context, configDir, flag string, envVars []string) (string, error) {
	if cliFlagExplicit(flag) {
		return ctx.String(flag), nil
	}

	values, err := readPackageEnv(configDir)
	if err != nil {
		return "", err
	}
	for _, envVar := range envVars {
		if value, ok := values[envVar]; ok {
			return value, nil
		}
	}

	return ctx.String(flag), nil
}

func targetBoolValue(ctx *cli.Context, configDir, flag string, envVars []string) (bool, error) {
	if cliFlagExplicit(flag) {
		return ctx.Bool(flag), nil
	}

	values, err := readPackageEnv(configDir)
	if err != nil {
		return false, err
	}
	for _, envVar := range envVars {
		if value, ok := values[envVar]; ok {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return false, err
			}
			return parsed, nil
		}
	}

	return ctx.Bool(flag), nil
}
