package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenDeploymentActions(deploymentActions []octopusdeploy.DeploymentAction) []map[string]interface{} {
	var flattenedDeploymentActions = make([]map[string]interface{}, len(deploymentActions))
	for key, deploymentAction := range deploymentActions {
		flattenedDeploymentActions[key] = map[string]interface{}{
			"action_type":                        deploymentAction.ActionType,
			"can_be_used_for_project_versioning": deploymentAction.CanBeUsedForProjectVersioning,
			"channels":                           deploymentAction.Channels,
			"environments":                       deploymentAction.Environments,
			"excluded_environments":              deploymentAction.ExcludedEnvironments,
			"id":                                 deploymentAction.ID,
			"is_disabled":                        deploymentAction.IsDisabled,
			"is_required":                        deploymentAction.IsRequired,
			"name":                               deploymentAction.Name,
			"packages":                           deploymentAction.Packages,
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
			"channels": {
				Description: "The channels that this step applies to",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"disabled": {
				Default:     false,
				Description: "Whether this step is disabled",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"excluded_environments": {
				Description: "The environments that this step will be skipped in",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"environments": {
				Description: "The environments that this step will run in",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"name": {
				Description: "The name of the action",
				Required:    true,
				Type:        schema.TypeString,
			},
			"required": {
				Default:     false,
				Description: "Whether this step is required and cannot be skipped",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"property": getPropertySchema(),
			"tenant_tags": {
				Description: "The tags for the tenants that this step applies to",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
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
		Description: "Which worker pool to run on",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func expandDeploymentAction(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := octopusdeploy.DeploymentAction{
		Channels:             getSliceFromTerraformTypeList(tfAction["channels"]),
		Environments:         getSliceFromTerraformTypeList(tfAction["environments"]),
		ExcludedEnvironments: getSliceFromTerraformTypeList(tfAction["excluded_environments"]),
		IsDisabled:           tfAction["disabled"].(bool),
		IsRequired:           tfAction["required"].(bool),
		Name:                 tfAction["name"].(string),
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
