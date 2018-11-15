package octopusdeploy

import (
	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func getDeploymentActionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "The name of the action",
					Required:    true,
				},
				"action_type": {
					Type:        schema.TypeString,
					Description: "The type of action",
					Required:    true,
				},
				"run_on_server": {
					Type:        schema.TypeBool,
					Description: "Whether this step is disabled",
					Optional:    true,
					Default: 	 false,
				},
				"disabled": {
					Type:        schema.TypeBool,
					Description: "Whether this step is disabled",
					Optional:    true,
					Default: 	 false,
				},
				"required": {
					Type:        schema.TypeBool,
					Description: "Whether this step is required and cannot be skipped",
					Optional:    true,
					Default: 	 false,
				},
				"worker_pool_id": {
					Type:        schema.TypeString,
					Description: "Which worker pool to run on",
					Optional:    true,
				},
				"environments": &schema.Schema{
					Description: "The environments that this step will run in",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"excluded_environments": &schema.Schema{
					Description: "The environments that this still will be skipped in",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"channels": &schema.Schema{
					Description: "The channels that this step applies to",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"tenant_tags": &schema.Schema{
					Description: "The tags for the tenants that this step applies to",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"property": getPropertySchema(),
				"primary_package": getPrimaryPackageSchema(),
				"package": getPackageSchema(),
			},
		},
	}
}


func buildDeploymentActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := octopusdeploy.DeploymentAction{
		Name:                 tfAction["name"].(string),
		ActionType:           tfAction["action_type"].(string),
		IsDisabled:           tfAction["disabled"].(bool),
		IsRequired:           tfAction["required"].(bool),
		Environments:         getSliceFromTerraformTypeList(tfAction["environments"]),
		ExcludedEnvironments: getSliceFromTerraformTypeList(tfAction["excluded_environments"]),
		Channels:             getSliceFromTerraformTypeList(tfAction["channels"]),
		TenantTags:           getSliceFromTerraformTypeList(tfAction["tenant_tags"]),
		Properties:           map[string]string{},
	}

	runOnServer := tfAction["run_on_server"]
	if runOnServer != nil {
		action.Properties["Octopus.Action.RunOnServer"] = strconv.FormatBool(runOnServer.(bool))
	}

	workerPoolId := tfAction["worker_pool_id"]
	if workerPoolId != nil {
		action.WorkerPoolId = workerPoolId.(string)
	}

	if primaryPackage, ok := tfAction["primary_package"]; ok {
		tfPrimaryPackage := primaryPackage.(*schema.Set).List()
		if(len(tfPrimaryPackage) > 0) {
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
		for _, tfProp := range  tfProps.(*schema.Set).List() {
			tfPropi := tfProp.(map[string]interface{})
			action.Properties[tfPropi["key"].(string)] = tfPropi["value"].(string)
		}
	}

	return action;
}

