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
	mapPropertyToStateBool(action, newAction, "Octopus.Action.RunOnServer", "run_on_server", false)
	newAction["worker_pool_id"] = types.StringValue(action.WorkerPool)
	newAction["worker_pool_variable"] = types.StringValue(action.WorkerPoolVariable)
	return nil
}

func (a Action) ToDeploymentAction(actionAttribute attr.Value) *deployments.DeploymentAction {

	return GetBaseAction(actionAttribute)

}

var _ MappableAction = &Action{}
