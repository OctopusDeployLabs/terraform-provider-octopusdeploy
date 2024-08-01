package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/extensions"
	env "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandEnvironment(d *schema.ResourceData) *environments.Environment {
	name := d.Get("name").(string)

	environment := environments.NewEnvironment(name)
	environment.ID = d.Id()

	if v, ok := d.GetOk("allow_dynamic_infrastructure"); ok {
		environment.AllowDynamicInfrastructure = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		environment.Description = v.(string)
	}

	if v, ok := d.GetOk("jira_extension_settings"); ok {
		environment.ExtensionSettings = append(environment.ExtensionSettings, env.ExpandJiraExtensionSettings(v))
	}

	if v, ok := d.GetOk("jira_service_management_extension_settings"); ok {
		environment.ExtensionSettings = append(environment.ExtensionSettings, env.ExpandJiraServiceManagementExtensionSettings(v))
	}

	if v, ok := d.GetOk("servicenow_extension_settings"); ok {
		environment.ExtensionSettings = append(environment.ExtensionSettings, env.ExpandServiceNowExtensionSettings(v))
	}

	if v, ok := d.GetOk("slug"); ok {
		environment.Slug = v.(string)
	}

	if v, ok := d.GetOk("sort_order"); ok {
		environment.SortOrder = v.(int)
	}

	if v, ok := d.GetOk("space_id"); ok {
		environment.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("use_guided_failure"); ok {
		environment.UseGuidedFailure = v.(bool)
	}

	return environment
}

func flattenEnvironment(environment *environments.Environment) map[string]interface{} {
	if environment == nil {
		return nil
	}

	environmentMap := map[string]interface{}{
		"allow_dynamic_infrastructure": environment.AllowDynamicInfrastructure,
		"description":                  environment.Description,
		"id":                           environment.GetID(),
		"name":                         environment.Name,
		"slug":                         environment.Slug,
		"sort_order":                   environment.SortOrder,
		"space_id":                     environment.SpaceID,
		"use_guided_failure":           environment.UseGuidedFailure,
	}

	if len(environment.ExtensionSettings) != 0 {
		for _, extensionSettings := range environment.ExtensionSettings {
			switch extensionSettings.ExtensionID() {
			case extensions.JiraExtensionID:
				if jiraExtensionSettings, ok := extensionSettings.(*environments.JiraExtensionSettings); ok {
					environmentMap["jira_extension_settings"] = env.FlattenJiraExtensionSettings(jiraExtensionSettings)
				}
			case extensions.JiraServiceManagementExtensionID:
				if jiraServiceManagementExtensionSettings, ok := extensionSettings.(*environments.JiraServiceManagementExtensionSettings); ok {
					environmentMap["jira_service_management_extension_settings"] = env.FlattenJiraServiceManagementExtensionSettings(jiraServiceManagementExtensionSettings)
				}
			case extensions.ServiceNowExtensionID:
				if serviceNowExtensionSettings, ok := extensionSettings.(*environments.ServiceNowExtensionSettings); ok {
					environmentMap["servicenow_extension_settings"] = env.FlattenServiceNowExtensionSettings(serviceNowExtensionSettings)
				}
			}
		}
	}

	return environmentMap
}

func getEnvironmentDataSchema() map[string]*schema.Schema {
	dataSchema := getEnvironmentSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"environments": {
			Computed:    true,
			Description: "A list of environments that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    false,
			Type:        schema.TypeList,
		},
		"id":           getDataSchemaID(),
		"ids":          getQueryIDs(),
		"name":         getQueryName(),
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
		"space_id":     getSpaceIDSchema(),
	}
}

func getEnvironmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_dynamic_infrastructure": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"description": getDescriptionSchema("environment"),
		"id":          getIDSchema(),
		"jira_extension_settings": {
			Description: "Provides extension settings for the Jira integration for this environment.",
			Elem:        &schema.Resource{Schema: env.GetJiraExtensionSettingsSchema()},
			MaxItems:    1,
			Optional:    true,
			Type:        schema.TypeList,
		},
		"jira_service_management_extension_settings": {
			Description: "Provides extension settings for the Jira Service Management (JSM) integration for this environment.",
			Elem:        &schema.Resource{Schema: env.GetJiraServiceManagementExtensionSettingsSchema()},
			MaxItems:    1,
			Optional:    true,
			Type:        schema.TypeList,
		},
		"name": getNameSchema(true),
		"servicenow_extension_settings": {
			Description: "Provides extension settings for the ServiceNow integration for this environment.",
			Elem:        &schema.Resource{Schema: env.GetServiceNowExtensionSettingsSchema()},
			MaxItems:    1,
			Optional:    true,
			Type:        schema.TypeList,
		},
		"slug": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"sort_order": {
			Computed:    true,
			Description: "The order number to sort an environment.",
			Optional:    true,
			Type:        schema.TypeInt,
		},
		"space_id": {
			Computed:         true,
			Description:      "The space ID associated with this environment.",
			Optional:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"use_guided_failure": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}

func setEnvironment(ctx context.Context, d *schema.ResourceData, environment *environments.Environment) error {
	d.Set("allow_dynamic_infrastructure", environment.AllowDynamicInfrastructure)
	d.Set("description", environment.Description)

	if len(environment.ExtensionSettings) != 0 {
		if err := env.SetExtensionSettings(d, environment.ExtensionSettings); err != nil {
			return fmt.Errorf("error setting extension settings: %s", err)
		}
	}

	d.Set("name", environment.Name)
	d.Set("slug", environment.Slug)
	d.Set("sort_order", environment.SortOrder)
	d.Set("space_id", environment.SpaceID)
	d.Set("use_guided_failure", environment.UseGuidedFailure)

	d.SetId(environment.GetID())

	return nil
}
