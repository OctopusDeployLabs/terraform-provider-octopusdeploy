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
	resourceDeploymentStepAddDefaultSchema(schemaRes, true)
	resourceDeploymentStepAddPackageSchema(schemaRes)

	/* Return Schema */
	return schemaRes
}

func buildDeployPackageDeploymentStep(d *schema.ResourceData) *octopusdeploy.DeploymentStep {
	/* Create Basic Deployment Step */
	deploymentStep := resourceDeploymentStepCreateBasicStep(d, "Octopus.TentaclePackage")

	/* Add Shared Properties */
	resourceDeploymentStepAddPackageProperties(d, deploymentStep)

	/* Return Deployment Step */
	return deploymentStep
}

func setDeployPackageSchema(d *schema.ResourceData, deploymentStep octopusdeploy.DeploymentStep) {
	resourceDeploymentStepSetBasicSchema(d, deploymentStep)
	resourceDeploymentStepSetPackageSchema(d, deploymentStep)
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
