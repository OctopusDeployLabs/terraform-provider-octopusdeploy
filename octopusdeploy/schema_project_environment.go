package octopusdeploy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandProjectEnvironments(value interface{}) map[string][]string {
	expandedProjectEnvironments := map[string][]string{}

	set := value.(*schema.Set)
	for _, item := range set.List() {
		projectEnvironment := item.(map[string]interface{})
		projectID := projectEnvironment["project_id"].(string)
		environments := []string{}
		for _, e := range projectEnvironment["environments"].([]interface{}) {
			environments = append(environments, e.(string))
		}

		expandedProjectEnvironments[projectID] = environments
	}

	return expandedProjectEnvironments
}

func flattenProjectEnvironments(projectEnvironments map[string][]string) []interface{} {
	if projectEnvironments == nil {
		return nil
	}

	flattenedProjectEnvironments := []interface{}{}
	for projectID, enviroments := range projectEnvironments {
		rawProjectEnvironment := map[string]interface{}{
			"project_id":   projectID,
			"environments": enviroments,
		}
		flattenedProjectEnvironments = append(flattenedProjectEnvironments, rawProjectEnvironment)
	}

	return flattenedProjectEnvironments
}
