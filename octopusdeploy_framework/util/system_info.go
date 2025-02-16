package util

import (
	"os"
	"path/filepath"
)

type SystemInfo struct {
	TerraformVersion string
	OctopusVersion   string
}

func (s *SystemInfo) GetExecutableName() string {
	executable, err := os.Executable()
	if err != nil {
		// We don't want to cause more errors while trying to add context to an error,
		// so we return "Unknown" if we can't get the executable name.
		return "Unknown"
	}
	return filepath.Base(executable)
}
