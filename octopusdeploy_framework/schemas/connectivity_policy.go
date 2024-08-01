package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var runbookConnectivityPolicySchemeAttributeNames = struct {
	AllowDeploymentsToNoTargets string
	ExcludeUnhealthyTargets     string
	SkipMachineBehavior         string
	TargetRoles                 string
}{
	AllowDeploymentsToNoTargets: "allow_deployments_to_no_targets",
	ExcludeUnhealthyTargets:     "exclude_unhealthy_targets",
	SkipMachineBehavior:         "skip_machine_behavior",
	TargetRoles:                 "target_roles",
}

var skipMachineBehaviorNames = struct {
	SkipUnavailableMachines string
	None                    string
}{
	SkipUnavailableMachines: "SkipUnavailableMachines",
	None:                    "None",
}

var skipMachineBehaviors = []string{
	skipMachineBehaviorNames.SkipUnavailableMachines,
	skipMachineBehaviorNames.None,
}

func GetConnectivityPolicyObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		runbookConnectivityPolicySchemeAttributeNames.AllowDeploymentsToNoTargets: types.BoolType,
		runbookConnectivityPolicySchemeAttributeNames.ExcludeUnhealthyTargets:     types.BoolType,
		runbookConnectivityPolicySchemeAttributeNames.SkipMachineBehavior:         types.StringType,
		runbookConnectivityPolicySchemeAttributeNames.TargetRoles:                 types.ListType{ElemType: types.StringType},
	}
}

func getConnectivityPolicySchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		runbookConnectivityPolicySchemeAttributeNames.AllowDeploymentsToNoTargets: resourceSchema.BoolAttribute{
			Computed: true,
			Optional: true,
		},
		runbookConnectivityPolicySchemeAttributeNames.ExcludeUnhealthyTargets: resourceSchema.BoolAttribute{
			Computed: true,
			Optional: true,
		},
		runbookConnectivityPolicySchemeAttributeNames.SkipMachineBehavior: resourceSchema.StringAttribute{
			Optional: true,
			Default:  stringdefault.StaticString(skipMachineBehaviorNames.None),
			Validators: []validator.String{
				stringvalidator.OneOf(
					skipMachineBehaviors...,
				),
			},
		},
		runbookConnectivityPolicySchemeAttributeNames.TargetRoles: resourceSchema.ListAttribute{
			Computed:    true,
			Optional:    true,
			ElementType: types.StringType,
		},
	}
}

func MapFromConnectivityPolicy(connectivityPolicy *core.ConnectivityPolicy) attr.Value {
	if connectivityPolicy == nil {
		return nil
	}

	attrs := map[string]attr.Value{
		runbookConnectivityPolicySchemeAttributeNames.AllowDeploymentsToNoTargets: types.BoolValue(connectivityPolicy.AllowDeploymentsToNoTargets),
		runbookConnectivityPolicySchemeAttributeNames.ExcludeUnhealthyTargets:     types.BoolValue(connectivityPolicy.ExcludeUnhealthyTargets),
		runbookConnectivityPolicySchemeAttributeNames.SkipMachineBehavior:         types.StringValue(string(connectivityPolicy.SkipMachineBehavior)),
		runbookConnectivityPolicySchemeAttributeNames.TargetRoles:                 util.FlattenStringList(connectivityPolicy.TargetRoles),
	}

	return types.ObjectValueMust(GetConnectivityPolicyObjectType(), attrs)
}

func MapToConnectivityPolicy(flattenedConnectivityPolicy types.List) *core.ConnectivityPolicy {
	if flattenedConnectivityPolicy.IsNull() || len(flattenedConnectivityPolicy.Elements()) == 0 {
		return nil
	}

	obj := flattenedConnectivityPolicy.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	var connectivityPolicy core.ConnectivityPolicy
	if allowDeploymentsToNoTargets, ok := attrs[runbookConnectivityPolicySchemeAttributeNames.AllowDeploymentsToNoTargets].(types.Bool); ok && !allowDeploymentsToNoTargets.IsNull() {
		connectivityPolicy.AllowDeploymentsToNoTargets = allowDeploymentsToNoTargets.ValueBool()
	}
	if excludeUnhealthyTargets, ok := attrs[runbookConnectivityPolicySchemeAttributeNames.ExcludeUnhealthyTargets].(types.Bool); ok && !excludeUnhealthyTargets.IsNull() {
		connectivityPolicy.ExcludeUnhealthyTargets = excludeUnhealthyTargets.ValueBool()
	}
	if skipMachineBehavior, ok := attrs[runbookConnectivityPolicySchemeAttributeNames.SkipMachineBehavior].(types.String); ok && !skipMachineBehavior.IsNull() {
		connectivityPolicy.SkipMachineBehavior = core.SkipMachineBehavior(skipMachineBehavior.ValueString())
	}
	if targetRoles, ok := attrs[runbookConnectivityPolicySchemeAttributeNames.TargetRoles].(types.List); ok && !targetRoles.IsNull() {
		connectivityPolicy.TargetRoles = util.ExpandStringList(targetRoles)
	}

	return &connectivityPolicy
}
