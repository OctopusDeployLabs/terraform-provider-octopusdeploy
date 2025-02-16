package logging

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"os"
	"path/filepath"
)

// AddDiagnosticError is used to wrap calls to Diagnostics.AddError with additional information about the executable and versions
func AddDiagnosticError(resp *resource.ReadResponse, config *octopusdeploy_framework.Config, message string, err error) {
	suffix := "\nPlease ensure these details are included in any error report you raise.\n" +
		"Executable: " + getExecutableName() + "\n" +
		"Terraform Version: " + config.TerraformVersion + "\n" +
		"Octopus Version: " + config.OctopusVersion

	resp.Diagnostics.AddError(message+suffix, err.Error())
}

func getExecutableName() string {
	executable, err := os.Executable()
	if err != nil {
		// We don't want to cause more errors while trying to add context to an error,
		// so we return "Unknown" if we can't get the executable name.
		return "Unknown"
	}
	return filepath.Base(executable)
}
