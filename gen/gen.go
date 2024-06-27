package gen

import (
	"os"
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
func VersionFile(packageName, version string, local bool, outputFile string) (err error) {
	tmpl, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return
	}
	defer f.Close()

	return tmpl.Execute(f, TemplateData{
		PackageName: packageName,
		Version:     version,
		Local:       local,
	})
}
