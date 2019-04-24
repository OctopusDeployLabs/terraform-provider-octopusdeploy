package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func getDeploymentActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addActionTypeSchema(element)
	addExecutionLocationSchema(element)
	element.Schema["action_type"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The type of action",
		Required:    true,
	}
	addWorkerPoolSchema(element)
	addPackagesSchema(element, false)

	return actionSchema
}

func getCommonDeploymentActionSchema() (*schema.Schema, *schema.Resource) {
	element := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the action",
				Required:    true,
			},
			"disabled": {
				Type:        schema.TypeBool,
				Description: "Whether this step is disabled",
				Optional:    true,
				Default:     false,
			},
			"required": {
				Type:        schema.TypeBool,
				Description: "Whether this step is required and cannot be skipped",
				Optional:    true,
				Default:     false,
			},
			"environments": {
				Description: "The environments that this step will run in",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"excluded_environments": {
				Description: "The environments that this step will be skipped in",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"channels": {
				Description: "The channels that this step applies to",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenant_tags": {
				Description: "The tags for the tenants that this step applies to",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"property": getPropertySchema(),
		},
	}

	actionSchema := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     element,
	}

	return actionSchema, element
}

func addExecutionLocationSchema(element *schema.Resource) {
	element.Schema["run_on_server"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Whether this step runs on a worker or on the target",
		Optional:    true,
		Default:     false,
	}
}

func addActionTypeSchema(element *schema.Resource) {
	element.Schema["action_type"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The type of action",
		Required:    true,
	}
}

func addWorkerPoolSchema(element *schema.Resource) {
	element.Schema["worker_pool_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Which worker pool to run on",
		Optional:    true,
	}
}

func buildDeploymentActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := octopusdeploy.DeploymentAction{
		Name:                 tfAction["name"].(string),
		IsDisabled:           tfAction["disabled"].(bool),
		IsRequired:           tfAction["required"].(bool),
		Environments:         getSliceFromTerraformTypeList(tfAction["environments"]),
		ExcludedEnvironments: getSliceFromTerraformTypeList(tfAction["excluded_environments"]),
		Channels:             getSliceFromTerraformTypeList(tfAction["channels"]),
		TenantTags:           getSliceFromTerraformTypeList(tfAction["tenant_tags"]),
		Properties:           map[string]string{},
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
		action.WorkerPoolId = workerPoolID.(string)
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
