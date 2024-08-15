package actions

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type PackageActionMapper struct{}

func (p PackageActionMapper) ToState(ctx context.Context, actionState attr.Value, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	diag := mapBaseDeploymentActionToState(ctx, actionState, action, newAction)
	if diag.HasError() {
		return diag
	}

	if v, ok := action.Properties["Octopus.Action.EnabledFeatures"]; ok {
		if strings.Contains(v.Value, "Octopus.Features.WindowsService") {
			attrs := make(map[string]attr.Value)
			mapWindowsServicePropertiesToState(action, attrs)
			list := make([]attr.Value, 1)
			list[0] = types.ObjectValueMust(getWindowsServiceAttrTypes(), attrs)
			newAction["windows_service"] = types.ListValueMust(types.ObjectType{AttrTypes: getWindowsServiceAttrTypes()}, list)
		} else {
			newAction["windows_service"] = types.ListNull(types.ObjectType{AttrTypes: getWindowsServiceAttrTypes()})
		}
	}

	return nil
}

func (p PackageActionMapper) ToDeploymentAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := GetActionAttributes(actionAttribute)
	if actionAttrs == nil {
		return nil
	}

	action := GetBaseAction(actionAttribute)
	if action == nil {
		return nil
	}

	action.ActionType = "Octopus.TentaclePackage"

	if v, ok := actionAttrs["windows_service"]; ok {
		list := v.(types.List).Elements()
		for _, item := range list {
			mapWindowsServiceProperties(action, item.(types.Object).Attributes())
			return action
		}
	}

	return action
}

var _ MappableAction = &PackageActionMapper{}
