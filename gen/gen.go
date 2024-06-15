package gen

import (
	"os"
	"text/template"
)

// perm is the file permission used for writing.
const perm = 0644

// tmplFile is the template file.
const tmplFile = "version.tmpl"

// TemplateData	is the template data.
type TemplateData struct {
	PackageName string
	Version     string
}

// VersionFile generates a file containing the package version.
func VersionFile(packageName, version, outputFile string) (err error) {
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		return err
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, TemplateData{
		PackageName: packageName,
		Version:     version,
	})
}
