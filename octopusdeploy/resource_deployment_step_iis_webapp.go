package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDeploymentStepIisWebapp() *schema.Resource {
	schemaRes := &schema.Resource{
		Create: resourceDeploymentStepIisWebappCreate,
		Read:   resourceDeploymentStepIisWebappRead,
		Update: resourceDeploymentStepIisWebappUpdate,
		Delete: resourceDeploymentStepIisWebappDelete,

		Schema: map[string]*schema.Schema{
			"deployment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"website_name": {
				Type:        schema.TypeString,
				Description: "The name of the Website to be add web application to",
				Required:    true,
			},
			"virtual_path": {
				Type:        schema.TypeString,
				Description: "Virtual Path for the Web Application",
				Required:    true,
			},
			"path_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"relative_path": {
				Type:        schema.TypeString,
				Description: "Relative Path to package Root for the physical Path",
				Optional:    true,
			},
		},
	}

	/* Add Shared Schema's */
	resourceDeploymentStep_AddDefaultSchema(schemaRes, true)
	resourceDeploymentStep_AddPackageSchema(schemaRes)
	resourceDeploymentStep_AddIisAppPoolSchema(schemaRes)

	/* Return Schema */
	return schemaRes
}

func buildIisWebappDeploymentStep(d *schema.ResourceData) *octopusdeploy.DeploymentStep {
	/* Set Computed Values */
	d.Set("deployment_type", "webApplication")

	/* Create Basic Deployment Step */
	deploymentStep := resourceDeploymentStep_CreateBasicStep(d, "Octopus.IIS")

	/* Enable IIS Web Site Features */
	deploymentStep.Actions[0].Properties["Octopus.Action.EnabledFeatures"] = "Octopus.Features.IISWebSite"

	/* Add Shared Properties */
	resourceDeploymentStep_AddPackageProperties(d, deploymentStep)
	resourceDeploymentStep_AddIisAppPoolProperties(d, deploymentStep, "WebApplication")

	/* Add Web Site Properties */
	deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.DeploymentType"] = d.Get("deployment_type").(string)
	deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.StartWebSite"] = "true"
	deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.CreateOrUpdateWebSite"] = "False"
	deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.WebApplication.CreateOrUpdate"] = "True"
	deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.VirtualDirectory.CreateOrUpdate"] = "False"

	if relativePath, ok := d.GetOk("relative_path"); ok {
		d.Set("path_type", "relativeToPackageRoot")
		deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.WebApplication.PhysicalPath"] = relativePath.(string)
	} else {
		d.Set("path_type", "packageRoot")
	}
	deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.WebRootType"] = d.Get("path_type").(string)
	deploymentStep.Actions[0].Properties["Octopus.Action.WebApplication.WebRootType"] = d.Get("path_type").(string)

	deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.WebApplication.WebSiteName"] = d.Get("website_name").(string)
	deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.WebApplication.VirtualPath"] = d.Get("virtual_path").(string)

	/* Return Deployment Step */
	return deploymentStep
}

func setIisWebappSchema(d *schema.ResourceData, deploymentStep octopusdeploy.DeploymentStep) {
	resourceDeploymentStep_SetBasicSchema(d, deploymentStep)
	resourceDeploymentStep_SetPackageSchema(d, deploymentStep)
	resourceDeploymentStep_SetIisAppPoolSchema(d, deploymentStep, "WebApplication")

	/* Get Web Site Properties */
	d.Set("deployment_type", deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.DeploymentType"])

	if pathType, ok := deploymentStep.Actions[0].Properties["Octopus.Action.WebApplication.WebRootType"]; ok {
		d.Set("path_type", pathType)
	}

	if relativePath, ok := deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.WebApplication.PhysicalPath"]; ok {
		d.Set("relative_path", relativePath)
	}

	if websiteName, ok := deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.WebApplication.WebSiteName"]; ok {
		d.Set("website_name", websiteName)
	}

	if virtualPath, ok := deploymentStep.Actions[0].Properties["Octopus.Action.IISWebSite.WebApplication.VirtualPath"]; ok {
		d.Set("virtual_path", virtualPath)
	}
}

func resourceDeploymentStepIisWebappCreate(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepCreate(d, m, buildIisWebappDeploymentStep)
}

func resourceDeploymentStepIisWebappRead(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepRead(d, m, setIisWebappSchema)
}

func resourceDeploymentStepIisWebappUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepUpdate(d, m, buildIisWebappDeploymentStep)
}

func resourceDeploymentStepIisWebappDelete(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepDelete(d, m)
}
