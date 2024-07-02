package test

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
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
	options := &ProjectTestOptions{
		Lifecycle:    lifecycle,
		ProjectGroup: projectGroup,
		TestOptions:  *NewTestOptions[projects.Project]("project"),
	}

	// remove comments if testing against Git-backed persistence

	// basePath := ".octopus/" + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	// defaultBranch := "main"
	// url, _ := url.Parse("https://example.com")

	// password := core.NewSensitiveValue(acctest.RandStringFromCharSet(20, acctest.CharSetAlpha))
	// username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	// credentials := projects.NewUsernamePasswordGitCredential(username, password)

	project := projects.NewProject(acctest.RandStringFromCharSet(20, acctest.CharSetAlpha), lifecycle.Resource.ID, projectGroup.Resource.ID)
	// project.PersistenceSettings = projects.NewGitPersistenceSettings(basePath, credentials, defaultBranch, url)
	project.Description = acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	options.Resource = project

	return options
}

func ProjectConfiguration(options *ProjectTestOptions) string {
	configuration := fmt.Sprintf(`resource "%s" "%s" {`, options.ResourceName, options.LocalName) + "\n"

	if len(options.Resource.Description) > 0 {
		configuration += fmt.Sprintf(`description = "%s"`, options.Resource.Description) + "\n"
	}

	configuration += fmt.Sprintf(`lifecycle_id = %s.id`, options.Lifecycle.QualifiedName) + "\n"
	configuration += fmt.Sprintf(`name = "%s"`, options.Resource.Name) + "\n"
	configuration += fmt.Sprintf(`project_group_id = %s.id`, options.ProjectGroup.QualifiedName) + "\n"

	if options.Resource.IsDisabled {
		configuration += fmt.Sprintf(`is_disabled = %v`, options.Resource.IsDisabled) + "\n"
	}

	if options.Resource.PersistenceSettings != nil {
		if options.Resource.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
			gitPersistenceSettings := options.Resource.PersistenceSettings.(projects.GitPersistenceSettings)
			configuration += `git_persistence_settings {` + "\n"
			configuration += fmt.Sprintf(`base_path = "%s"`, gitPersistenceSettings.BasePath()) + "\n"
			configuration += fmt.Sprintf(`url = "%s"`, gitPersistenceSettings.URL().String()) + "\n"
			configuration += "}" + "\n"
		}
	}

	if options.Resource.SpaceID != "" {
		configuration += fmt.Sprintf(`space_id = "%s"`, options.Resource.SpaceID) + "\n"
	}

	configuration += "}"
	return configuration
}
