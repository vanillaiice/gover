package load

import (
	"encoding/json"
	"os"
)

type packageData struct {
	Version string `json:"version"`
}

func FromFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	var versionData packageData
	if err = json.Unmarshal(data, &versionData); err != nil {
		return "", err
	}

	return versionData.Version, nil
}

func FromFilePanic(file string) string {
	v, err := FromFile(file)
	if err != nil {
		panic(err)
	}
	return v
}
