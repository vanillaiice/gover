package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"
)

// tmpl is the template for the Go version file.
const tmpl = `package {{ .PackageName }}

// {{if .Local }}version{{ else }}Version{{ end }} is the current version of the package.
{{ if .Local }}const version = "{{ .Version }}"{{ else }}const Version = "{{ .Version }}"{{ end }}`

// TemplateData	is the template data.
type TemplateData struct {
	PackageName string
	Version     string
	Local       bool
}

// VersionFile generates a file containing the package version.
func VersionFile(packageName, version string, local bool) ([]byte, error) {
	tmpl, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, TemplateData{
		PackageName: packageName,
		Version:     version,
		Local:       local,
	}); err != nil {
		return []byte{}, err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return []byte{}, fmt.Errorf("generated invalid Go source: %w", err)
	}

	return formatted, nil
}
