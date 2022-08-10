package test

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

type ProjectCreateTestOptions struct {
	Lifecycle    *LifecycleTestOptions
	ProjectGroup *ProjectGroupTestOptions
}

type ProjectTestOptions struct {
	Lifecycle    *LifecycleTestOptions
	ProjectGroup *ProjectGroupTestOptions
	TestOptions[projects.Project]
}

func NewProjectTestOptions(lifecycle *LifecycleTestOptions, projectGroup *ProjectGroupTestOptions) *ProjectTestOptions {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	options := &ProjectTestOptions{
		Lifecycle:    lifecycle,
		ProjectGroup: projectGroup,
		TestOptions:  *NewTestOptions[projects.Project]("project"),
	}
	options.Resource = projects.NewProject(name, lifecycle.Resource.ID, projectGroup.Resource.ID)
	options.Resource.Description = description

	return options
}

func ProjectConfiguration(options *ProjectTestOptions) string {
	configuration := fmt.Sprintf(`resource "%s" "%s" {`, options.ResourceName, options.LocalName)
	configuration += "\n"

	if len(options.Resource.Description) > 0 {
		configuration += fmt.Sprintf(`description = "%s"`, options.Resource.Description)
		configuration += "\n"
	}

	configuration += fmt.Sprintf(`lifecycle_id = %s.id`, options.Lifecycle.QualifiedName)
	configuration += "\n"
	configuration += fmt.Sprintf(`name = "%s"`, options.Resource.Name)
	configuration += "\n"
	configuration += fmt.Sprintf(`project_group_id = %s.id`, options.ProjectGroup.QualifiedName)
	configuration += "\n"

	if options.Resource.IsDisabled {
		configuration += fmt.Sprintf(`is_disabled = %v`, options.Resource.IsDisabled)
		configuration += "\n"
	}

	configuration += "}"
	return configuration
}
