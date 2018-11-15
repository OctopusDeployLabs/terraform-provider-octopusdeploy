package octopusdeploy

import (
	"fmt"
	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceDeploymentProcess() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentProcessCreate,
		Read:   resourceDeploymentProcessRead,
		Update: resourceDeploymentProcessUpdate,
		Delete: resourceDeploymentProcessDelete,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"step": getStepSchema(),
		},
	}
}


func getStepSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "The name of the step",
					Required:    true,
				},
				"target_roles": &schema.Schema{
					Description: "The roles that this step run against, or runs on behalf of",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"package_requirement": {
					Type:        schema.TypeString,
					Description: "Whether to run this step before or after package acquisition (if possible)",
					Optional:    true,
					Default:     (string)(octopusdeploy.DeploymentStepPackageRequirement_LetOctopusDecide),
					ValidateFunc: validateValueFunc([]string{
						(string)(octopusdeploy.DeploymentStepPackageRequirement_LetOctopusDecide),
						(string)(octopusdeploy.DeploymentStepPackageRequirement_BeforePackageAcquisition),
						(string)(octopusdeploy.DeploymentStepPackageRequirement_AfterPackageAcquisition),
					}),
				},
				"condition": {
					Type:        schema.TypeString,
					Description: "When to run the step, one of 'Success', 'Failure', 'Always' or 'Variable'",
					Optional:    true,
					Default:     (string)(octopusdeploy.DeploymentStepCondition_Success),
					ValidateFunc: validateValueFunc([]string{
						(string)(octopusdeploy.DeploymentStepCondition_Success),
						(string)(octopusdeploy.DeploymentStepCondition_Failure),
						(string)(octopusdeploy.DeploymentStepCondition_Always),
						(string)(octopusdeploy.DeploymentStepCondition_Variable),
					}),
				},
				"condition_expression": {
					Type:        schema.TypeString,
					Description: "The expression to evaluate to determine whether to run this step when 'condition' is 'Variable'",
					Optional:    true,
				},
				"start_trigger": {
					Type:        schema.TypeString,
					Description: "Whether to run this step after the previous step ('StartAfterPrevious') or at the same time as the previous step ('StartWithPrevious')",
					Optional:    true,
					Default:     (string)(octopusdeploy.DeploymentStepStartTrigger_StartAfterPrevious),
					ValidateFunc: validateValueFunc([]string{
						(string)(octopusdeploy.DeploymentStepStartTrigger_StartAfterPrevious),
						(string)(octopusdeploy.DeploymentStepStartTrigger_StartWithPrevious),
					}),
				},
				"window_size": {
					Type:        schema.TypeString,
					Description: "The maximum number of targets to deploy to simultaneously",
					Optional:    true,
				},
				"action": getActionSchema(),
			},
		},
	}
}

func getActionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
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
				"disabled": {
					Type:        schema.TypeString,
					Description: "Whether this step is disabled",
					Optional:    true,
					Default: 	 false,
				},
				"required": {
					Type:        schema.TypeString,
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

func getPrimaryPackageSchema() *schema.Schema {
	return &schema.Schema{
		Description: "The primary package for the action",
		Type:        schema.TypeSet,
		Optional:    true,
		MaxItems:	 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"package_id": {
					Type:        schema.TypeString,
					Description: "The ID of the package",
					Required:    true,
				},
				"feed_id": {
					Type:        schema.TypeString,
					Description: "The feed to retrieve the package from",
					Default: 	"feeds-builtin",
					Optional:    true,
				},
				"acquisition_location": {
					Type:        schema.TypeString,
					Description: "Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression",
					Default:     (string)(octopusdeploy.PackageAcquisitionLocation_Server),
					Optional:    true,
				},
			},
		},
	}
}

func getPackageSchema() *schema.Schema {
	s := getPrimaryPackageSchema();
	elementSchema := s.Elem.(*schema.Resource).Schema
	elementSchema["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The name of the package",
		Required:    true,
	}
	elementSchema["extract_during_deployment"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Whether to extract the package during deployment",
		Optional:    true,
	}
	elementSchema["property"] = getPropertySchema()
	return s
}

func getPropertySchema() *schema.Schema {
	return &schema.Schema{
		Description: "The tags for the tenants that this step applies to",
		Type:        schema.TypeSet,
		Optional:    true,
		Elem:  &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:        schema.TypeString,
					Description: "The name of the action",
					Required:    true,
				},
				"value": {
					Type:        schema.TypeString,
					Description: "The type of action",
					Required:    true,
				},
			},
		},
	}
}

func resourceDeploymentProcessCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newDeploymentProcess := buildDeploymentProcessResource(d)

	project, err := client.Project.Get(newDeploymentProcess.ProjectID)
	if err != nil {
		return fmt.Errorf("error getting project %s: %s", newDeploymentProcess.ProjectID, err.Error())
	}

	newDeploymentProcess.ID = project.DeploymentProcessID
	createdDeploymentProcess, err := client.DeploymentProcess.Update(newDeploymentProcess)

	if err != nil {
		return fmt.Errorf("error creating deployment process: %s", err.Error())
	}

	d.SetId(createdDeploymentProcess.ID)

	return nil
}

func buildDeploymentProcessResource(d *schema.ResourceData) *octopusdeploy.DeploymentProcess {
	deploymentProcess := &octopusdeploy.DeploymentProcess {
		ProjectID: d.Get("project_id").(string),
	}

	if attr, ok := d.GetOk("step"); ok {
		tfSteps := attr.(*schema.Set)

		for _, tfStep := range tfSteps.List() {
			step := buildStepResource(tfStep.(map[string]interface{}))
			deploymentProcess.Steps = append(deploymentProcess.Steps, step)
		}
	}

	return deploymentProcess
}

