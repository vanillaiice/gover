package load

import (
	"fmt"
	"os"
	"strings"
)

func FromFile(file string) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	version, _, _, err := packageVersionRange(content)
	if err != nil {
		return "", fmt.Errorf("%s: %w", file, err)
	}

	return version, nil
}

func FromFilePanic(file string) string {
	v, err := FromFile(file)
	if err != nil {
		panic(err)
	}
	return v
}

func packageVersionRange(content []byte) (string, int, int, error) {
	inPackage := false
	offset := 0
	for _, line := range strings.SplitAfter(string(content), "\n") {
		lineBody := strings.TrimSuffix(line, "\n")
		lineBody = strings.TrimSuffix(lineBody, "\r")
		trimmed := strings.TrimSpace(lineBody)

		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			inPackage = trimmed == "[package]"
			offset += len(line)
			continue
		}

		if inPackage && !strings.HasPrefix(trimmed, "#") {
			value, start, end, ok := tomlStringValueRange(lineBody, "version")
			if ok {
				return value, offset + start, offset + end, nil
			}
		}

		offset += len(line)
	}

	return "", 0, 0, fmt.Errorf("package version not found")
}

func PackageVersionRangeForUpdate(content []byte) (string, int, int, error) {
	return packageVersionRange(content)
}

func tomlStringValueRange(line, key string) (string, int, int, bool) {
	trimmedLeft := strings.TrimLeft(line, " \t")
	keyStart := len(line) - len(trimmedLeft)
	if !strings.HasPrefix(trimmedLeft, key) {
		return "", 0, 0, false
	}

	i := keyStart + len(key)
	for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
		i++
	}
	if i >= len(line) || line[i] != '=' {
		return "", 0, 0, false
	}

	i++
	for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
		i++
	}
	if i >= len(line) || line[i] != '"' {
		return "", 0, 0, false
	}

	start := i
	escaped := false
	for i++; i < len(line); i++ {
		switch {
		case escaped:
			escaped = false
		case line[i] == '\\':
			escaped = true
		case line[i] == '"':
			return line[start+1 : i], start, i + 1, true
		}
	}

	return "", 0, 0, false
}
