package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
)

func expandGitPersistenceSettings(ctx context.Context, values interface{}, callback func(ctx context.Context, flattenedMap map[string]interface{}) credentials.GitCredential) projects.GitPersistenceSettings {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 || flattenedValues[0] == nil {
		return nil
	}

	flattenedMap := flattenedValues[0].(map[string]interface{})

	tflog.Info(ctx, "expanding Git credentials")

	gitUrl, err := url.Parse(flattenedMap["url"].(string))
	if err != nil {
		return nil
	}

	gitCredential := callback(ctx, flattenedMap)

	return projects.NewGitPersistenceSettings(
		flattenedMap["base_path"].(string),
		gitCredential,
		flattenedMap["default_branch"].(string),
		[]string{},
		gitUrl,
	)
}

func expandLibraryGitCredential(ctx context.Context, flattenedMap map[string]interface{}) credentials.GitCredential {
	tflog.Info(ctx, "expanding reference credential")
	return credentials.NewReference(flattenedMap["git_credential_id"].(string))
}

func expandUsernamePasswordGitCredential(ctx context.Context, flattenedMap map[string]interface{}) credentials.GitCredential {
	tflog.Info(ctx, "expanding U/P credential")
	return credentials.NewUsernamePassword(
		flattenedMap["username"].(string),
		core.NewSensitiveValue(flattenedMap["password"].(string)),
	)
}

func expandAnonymousGitCredential(ctx context.Context, flattenedMap map[string]interface{}) credentials.GitCredential {
	tflog.Info(ctx, "expanding Anonymous credential")
	return credentials.NewAnonymous()
}

func flattenGitPersistenceSettings(ctx context.Context, persistenceSettings projects.PersistenceSettings) []interface{} {
	if persistenceSettings == nil || persistenceSettings.Type() == projects.PersistenceSettingsTypeDatabase {
		return nil
	}

	gitPersistenceSettings := persistenceSettings.(projects.GitPersistenceSettings)

	flattenedGitPersistenceSettings := make(map[string]interface{})
	flattenedGitPersistenceSettings["base_path"] = gitPersistenceSettings.BasePath()
	flattenedGitPersistenceSettings["default_branch"] = gitPersistenceSettings.DefaultBranch()

	credential := gitPersistenceSettings.Credential()
	switch credential.Type() {
	case credentials.GitCredentialTypeReference:
		tflog.Info(ctx, "flatten reference credential")
		flattenedGitPersistenceSettings["git_credential_id"] = credential.(*credentials.Reference).Id
	case credentials.GitCredentialTypeUsernamePassword:
		tflog.Info(ctx, "flatten U/P credential")
		flattenedGitPersistenceSettings["username"] = credential.(*credentials.UsernamePassword).Username
	}

	if gitPersistenceSettings.URL() != nil {
		flattenedGitPersistenceSettings["url"] = gitPersistenceSettings.URL().String()
	}

	return []interface{}{flattenedGitPersistenceSettings}
}
