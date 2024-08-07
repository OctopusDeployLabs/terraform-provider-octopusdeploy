package actions

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Action struct{}

func (a Action) ToState(ctx context.Context, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	diags := mapBaseDeploymentActionToState(ctx, action, newAction)
	if diags.HasError() {
		return diags
	}

	newAction["action_type"] = types.StringValue(action.ActionType)
	return nil
}

func (a Action) ToDeploymentAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	return GetBaseAction(actionAttribute)
}

var _ MappableAction = &Action{}
