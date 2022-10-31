package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandVersionControlSettingsForProjectConversion(ctx context.Context, d *schema.ResourceData) projects.GitPersistenceSettings {

	var persistenceSettings projects.GitPersistenceSettings
	if v, ok := d.GetOk("git_library_persistence_settings"); ok {
		persistenceSettings = expandGitPersistenceSettings(ctx, v, expandLibraryGitCredential)
	}
	if v, ok := d.GetOk("git_username_password_persistence_settings"); ok {
		persistenceSettings = expandGitPersistenceSettings(ctx, v, expandUsernamePasswordGitCredential)
	}
	if v, ok := d.GetOk("git_anonymous_persistence_settings"); ok {
		persistenceSettings = expandGitPersistenceSettings(ctx, v, expandAnonymousGitCredential)
	}

	return persistenceSettings
}
