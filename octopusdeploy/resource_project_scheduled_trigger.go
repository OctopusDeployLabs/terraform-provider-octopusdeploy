package octopusdeploy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceProjectScheduledTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectDeploymentTargetTriggerCreate,
		DeleteContext: resourceProjectDeploymentTargetTriggerDelete,
		Importer:      getImporter(),
		ReadContext:   resourceProjectDeploymentTargetTriggerRead,
		Schema:        getProjectDeploymentTargetTriggerSchema(),
		UpdateContext: resourceProjectDeploymentTargetTriggerUpdate,
	}
}

func resourceProjectScheduledTriggerRead()
func resourceProjectScheduledTriggerCreate()
func resourceProjectScheduledTriggerUpdate()
func resourceProjectScheduledTriggerDelete()
