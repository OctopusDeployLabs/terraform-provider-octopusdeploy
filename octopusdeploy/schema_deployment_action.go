package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenDeploymentAction(action octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedDeploymentAction := flattenAction(action)

	flattenedDeploymentAction["action_type"] = action.ActionType
	flattenedDeploymentAction["worker_pool_id"] = action.WorkerPoolID

	return flattenedDeploymentAction
}

func flattenAction(action octopusdeploy.DeploymentAction) map[string]interface{} {
	actionProperties := action.Properties

	flattenedAction := map[string]interface{}{
		"can_be_used_for_project_versioning": action.CanBeUsedForProjectVersioning,
		"channels":                           action.Channels,
		"container":                          flattenDeploymentActionContainer(action.Container),
		"condition":                          action.Condition,
		"environments":                       action.Environments,
		"excluded_environments":              action.ExcludedEnvironments,
		"id":                                 action.ID,
		"is_disabled":                        action.IsDisabled,
		"is_required":                        action.IsRequired,
		"name":                               action.Name,
		"notes":                              action.Notes,
		"tenant_tags":                        action.TenantTags,
	}

	flattenedPackageReferences := []interface{}{}
	for _, packageReference := range action.Packages {
		flattenedPackageReference := flattenPackageReference(packageReference)
		if len(packageReference.Name) == 0 {
			flattenedAction["primary_package"] = []interface{}{flattenedPackageReference}
			actionProperties["Octopus.Action.Package.DownloadOnTentacle"] = packageReference.AcquisitionLocation
			flattenedAction["properties"] = actionProperties
			continue
		}

		if v, ok := packageReference.Properties["Extract"]; ok {
			extractDuringDeployment, _ := strconv.ParseBool(v)
			flattenedPackageReference["extract_during_deployment"] = extractDuringDeployment
		}

		flattenedPackageReferences = append(flattenedPackageReferences, flattenedPackageReference)
	}
	flattenedAction["package"] = flattenedPackageReferences
	flattenedAction["properties"] = actionProperties

	return flattenedAction
}

func getDeploymentActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addActionTypeSchema(element)
	addExecutionLocationSchema(element)
	element.Schema["action_type"] = &schema.Schema{
		Description: "The type of action",
		Required:    true,
		Type:        schema.TypeString,
	}
	addWorkerPoolSchema(element)
	addPackagesSchema(element, false)

	return actionSchema
}

func getCommonDeploymentActionSchema() (*schema.Schema, *schema.Resource) {
	element := &schema.Resource{
		Schema: map[string]*schema.Schema{
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
			"environments": {
				Computed:    true,
				Description: "The environments within which this deployment action will run.",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Optional:    true,
				Type:        schema.TypeMap,
			},
			"tenant_tags": getTenantTagsSchema(),
		},
	}

	actionSchema := &schema.Schema{
		Computed: true,
		Elem:     element,
		Optional: true,
		Type:     schema.TypeList,
	}

	return actionSchema, element
}

func addExecutionLocationSchema(element *schema.Resource) {
	element.Schema["run_on_server"] = &schema.Schema{
		Default:     false,
		Description: "Whether this step runs on a worker or on the target",
		Optional:    true,
		Type:        schema.TypeBool,
	}
}

func addActionTypeSchema(element *schema.Resource) {
	element.Schema["action_type"] = &schema.Schema{
		Description: "The type of action",
		Required:    true,
		Type:        schema.TypeString,
	}
}

func addWorkerPoolSchema(element *schema.Resource) {
	element.Schema["worker_pool_id"] = &schema.Schema{
		Description: "The worker pool associated with this deployment action.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func expandAction(flattenedAction map[string]interface{}) octopusdeploy.DeploymentAction {
	name := flattenedAction["name"].(string)
	action := octopusdeploy.NewDeploymentAction(name)

	if v, ok := flattenedAction["action_type"]; ok {
		if actionType := v.(string); len(actionType) > 0 {
			action.ActionType = actionType
		}
	}

	if v, ok := flattenedAction["channels"]; ok {
		action.Channels = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedAction["condition"]; ok {
		action.Condition = v.(string)
	}

	if v, ok := flattenedAction["container"]; ok {
		action.Container = expandDeploymentActionContainer(v)
	}

	if v, ok := flattenedAction["environments"]; ok {
		action.Environments = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedAction["excluded_environments"]; ok {
		action.ExcludedEnvironments = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedAction["is_disabled"]; ok {
		action.IsDisabled = v.(bool)
	}

	if v, ok := flattenedAction["is_required"]; ok {
		action.IsRequired = v.(bool)
	}

	if v, ok := flattenedAction["notes"]; ok {
		action.Notes = v.(string)
	}

	if v, ok := flattenedAction["properties"]; ok {
		action.Properties = expandProperties(v)
	}

	// Even though not all actions have these properties, we'll keep them here.
	// They will just be ignored if the action doesn't have it
	if v, ok := flattenedAction["run_on_server"]; ok {
		action.Properties["Octopus.Action.RunOnServer"] = strconv.FormatBool(v.(bool))
	}

	if v, ok := flattenedAction["tenant_tags"]; ok {
		action.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedAction["worker_pool_id"]; ok {
		action.WorkerPoolID = v.(string)
	}

	if v, ok := flattenedAction["primary_package"]; ok {
		primaryPackages := v.([]interface{})
		for _, primaryPackage := range primaryPackages {
			primaryPackageReference := expandPackageReference(primaryPackage.(map[string]interface{}))

			action.Properties["Octopus.Action.Package.DownloadOnTentacle"] = primaryPackageReference.AcquisitionLocation

			if len(primaryPackageReference.PackageID) > 0 {
				action.Properties["Octopus.Action.Package.PackageId"] = primaryPackageReference.PackageID
			}

			if len(primaryPackageReference.FeedID) > 0 {
				action.Properties["Octopus.Action.Package.FeedId"] = primaryPackageReference.FeedID
			}

			action.Packages = append(action.Packages, primaryPackageReference)
		}
	}

	if v, ok := flattenedAction["package"]; ok {
		packageReferences := v.([]interface{})
		for _, packageReference := range packageReferences {
			action.Packages = append(action.Packages, expandPackageReference(packageReference.(map[string]interface{})))
		}
	}

	return *action
}
