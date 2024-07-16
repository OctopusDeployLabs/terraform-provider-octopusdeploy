package test

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

type ProjectGroupTestOptions struct {
	TestOptions[projectgroups.ProjectGroup]
}

func NewProjectGroupTestOptions() *ProjectGroupTestOptions {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	options := &ProjectGroupTestOptions{
		TestOptions: *NewTestOptions[projectgroups.ProjectGroup]("project_group"),
	}
	options.Resource = projectgroups.NewProjectGroup(name)
	options.Resource.Description = description

	return options
}

func ProjectGroupConfiguration(options *ProjectGroupTestOptions) string {
	configuration := fmt.Sprintf(`resource "%s" "%s" {`, options.ResourceName, options.LocalName)
	configuration += "\n"

	if len(options.Resource.Description) > 0 {
		configuration += fmt.Sprintf(`description = "%s"`, options.Resource.Description)
		configuration += "\n"
	}

	configuration += fmt.Sprintf(`name = "%s"`, options.Resource.Name)
	configuration += "\n"

	configuration += "}"
	return configuration
}
