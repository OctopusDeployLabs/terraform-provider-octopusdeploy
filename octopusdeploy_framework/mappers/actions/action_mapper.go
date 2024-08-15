package actions

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type MappableAction interface {
	ToState(ctx context.Context, actionState attr.Value, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics
	ToDeploymentAction(actionAttribute attr.Value) *deployments.DeploymentAction
}
