package lang

import "fmt"

type Lang string

const (
	Plain      = "plain"
	Go    Lang = "go"
	JS    Lang = "js"
	TS    Lang = "ts"
	Rust  Lang = "rust"
	PHP   Lang = "php"
)

func ParseLang(lang string) (Lang, error) {
	l := Lang(lang)

	switch l {
	case Go, JS, TS, Rust, PHP, Plain:
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
	case Rust:
		file = "Cargo.toml"
	case PHP:
		file = "composer.json"
	case Plain:
		file = "version.txt"
	default:
		err = fmt.Errorf("invalid lang %q", lang)
	}

	return
}
