package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDeploymentStepDeployPackage() *schema.Resource {
	schemaRes := &schema.Resource{
		Create: resourceDeploymentStepDeployPackageCreate,
		Read:   resourceDeploymentStepDeployPackageRead,
		Update: resourceDeploymentStepDeployPackageUpdate,
		Delete: resourceDeploymentStepDeployPackageDelete,

		Schema: map[string]*schema.Schema{},
	}

	/* Add Shared Schema's */
	resourceDeploymentStep_AddDefaultSchema(schemaRes, true)
	resourceDeploymentStep_AddPackageSchema(schemaRes)

	/* Return Schema */
	return schemaRes
}

func buildDeployPackageDeploymentStep(d *schema.ResourceData) *octopusdeploy.DeploymentStep {
	/* Create Basic Deployment Step */
	deploymentStep := resourceDeploymentStep_CreateBasicStep(d, "Octopus.TentaclePackage")	

	/* Add Shared Properties */
	resourceDeploymentStep_AddPackageProperties(d, deploymentStep)

	/* Return Deployment Step */
	return deploymentStep
}

func setDeployPackageSchema(d *schema.ResourceData, deploymentStep octopusdeploy.DeploymentStep) {
	resourceDeploymentStep_SetBasicSchema(d, deploymentStep);
	resourceDeploymentStep_SetPackageSchema(d, deploymentStep);
}

func resourceDeploymentStepDeployPackageCreate(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepCreate(d, m, buildDeployPackageDeploymentStep)
}

func resourceDeploymentStepDeployPackageRead(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepRead(d, m, setDeployPackageSchema)
}

func resourceDeploymentStepDeployPackageUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepUpdate(d, m, buildDeployPackageDeploymentStep)
}

func resourceDeploymentStepDeployPackageDelete(d *schema.ResourceData, m interface{}) error {
	return resourceDeploymentStepDelete(d, m)
}
