package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenVersionControlSettings(versionControlSettings *octopusdeploy.VersionControlSettings) []interface{} {
	if versionControlSettings == nil {
		return nil
	}

	flattenedVersionControlSettings := make(map[string]interface{})
	flattenedVersionControlSettings["default_branch"] = versionControlSettings.DefaultBranch
	flattenedVersionControlSettings["password"] = versionControlSettings.Password
	flattenedVersionControlSettings["url"] = versionControlSettings.URL
	flattenedVersionControlSettings["username"] = versionControlSettings.Username
	return []interface{}{flattenedVersionControlSettings}
}

func getVersionControlSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_branch": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"password": {
			Computed:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"url": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"username": {
			Computed:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
	}
}
