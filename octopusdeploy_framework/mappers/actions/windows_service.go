package actions

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsServiceActionMapper struct{}

func (w WindowsServiceActionMapper) ToState(ctx context.Context, actionState attr.Value, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	diag := mapBaseDeploymentActionToState(ctx, actionState, action, newAction)
	if diag.HasError() {
		return diag
	}

	mapWindowsServicePropertiesToState(action, newAction)

	return nil
}

func (w WindowsServiceActionMapper) ToDeploymentAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := GetActionAttributes(actionAttribute)
	if actionAttrs == nil {
		return nil
	}

	action := GetBaseAction(actionAttribute)
	if action == nil {
		return nil
	}

	mapWindowsServiceProperties(action, actionAttrs)
	return action
}

func mapWindowsServiceProperties(action *deployments.DeploymentAction, actionAttrs map[string]attr.Value) {
	ensureFeatureIsEnabled(action, "Octopus.Features.WindowsService")
	mapBooleanAttributeToProperty(action, actionAttrs, "create_or_update_service", "Octopus.Action.WindowsService.CreateOrUpdateService")

	mapAttributeToProperty(action, actionAttrs, "service_name", "Octopus.Action.WindowsService.ServiceName")
	mapAttributeToProperty(action, actionAttrs, "display_name", "Octopus.Action.WindowsService.DisplayName")
	mapAttributeToProperty(action, actionAttrs, "description", "Octopus.Action.WindowsService.Description")
	mapAttributeToProperty(action, actionAttrs, "executable_path", "Octopus.Action.WindowsService.ExecutablePath")
	mapAttributeToProperty(action, actionAttrs, "arguments", "Octopus.Action.WindowsService.Arguments")
	mapAttributeToProperty(action, actionAttrs, "service_account", "Octopus.Action.WindowsService.ServiceAccount")
	mapAttributeToProperty(action, actionAttrs, "custom_account_name", "Octopus.Action.WindowsService.CustomAccountName")
	mapAttributeToProperty(action, actionAttrs, "custom_account_password", "Octopus.Action.WindowsService.CustomAccountPassword")
	mapAttributeToProperty(action, actionAttrs, "start_mode", "Octopus.Action.WindowsService.StartMode")
	mapAttributeToProperty(action, actionAttrs, "dependencies", "Octopus.Action.WindowsService.Dependencies")
}

func mapWindowsServicePropertiesToState(action *deployments.DeploymentAction, newAction map[string]attr.Value) {
	mapPropertyToStateString(action, newAction, "Octopus.Action.WindowsService.Arguments", "arguments")
	mapPropertyToStateBool(action, newAction, "Octopus.Action.WindowsService.CreateOrUpdateService", "create_or_update_service", false)
	mapPropertyToStateString(action, newAction, "Octopus.Action.WindowsService.CustomAccountName", "custom_account_name")
	mapPropertyToStateString(action, newAction, "Octopus.Action.WindowsService.CustomAccountPassword", "custom_account_password")
	mapPropertyToStateString(action, newAction, "Octopus.Action.WindowsService.Dependencies", "dependencies")
	mapPropertyToStateString(action, newAction, "Octopus.Action.WindowsService.Description", "description")
	mapPropertyToStateString(action, newAction, "Octopus.Action.WindowsService.DisplayName", "display_name")
	mapPropertyToStateString(action, newAction, "Octopus.Action.WindowsService.ExecutablePath", "executable_path")
	mapPropertyToStateString(action, newAction, "Octopus.Action.WindowsService.ServiceAccount", "service_account")
	mapPropertyToStateString(action, newAction, "Octopus.Action.WindowsService.ServiceName", "service_name")
	mapPropertyToStateString(action, newAction, "Octopus.Action.WindowsService.StartMode", "start_mode")
}

func getWindowsServiceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"arguments":                types.StringType,
		"create_or_update_service": types.BoolType,
		"custom_account_name":      types.StringType,
		"custom_account_password":  types.StringType,
		"dependencies":             types.StringType,
		"display_name":             types.StringType,
		"executable_path":          types.StringType,
		"service_account":          types.StringType,
		"service_name":             types.StringType,
		"start_mode":               types.StringType,
	}
}

var _ MappableAction = &WindowsServiceActionMapper{}
