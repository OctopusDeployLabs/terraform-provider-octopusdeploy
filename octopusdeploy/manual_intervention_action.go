package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func getManualInterventionActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()

	element.Schema[constInstructions] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The instructions for the user to follow",
		Required:    true,
	}

	element.Schema[constResponsibleTeams] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The teams responsible to resolve this step. If no teams are specified, all users who have permission to deploy the project can resolve it.",
		Optional:    true,
	}

	return actionSchema
}

func buildManualInterventionActionResource(tfAction map[string]interface{}) model.DeploymentAction {
	resource := buildDeploymentActionResource(tfAction)
	resource.ActionType = "Octopus.Manual"
	resource.Properties["Octopus.Action.Manual.Instructions"] = tfAction[constInstructions].(string)

	responsibleTeams := tfAction[constResponsibleTeams]
	if responsibleTeams != nil {
		resource.Properties["Octopus.Action.Manual.ResponsibleTeamIds"] = responsibleTeams.(string)
	}

	return resource
}
