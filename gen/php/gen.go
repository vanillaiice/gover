package gen

import genJSON "github.com/vanillaiice/gover/v3/gen/js"

func UpdateComposerVersion(filePath string, version string) ([]byte, error) {
	return genJSON.UpdatePackageVersion(filePath, version)
}
