package deployment_process

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type ManualInterventionActionMapper struct{}

var _ MappableAction = &ManualInterventionActionMapper{}

func (m *ManualInterventionActionMapper) ToState(ctx context.Context, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	diag := mapBaseDeploymentActionToState(ctx, action, newAction)
	if diag.HasError() {
		return diag
	}

	mapPropertyToStateString(action, newAction, "Octopus.Action.Manual.Instructions", "instructions")
	mapPropertyToStateString(action, newAction, "Octopus.Action.Manual.ResponsibleTeamIds", "responsible_teams")

	return nil
}

func (m *ManualInterventionActionMapper) ToDeploymentAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := GetActionAttributes(actionAttribute)
	if actionAttrs == nil {
		return nil
	}

	action := GetBaseAction(actionAttribute)
	if action == nil {
		return nil
	}

	mapAttributeToProperty(action, actionAttrs, "instructions", "Octopus.Action.Manual.Instructions")
	mapAttributeToProperty(action, actionAttrs, "responsible_teams", "Octopus.Action.Manual.ResponsibleTeamIds")

	return action
}
