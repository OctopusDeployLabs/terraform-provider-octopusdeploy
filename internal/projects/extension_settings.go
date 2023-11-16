package projects

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/extensions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ExpandJiraExtensionSettings deserializes the project extension settings for Jira Service Management (JSM) integration from its HCL representation.
func ExpandJiraServiceManagementExtensionSettings(extensionSettings interface{}) extensions.ExtensionSettings {
	values := extensionSettings.([]interface{})
	valuesMap := values[0].(map[string]interface{})
	return projects.NewJiraServiceManagementExtensionSettings(
		valuesMap["connection_id"].(string),
		valuesMap["is_enabled"].(bool),
		valuesMap["service_desk_project_name"].(string),
	)
}

// ExpandJiraExtensionSettings deserializes the project extension settings for ServiceNow integration from its HCL representation.
func ExpandServiceNowExtensionSettings(extensionSettings interface{}) extensions.ExtensionSettings {
	values := extensionSettings.([]interface{})
	valuesMap := values[0].(map[string]interface{})
	return projects.NewServiceNowExtensionSettings(
		valuesMap["connection_id"].(string),
		valuesMap["is_enabled"].(bool),
		valuesMap["standard_change_template_name"].(string),
		valuesMap["is_state_automatically_transitioned"].(bool),
	)
}

// FlattenJiraServiceManagementExtensionSettings serializes the project extension settings for Jira Service Management (JSM) integration into its HCL representation.
func FlattenJiraServiceManagementExtensionSettings(jiraServiceManagementExtensionSettings *projects.JiraServiceManagementExtensionSettings) []interface{} {
	if jiraServiceManagementExtensionSettings == nil {
		return nil
	}

	flattenedJiraServiceManagementExtensionSettings := make(map[string]interface{})
	flattenedJiraServiceManagementExtensionSettings["connection_id"] = jiraServiceManagementExtensionSettings.ConnectionID()
	flattenedJiraServiceManagementExtensionSettings["is_enabled"] = jiraServiceManagementExtensionSettings.IsChangeControlled()
	flattenedJiraServiceManagementExtensionSettings["service_desk_project_name"] = jiraServiceManagementExtensionSettings.ServiceDeskProjectName
	return []interface{}{flattenedJiraServiceManagementExtensionSettings}
}

// FlattenServiceNowExtensionSettings serializes the project extension settings for ServiceNow integration into its HCL representation.
func FlattenServiceNowExtensionSettings(serviceNowExtensionSettings *projects.ServiceNowExtensionSettings) []interface{} {
	if serviceNowExtensionSettings == nil {
		return nil
	}

	flattenedServiceNowExtensionSettings := make(map[string]interface{})
	flattenedServiceNowExtensionSettings["connection_id"] = serviceNowExtensionSettings.ConnectionID()
	flattenedServiceNowExtensionSettings["is_enabled"] = serviceNowExtensionSettings.IsChangeControlled()
	flattenedServiceNowExtensionSettings["is_state_automatically_transitioned"] = serviceNowExtensionSettings.IsStateAutomaticallyTransitioned
	flattenedServiceNowExtensionSettings["standard_change_template_name"] = serviceNowExtensionSettings.StandardChangeTemplateName
	return []interface{}{flattenedServiceNowExtensionSettings}
}

// GetJiraServiceManagementExtensionSettingsSchema returns the Terraform schema for Jira Service Management (JSM) integration with projects.
func GetJiraServiceManagementExtensionSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_id": {
			Description:      "The connection identifier associated with the extension settings.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"is_enabled": {
			Description: "Specifies whether or not this extension is enabled for this project.",
			Required:    true,
			Type:        schema.TypeBool,
		},
		"service_desk_project_name": {
			Description:      "The project name associated with this extension.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
	}
}

// GetServiceNowExtensionSettingsSchema returns the Terraform schema for ServiceNow integration with projects.
func GetServiceNowExtensionSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_id": {
			Description:      "The connection identifier associated with the extension settings.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"is_enabled": {
			Description: "Specifies whether or not this extension is enabled for this project.",
			Required:    true,
			Type:        schema.TypeBool,
		},
		"is_state_automatically_transitioned": {
			Description: "Specifies whether or not this extension will automatically transition the state of a deployment for this project.",
			Required:    true,
			Type:        schema.TypeBool,
		},
		"standard_change_template_name": {
			Description:      "The name of the standard change template associated with this extension.",
			Required:         false,
			Type:             schema.TypeString
		},
	}
}

// SetExtensionSettings sets the Terraform state of project settings collection for extensions.
func SetExtensionSettings(d *schema.ResourceData, extensionSettingsCollection []extensions.ExtensionSettings) error {
	for _, extensionSettings := range extensionSettingsCollection {
		switch extensionSettings.ExtensionID() {
		case extensions.JiraServiceManagementExtensionID:
			if jiraServiceManagementExtensionSettings, ok := extensionSettings.(*projects.JiraServiceManagementExtensionSettings); ok {
				if err := d.Set("jira_service_management_extension_settings", FlattenJiraServiceManagementExtensionSettings(jiraServiceManagementExtensionSettings)); err != nil {
					return fmt.Errorf("error setting extension settings for Jira Service Management (JSM): %s", err)
				}
			}
		case extensions.ServiceNowExtensionID:
			if serviceNowExtensionSettings, ok := extensionSettings.(*projects.ServiceNowExtensionSettings); ok {
				if err := d.Set("servicenow_extension_settings", FlattenServiceNowExtensionSettings(serviceNowExtensionSettings)); err != nil {
					return fmt.Errorf("error setting extension settings for ServiceNow: %s", err)
				}
			}
		}
	}

	return nil
}
