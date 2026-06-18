package lang

import "fmt"

type Lang string

const (
	Go Lang = "go"
	JS Lang = "js"
	TS Lang = "ts"
)

func ParseLang(lang string) (Lang, error) {
	l := Lang(lang)

	switch l {
	case Go, JS, TS:
	default:
		return "", fmt.Errorf("invalid lang %q", lang)
	}

	return l, nil
}

func DefaultVersionFilePath(lang Lang) (file string, err error) {
	switch lang {
	case Go:
		file = "version/version.go"
	case JS, TS:
		file = "package.json"
	default:
		err = fmt.Errorf("invalid lang %q", lang)
	}

	return
}