func buildStepResource(tfStep map[string]interface{}) octopusdeploy.DeploymentStep {
	step := octopusdeploy.DeploymentStep{
		Name: tfStep["name"].(string),
		PackageRequirement: octopusdeploy.DeploymentStepPackageRequirement(tfStep["package_requirement"].(string)),
		Condition: octopusdeploy.DeploymentStepCondition(tfStep["package_requirement"].(string)),
		StartTrigger: octopusdeploy.DeploymentStepStartTrigger(tfStep["start_trigger"].(string)),
	}

	targetRoles := tfStep["target_roles"];
	if targetRoles != nil {
		step.Properties["Octopus.Action.TargetRoles"] = strings.Join(getSliceFromTerraformTypeList(targetRoles), ",")
	}

	conditionExpression := tfStep["condition_expression"]
	if conditionExpression != nil {
		step.Properties["Octopus.Action.ConditionVariableExpression"] = conditionExpression.(string)
	}

	windowSize := tfStep["window_size"]
	if windowSize != nil {
		step.Properties["Octopus.Action.ConditionVariableExpression"] = windowSize.(string)
	}

	if attr, ok := tfStep["action"]; ok {
		tfActions := attr.([]interface {})

		for _, tfAction := range tfActions {
			action := buildActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	return step;
}

func buildActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := octopusdeploy.DeploymentAction{
		Name:                 tfAction["name"].(string),
		ActionType:           tfAction["action_type"].(string),
		IsDisabled:           tfAction["disabled"].(bool),
		IsRequired:           tfAction["required"].(bool),
		Environments:         getSliceFromTerraformTypeList(tfAction["environments"]),
		ExcludedEnvironments: getSliceFromTerraformTypeList(tfAction["excluded_environments"]),
		Channels:             getSliceFromTerraformTypeList(tfAction["channels"]),
		TenantTags:           getSliceFromTerraformTypeList(tfAction["tenant_tags"]),
	}

	workerPoolId := tfAction["worker_pool_id"]
	if workerPoolId != nil {
		action.WorkerPoolId = workerPoolId.(string)
	}

	if primaryPackage, ok := tfAction["primary_package"]; ok {
		tfPrimaryPackage := primaryPackage.([]interface {})
		primaryPackage := buildPackageReferenceResource(tfPrimaryPackage[0].(map[string]interface{}))
		action.Packages = append(action.Packages, primaryPackage)
	}

	if pkgAttr, ok := tfAction["package"]; ok {
		tfPkgs := pkgAttr.([]interface {})

		for _, tfPkg := range tfPkgs {
			pkg := buildPackageReferenceResource(tfPkg.(map[string]interface{}))
			action.Packages = append(action.Packages, pkg)
		}
	}

	if propAttr, ok := tfAction["property"]; ok {
		tfProps := propAttr.([]interface {})

		for _, tfProp := range tfProps {
			tfPropi := tfProp.(map[string]interface{})
			action.Properties[tfPropi["key"].(string)] = tfPropi["value"].(string)
		}
	}

	return action;
}


func buildPackageReferenceResource(tfPkg map[string]interface{}) octopusdeploy.PackageReference {
	pkg := octopusdeploy.PackageReference{
		PackageId:           tfPkg["package_id"].(string),
		FeedId:              tfPkg["feed_id"].(string),
		AcquisitionLocation: tfPkg["feed_id"].(string),
	}

	name := tfPkg["name"]
	if name != nil {
		pkg.Name = name.(string)
	}


	extract := tfPkg["extract_during_deployment"]
	if extract != nil {
		pkg.Properties["Extract"] = extract.(string)
	}

	if propAttr, ok := tfPkg["property"]; ok {
		tfProps := propAttr.([]interface {})

		for _, tfProp := range tfProps {
			tfPropi := tfProp.(map[string]interface{})
			pkg.Properties[tfPropi["key"].(string)] = tfPropi["value"].(string)
		}
	}

	return pkg;
}


func resourceDeploymentProcessRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	deploymentProcessID := d.Id()

	_, err := client.DeploymentProcess.Get(deploymentProcessID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading deployment process id %s: %s", deploymentProcessID, err.Error())
	}

	log.Printf("[DEBUG] deploymentProcess: %v", m)

	return nil
}


func resourceDeploymentProcessUpdate(d *schema.ResourceData, m interface{}) error {
	deploymentProcess := buildDeploymentProcessResource(d)
	deploymentProcess.ID = d.Id() // set deploymentProcess struct ID so octopus knows which deploymentProcess to update

	client := m.(*octopusdeploy.Client)

	deploymentProcess, err := client.DeploymentProcess.Update(deploymentProcess)

	if err != nil {
		return fmt.Errorf("error updating deployment process id %s: %s", d.Id(), err.Error())
	}

	d.SetId(deploymentProcess.ID)

	return nil
}

func resourceDeploymentProcessDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	deploymentProcess := &octopusdeploy.DeploymentProcess{
		ID: d.Id(),
	}
	deploymentProcess, err := client.DeploymentProcess.Update(deploymentProcess)

	if err != nil {
		return fmt.Errorf("error deleting deployment process with id %s: %s", deploymentProcess.ID, err.Error())
	}

	d.SetId("")
	return nil
}
