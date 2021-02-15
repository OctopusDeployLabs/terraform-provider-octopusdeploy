package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenWindowsService(actionMap map[string]interface{}, properties map[string]string) {
	for propertyName, propertyValue := range properties {
		switch propertyName {
		case "Octopus.Action.WindowsService.Arguments":
			actionMap["arguments"] = propertyValue
		case "Octopus.Action.WindowsService.CreateOrUpdateService":
			createOrUpdateService, _ := strconv.ParseBool(propertyValue)
			actionMap["create_or_update_service"] = createOrUpdateService
		case "Octopus.Action.WindowsService.CustomAccountName":
			actionMap["custom_account_name"] = propertyValue
		case "Octopus.Action.WindowsService.CustomAccountPassword":
			actionMap["custom_account_password"] = propertyValue
		case "Octopus.Action.WindowsService.Dependencies":
			actionMap["dependencies"] = propertyValue
		case "Octopus.Action.WindowsService.Description":
			actionMap["description"] = propertyValue
		case "Octopus.Action.WindowsService.DisplayName":
			actionMap["display_name"] = propertyValue
		case "Octopus.Action.WindowsService.ExecutablePath":
			actionMap["executable_path"] = propertyValue
		case "Octopus.Action.WindowsService.ServiceAccount":
			actionMap["service_account"] = propertyValue
		case "Octopus.Action.WindowsService.ServiceName":
			actionMap["service_name"] = propertyValue
		case "Octopus.Action.WindowsService.StartMode":
			actionMap["start_mode"] = propertyValue
		}
	}
}

func flattenWindowsServiceAction(deploymentAction octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedWindowsServiceAction := map[string]interface{}{
		"can_be_used_for_project_versioning": deploymentAction.CanBeUsedForProjectVersioning,
		"channels":                           deploymentAction.Channels,
		"condition":                          deploymentAction.Condition,
		"container":                          flattenDeploymentActionContainer(deploymentAction.Container),
		"environments":                       deploymentAction.Environments,
		"excluded_environments":              deploymentAction.ExcludedEnvironments,
		"id":                                 deploymentAction.ID,
		"is_disabled":                        deploymentAction.IsDisabled,
		"is_required":                        deploymentAction.IsRequired,
		"name":                               deploymentAction.Name,
		"notes":                              deploymentAction.Notes,
		"package":                            flattenPackageReferences(deploymentAction.Packages),
		"properties":                         deploymentAction.Properties,
		"tenant_tags":                        deploymentAction.TenantTags,
	}

	flattenWindowsService(flattenedWindowsServiceAction, deploymentAction.Properties)

	return flattenedWindowsServiceAction
}

func getWindowsServiceActionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"can_be_used_for_project_versioning": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"channels": {
			Computed:    true,
			Description: "The channels associated with this deployment action.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"condition": {
			Computed:    true,
			Description: "The condition associated with this deployment action.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"container": {
			Computed:    true,
			Description: "The deployment action container associated with this deployment action.",
			Elem:        &schema.Resource{Schema: getDeploymentActionContainerSchema()},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"excluded_environments": {
			Computed:    true,
			Description: "The environments that this step will be skipped in",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"environments": {
			Computed:    true,
			Description: "The environments within which this deployment action will run.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"id": getIDSchema(),
		"is_disabled": {
			Default:     false,
			Description: "Indicates the disabled status of this deployment action.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"is_required": {
			Default:     false,
			Description: "Indicates the required status of this deployment action.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"name": getNameSchema(true),
		"notes": {
			Description: "The notes associated with this deploymnt action.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"package": getPackageSchema(false),
		"properties": {
			Computed:    true,
			Description: "The properties associated with this deployment action.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeMap,
		},
		"tenant_tags": getTenantTagsSchema(),
	}
}
