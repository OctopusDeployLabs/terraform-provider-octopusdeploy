package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenDeploymentActions(deploymentActions []octopusdeploy.DeploymentAction) []map[string]interface{} {
	if deploymentActions == nil {
		return nil
	}

	var flattenedDeploymentActions = make([]map[string]interface{}, len(deploymentActions))
	for key, deploymentAction := range deploymentActions {
		flattenedDeploymentActions[key] = map[string]interface{}{
			"action_type":                        deploymentAction.ActionType,
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
			"package":                            deploymentAction.Packages,
			"properties":                         deploymentAction.Properties,
			"tenant_tags":                        deploymentAction.TenantTags,
			"worker_pool_id":                     deploymentAction.WorkerPoolID,
		}
	}

	return flattenedDeploymentActions
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
				Optional: true,
				Type:     schema.TypeBool,
			},
			"channels": {
				Description: "The channels associated with this deployment action.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"condition": {
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
				Description: "The environments that this step will be skipped in",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"environments": {
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
			"package": getPackageSchema(true),
			"properties": {
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Type:     schema.TypeMap,
			},
			"tenant_tags": getTenantTagsSchema(),
		},
	}

	actionSchema := &schema.Schema{
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

func expandDeploymentAction(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := octopusdeploy.DeploymentAction{
		Channels:             getSliceFromTerraformTypeList(tfAction["channels"]),
		Container:            tfAction["container"].(octopusdeploy.DeploymentActionContainer),
		Condition:            tfAction["condition"].(string),
		Environments:         getSliceFromTerraformTypeList(tfAction["environments"]),
		ExcludedEnvironments: getSliceFromTerraformTypeList(tfAction["excluded_environments"]),
		IsDisabled:           tfAction["is_disabled"].(bool),
		IsRequired:           tfAction["is_required"].(bool),
		Name:                 tfAction["name"].(string),
		Notes:                tfAction["notes"].(string),
		Properties:           map[string]string{},
		TenantTags:           getSliceFromTerraformTypeList(tfAction["tenant_tags"]),
	}

	actionType := tfAction["action_type"]
	if actionType != nil {
		action.ActionType = actionType.(string)
	}

	// Even though not all actions have these properties, we'll keep them here.
	// They will just be ignored if the action doesn't have it
	runOnServer := tfAction["run_on_server"]
	if runOnServer != nil {
		action.Properties["Octopus.Action.RunOnServer"] = strconv.FormatBool(runOnServer.(bool))
	}

	workerPoolID := tfAction["worker_pool_id"]
	if workerPoolID != nil {
		action.WorkerPoolID = workerPoolID.(string)
	}

	if primaryPackage, ok := tfAction["primary_package"]; ok {
		tfPrimaryPackage := primaryPackage.(*schema.Set).List()
		if len(tfPrimaryPackage) > 0 {
			primaryPackage := buildPackageReferenceResource(tfPrimaryPackage[0].(map[string]interface{}))
			action.Packages = append(action.Packages, primaryPackage)
		}
	}

	if tfPkgs, ok := tfAction["package"]; ok {
		for _, tfPkg := range tfPkgs.(*schema.Set).List() {
			pkg := buildPackageReferenceResource(tfPkg.(map[string]interface{}))
			action.Packages = append(action.Packages, pkg)
		}
	}

	if tfProps, ok := tfAction["property"]; ok {
		for _, tfProp := range tfProps.(*schema.Set).List() {
			tfPropi := tfProp.(map[string]interface{})
			action.Properties[tfPropi["key"].(string)] = tfPropi["value"].(string)
		}
	}

	return action
}
