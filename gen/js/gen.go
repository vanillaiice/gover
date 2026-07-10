package gen

import (
	"fmt"
	"os"
	"strconv"
)

// UpdatePackageVersion updates package.json with the new version.
func UpdatePackageVersion(filePath string, version string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading file: %w", err)
	}

	valueStart, valueEnd, err := topLevelStringFieldRange(content, "version")
	if err != nil {
		return []byte{}, fmt.Errorf("no version field found in %s", filePath)
	}

	quotedVersion := strconv.Quote(version)
	updatedContent := make([]byte, 0, len(content)-valueEnd+valueStart+len(quotedVersion))
	updatedContent = append(updatedContent, content[:valueStart]...)
	updatedContent = append(updatedContent, quotedVersion...)
	updatedContent = append(updatedContent, content[valueEnd:]...)

	return updatedContent, nil
}

func topLevelStringFieldRange(content []byte, field string) (int, int, error) {
	i := skipWhitespace(content, 0)
	if i >= len(content) || content[i] != '{' {
		return 0, 0, fmt.Errorf("expected JSON object")
	}

	depth := 0
	expectKey := false
	for i < len(content) {
		switch content[i] {
		case ' ', '\n', '\r', '\t':
			i++
		case '{':
			depth++
			if depth == 1 {
				expectKey = true
			}
			i++
		case '}':
			if depth == 1 {
				expectKey = false
			}
			depth--
			i++
		case '[', ']':
			if content[i] == '[' {
				depth++
			} else {
				depth--
			}
			i++
		case ',':
			if depth == 1 {
				expectKey = true
			}
			i++
		case ':':
			if depth == 1 {
				expectKey = false
			}
			i++
		case '"':
			stringEnd, err := jsonStringEnd(content, i)
			if err != nil {
				return 0, 0, err
			}

			if depth == 1 && expectKey {
				key, err := strconv.Unquote(string(content[i:stringEnd]))
				if err != nil {
					return 0, 0, err
				}

				colon := skipWhitespace(content, stringEnd)
				if colon >= len(content) || content[colon] != ':' {
					return 0, 0, fmt.Errorf("expected colon after key")
				}

				valueStart := skipWhitespace(content, colon+1)
				if key == field {
					if valueStart >= len(content) || content[valueStart] != '"' {
						return 0, 0, fmt.Errorf("field %q is not a string", field)
					}

					valueEnd, err := jsonStringEnd(content, valueStart)
					if err != nil {
						return 0, 0, err
					}

					return valueStart, valueEnd, nil
				}
			}

			i = stringEnd
		default:
			i++
		}
	}

	return 0, 0, fmt.Errorf("field %q not found", field)
}

func skipWhitespace(content []byte, start int) int {
	for start < len(content) {
		switch content[start] {
		case ' ', '\n', '\r', '\t':
			start++
		default:
			return start
		}
	}

	return start
}

func jsonStringEnd(content []byte, start int) (int, error) {
	escaped := false
	for i := start + 1; i < len(content); i++ {
		switch {
		case escaped:
			escaped = false
		case content[i] == '\\':
			escaped = true
		case content[i] == '"':
			return i + 1, nil
		}
	}

	return 0, fmt.Errorf("unterminated JSON string")
}
