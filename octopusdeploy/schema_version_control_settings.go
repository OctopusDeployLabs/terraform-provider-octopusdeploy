package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandVersionControlSettings(versionControlSettings interface{}) *octopusdeploy.VersionControlSettings {
	versionControlSettingsList := versionControlSettings.(*schema.Set).List()
	versionControlSettingsMap := versionControlSettingsList[0].(map[string]interface{})

	return &octopusdeploy.VersionControlSettings{
		DefaultBranch: versionControlSettingsMap["default_branch"].(string),
		Password:      octopusdeploy.NewSensitiveValue(versionControlSettingsMap["password"].(string)),
		URL:           versionControlSettingsMap["url"].(string),
		Username:      versionControlSettingsMap["username"].(string),
	}
}

func flattenVersionControlSettings(versionControlSettings *octopusdeploy.VersionControlSettings) []interface{} {
	if versionControlSettings == nil {
		return nil
	}

	flattenedVersionControlSettings := make(map[string]interface{})
	flattenedVersionControlSettings["default_branch"] = versionControlSettings.DefaultBranch
	flattenedVersionControlSettings["url"] = versionControlSettings.URL
	flattenedVersionControlSettings["username"] = versionControlSettings.Username
	return []interface{}{flattenedVersionControlSettings}
}

func getVersionControlSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_branch": {
			Description: "The default branch associated with these version control settings.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"password": {
			Description:      "The password associated with these version control settings.",
			Sensitive:        true,
			Optional:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.StringIsNotEmpty),
		},
		"url": {
			Description: "The URL associated with these version control settings.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"username": {
			Description:      "The username associated with these version control settings.",
			Optional:         true,
			Sensitive:        true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.StringIsNotEmpty),
		},
	}
}
