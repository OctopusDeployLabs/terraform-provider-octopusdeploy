package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getDeploymentActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addActionTypeSchema(element)
	addExecutionLocationSchema(element)
	element.Schema[constActionType] = &schema.Schema{
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
			constName: {
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
			constRequired: {
				Type:        schema.TypeBool,
				Description: "Whether this step is required and cannot be skipped",
				Optional:    true,
				Default:     false,
			},
			constEnvironments: {
				Description: "The environments that this step will run in",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constExcludedEnvironments: {
				Description: "The environments that this step will be skipped in",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constChannels: {
				Description: "The channels that this step applies to",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constTenantTags: {
				Description: "The tags for the tenants that this step applies to",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constProperty: getPropertySchema(),
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
	element.Schema[constRunOnServer] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Whether this step runs on a worker or on the target",
		Optional:    true,
		Default:     false,
	}
}

func addActionTypeSchema(element *schema.Resource) {
	element.Schema[constActionType] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The type of action",
		Required:    true,
	}
}

func addWorkerPoolSchema(element *schema.Resource) {
	element.Schema[constWorkerPoolID] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Which worker pool to run on",
		Optional:    true,
	}
}

func buildDeploymentActionResource(tfAction map[string]interface{}) model.DeploymentAction {
	action := model.DeploymentAction{
		Name:                 tfAction[constName].(string),
		IsDisabled:           tfAction["disabled"].(bool),
		IsRequired:           tfAction[constRequired].(bool),
		Environments:         getSliceFromTerraformTypeList(tfAction[constEnvironments]),
		ExcludedEnvironments: getSliceFromTerraformTypeList(tfAction[constExcludedEnvironments]),
		Channels:             getSliceFromTerraformTypeList(tfAction[constChannels]),
		TenantTags:           getSliceFromTerraformTypeList(tfAction[constTenantTags]),
		Properties:           map[string]string{},
	}

	actionType := tfAction[constActionType]
	if actionType != nil {
		action.ActionType = actionType.(string)
	}

	// Even though not all actions have these properties, we'll keep them here.
	// They will just be ignored if the action doesn't have it
	runOnServer := tfAction[constRunOnServer]
	if runOnServer != nil {
		action.Properties["Octopus.Action.RunOnServer"] = strconv.FormatBool(runOnServer.(bool))
	}

	workerPoolID := tfAction[constWorkerPoolID]
	if workerPoolID != nil {
		action.WorkerPoolID = workerPoolID.(string)
	}

	if primaryPackage, ok := tfAction[constPrimaryPackage]; ok {
		tfPrimaryPackage := primaryPackage.(*schema.Set).List()
		if len(tfPrimaryPackage) > 0 {
			primaryPackage := buildPackageReferenceResource(tfPrimaryPackage[0].(map[string]interface{}))
			action.Packages = append(action.Packages, primaryPackage)
		}
	}

	if tfPkgs, ok := tfAction[constPackage]; ok {
		for _, tfPkg := range tfPkgs.(*schema.Set).List() {
			pkg := buildPackageReferenceResource(tfPkg.(map[string]interface{}))
			action.Packages = append(action.Packages, pkg)
		}
	}

	if tfProps, ok := tfAction[constProperty]; ok {
		for _, tfProp := range tfProps.(*schema.Set).List() {
			tfPropi := tfProp.(map[string]interface{})
			action.Properties[tfPropi[constKey].(string)] = tfPropi[constValue].(string)
		}
	}

	return action
}
