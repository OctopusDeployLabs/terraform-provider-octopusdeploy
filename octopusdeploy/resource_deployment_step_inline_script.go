package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDeploymentStepInlineScript() *schema.Resource {
	schemaRes := &schema.Resource{
		Create: resourceDeploymentStepInlineScriptCreate,
		Read:   resourceDeploymentStepInlineScriptRead,
		Update: resourceDeploymentStepInlineScriptUpdate,
		Delete: resourceDeploymentStepInlineScriptDelete,

		Schema: map[string]*schema.Schema{
			"script_type": {
				Type:        schema.TypeString,
				Description: "The scripting language of the deployment step.",
				Required:    true,
				ValidateFunc: validateValueFunc([]string{
					"PowerShell",
					"CSharp",
					"Bash",
					"FSharp",
				}),
			},
			"script_body": {
				Type:        schema.TypeString,
				Description: "The script body.",
				Required:    true,
			},
			"script_source": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	/* Add Shared Schema's */
	resourceDeploymentStepAddDefaultSchema(schemaRes, false)

	/* Return Schema */
	return schemaRes
}

func buildInlineScriptDeploymentStep(d *schema.ResourceData) *octopusdeploy.DeploymentStep {
	/* Set Computed Values */
	d.Set("script_source", "Inline")

	/* Create Basic Deployment Step */
	deploymentStep := resourceDeploymentStepCreateBasicStep(d, "Octopus.Script")

	/* Add Script Properties */
	deploymentStep.Actions[0].Properties["Octopus.Action.Script.ScriptSource"] = d.Get("script_source").(string)
	deploymentStep.Actions[0].Properties["Octopus.Action.Script.Syntax"] = d.Get("script_type").(string)
	deploymentStep.Actions[0].Properties["Octopus.Action.Script.ScriptBody"] = d.Get("script_body").(string)

	/* Return Deployment Step */
	return deploymentStep
}

func setInlineScriptSchema(d *schema.ResourceData, deploymentStep octopusdeploy.DeploymentStep) {
	resourceDeploymentStepSetBasicSchema(d, deploymentStep)

	/* Get Script Properties */
	d.Set("script_source", deploymentStep.Actions[0].Properties["Octopus.Action.Script.ScriptSource"])
	d.Set("script_type", deploymentStep.Actions[0].Properties["Octopus.Action.Script.Syntax"])
	d.Set("script_body", deploymentStep.Actions[0].Properties["Octopus.Action.Script.ScriptBody"])
}

func resourceDeploymentStepInlineScriptCreate(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepCreate(d, m, buildInlineScriptDeploymentStep)
}

func resourceDeploymentStepInlineScriptRead(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepRead(d, m, setInlineScriptSchema)
}

func resourceDeploymentStepInlineScriptUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepUpdate(d, m, buildInlineScriptDeploymentStep)
}

func resourceDeploymentStepInlineScriptDelete(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepDelete(d, m)
}
