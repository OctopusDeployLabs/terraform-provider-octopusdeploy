package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenDeploymentAction(deploymentAction octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedDeploymentAction := flattenCommonDeploymentAction(deploymentAction)

	flattenedDeploymentAction["action_type"] = deploymentAction.ActionType
	flattenedDeploymentAction["worker_pool_id"] = deploymentAction.WorkerPoolID

	return flattenedDeploymentAction
}

func flattenCommonDeploymentAction(deploymentAction octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedDeploymentAction := map[string]interface{}{
		"can_be_used_for_project_versioning": deploymentAction.CanBeUsedForProjectVersioning,
		"channels":                           deploymentAction.Channels,
		"container":                          flattenDeploymentActionContainer(deploymentAction.Container),
		"condition":                          deploymentAction.Condition,
		"environments":                       deploymentAction.Environments,
		"excluded_environments":              deploymentAction.ExcludedEnvironments,
		"id":                                 deploymentAction.ID,
		"is_disabled":                        deploymentAction.IsDisabled,
		"is_required":                        deploymentAction.IsRequired,
		"name":                               deploymentAction.Name,
		"notes":                              deploymentAction.Notes,
		"properties":                         deploymentAction.Properties,
		"tenant_tags":                        deploymentAction.TenantTags,
	}

	flattenedPackageReferences := []interface{}{}
	for _, packageReference := range deploymentAction.Packages {
		flattenedPackageReference := flattenPackageReference(packageReference)
		if len(packageReference.Name) == 0 {
			flattenedDeploymentAction["primary_package"] = []interface{}{flattenedPackageReference}
			continue
		}

		if v, ok := packageReference.Properties["Extract"]; ok {
			extractDuringDeployment, _ := strconv.ParseBool(v)
			flattenedPackageReference["extract_during_deployment"] = extractDuringDeployment
		}

		flattenedPackageReferences = append(flattenedPackageReferences, flattenedPackageReference)
	}
	flattenedDeploymentAction["package"] = flattenedPackageReferences

	return flattenedDeploymentAction
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
				Elem:        &schema.Schema{Type: schema.TypeString},
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

func expandDeploymentAction(flattenedDeploymentAction map[string]interface{}) octopusdeploy.DeploymentAction {
	name := flattenedDeploymentAction["name"].(string)
	action := octopusdeploy.NewDeploymentAction(name)

	if v, ok := flattenedDeploymentAction["action_type"]; ok {
		if actionType := v.(string); len(actionType) > 0 {
			action.ActionType = actionType
		}
	}

	if channels, ok := flattenedDeploymentAction["channels"]; ok {
		action.Channels = getSliceFromTerraformTypeList(channels)
	}

	if condition, ok := flattenedDeploymentAction["condition"]; ok {
		action.Condition = condition.(string)
	}

	if container, ok := flattenedDeploymentAction["container"]; ok {
		action.Container = expandDeploymentActionContainer(container)
	}

	if environments, ok := flattenedDeploymentAction["environments"]; ok {
		action.Environments = getSliceFromTerraformTypeList(environments)
	}

	if excludedEnvironments, ok := flattenedDeploymentAction["excluded_environments"]; ok {
		action.ExcludedEnvironments = getSliceFromTerraformTypeList(excludedEnvironments)
	}

	if isDisabled, ok := flattenedDeploymentAction["is_disabled"]; ok {
		action.IsDisabled = isDisabled.(bool)
	}

	if isRequired, ok := flattenedDeploymentAction["is_required"]; ok {
		action.IsRequired = isRequired.(bool)
	}

	if notes, ok := flattenedDeploymentAction["notes"]; ok {
		action.Notes = notes.(string)
	}

	if properties, ok := flattenedDeploymentAction["properties"]; ok {
		action.Properties = expandProperties(properties)
	}

	// Even though not all actions have these properties, we'll keep them here.
	// They will just be ignored if the action doesn't have it
	if runOnServer, ok := flattenedDeploymentAction["run_on_server"]; ok {
		action.Properties["Octopus.Action.RunOnServer"] = strconv.FormatBool(runOnServer.(bool))
	}

	if tenantTags, ok := flattenedDeploymentAction["tenant_tags"]; ok {
		action.TenantTags = getSliceFromTerraformTypeList(tenantTags)
	}

	if workerPoolID, ok := flattenedDeploymentAction["worker_pool_id"]; ok {
		action.WorkerPoolID = workerPoolID.(string)
	}

	if v, ok := flattenedDeploymentAction["primary_package"]; ok {
		primaryPackages := v.([]interface{})
		for _, primaryPackage := range primaryPackages {
			action.Packages = append(action.Packages, expandPackageReference(primaryPackage.(map[string]interface{})))
		}
	}

	if v, ok := flattenedDeploymentAction["package"]; ok {
		packageReferences := v.([]interface{})
		for _, packageReference := range packageReferences {
			action.Packages = append(action.Packages, expandPackageReference(packageReference.(map[string]interface{})))
		}
	}

	return *action
}
