package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandVersionControlSettings(flattenedVersionControlSettings interface{}) *octopusdeploy.VersionControlSettings {
	versionControlSettingsList := flattenedVersionControlSettings.(*schema.Set).List()
	versionControlSettingsMap := versionControlSettingsList[0].(map[string]interface{})

	versionControlSettings := &octopusdeploy.VersionControlSettings{
		BasePath:      versionControlSettingsMap["base_path"].(string),
		DefaultBranch: versionControlSettingsMap["default_branch"].(string),
		HasValue:      versionControlSettingsMap["has_value"].(bool),
		URL:           versionControlSettingsMap["url"].(string),
		Username:      versionControlSettingsMap["username"].(string),
	}

	if password := versionControlSettingsMap["password"].(string); len(password) > 0 {
		versionControlSettings.Password = octopusdeploy.NewSensitiveValue(password)
	}

	return versionControlSettings
}

func flattenVersionControlSettings(versionControlSettings *octopusdeploy.VersionControlSettings) []interface{} {
	if versionControlSettings == nil {
		return nil
	}

	flattenedVersionControlSettings := make(map[string]interface{})
	flattenedVersionControlSettings["base_path"] = versionControlSettings.BasePath
	flattenedVersionControlSettings["default_branch"] = versionControlSettings.DefaultBranch
	flattenedVersionControlSettings["has_value"] = versionControlSettings.HasValue
	flattenedVersionControlSettings["url"] = versionControlSettings.URL
	flattenedVersionControlSettings["username"] = versionControlSettings.Username
	return []interface{}{flattenedVersionControlSettings}
}

func getVersionControlSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"base_path": {
			Description: "The base path associated with these version control settings.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"default_branch": {
			Description: "The default branch associated with these version control settings.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"has_value": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"password": {
			Description:      "The password associated with these version control settings.",
			Sensitive:        true,
			Optional:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
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
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
	}
}
