package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandGitCredential(values interface{}) credentials.IGitCredential {
	if values == nil {
		return credentials.NewAnonymous()
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 || flattenedValues[0] == nil {
		return credentials.NewAnonymous()
	}

	flattenedMap := flattenedValues[0].(map[string]interface{})

	return credentials.NewUsernamePassword(
		flattenedMap["username"].(string),
		core.NewSensitiveValue(flattenedMap["password"].(string)),
	)
}

func flattenGitCredential(ctx context.Context, d *schema.ResourceData, gitCredential credentials.IGitCredential) []interface{} {
	if gitCredential == nil {
		return nil
	}

	switch gitCredential.GetType() {
	case "UsernamePassword":
		usernamePasswordCredential := gitCredential.(*credentials.UsernamePassword)
		return []interface{}{map[string]interface{}{
			"password": d.Get("git_persistence_settings.0.credentials.0.password").(string),
			"username": usernamePasswordCredential.Username,
		}}
	}

	return []interface{}{}
}
