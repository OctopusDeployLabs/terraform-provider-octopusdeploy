package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type channelResource struct {
	*Config
}

func NewChannelResource() resource.Resource {
	return &channelResource{}
}

var _ resource.ResourceWithImportState = &channelResource{}

func (r *channelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("channel")
}

func (r *channelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ChannelSchema{}.GetResourceSchema()
}

func (r *channelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *channelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.ChannelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating channel", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	channel := expandChannel(ctx, plan)
	createdChannel, err := channels.Add(r.Config.Client, channel)
	if err != nil {
		resp.Diagnostics.AddError("Error creating channel", err.Error())
		return
	}

	state := flattenChannel(ctx, createdChannel, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *channelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.ChannelModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channel, err := channels.GetByID(r.Client, state.SpaceId.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "channelResource"); err != nil {
			resp.Diagnostics.AddError("unable to load channel", err.Error())
		}
		return
	}

	newState := flattenChannel(ctx, channel, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
	return
}

func (r *channelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.ChannelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channel := expandChannel(ctx, plan)
	updatedChannel, err := channels.Update(r.Client, channel)
	if err != nil {
		resp.Diagnostics.AddError("Error updating channel", err.Error())
		return
	}

	state := flattenChannel(ctx, updatedChannel, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	return
}

func (r *channelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.ChannelModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := channels.DeleteByID(r.Client, state.SpaceId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting channel", err.Error())
		return
	}
}

func (*channelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func expandChannel(ctx context.Context, model schemas.ChannelModel) *channels.Channel {
	var channelName = model.Name.ValueString()
	var projectId = model.ProjectId.ValueString()

	channel := channels.NewChannel(channelName, projectId)

	channel.ID = model.ID.ValueString()
	channel.Description = model.Description.ValueString()
	channel.IsDefault = model.IsDefault.ValueBool()
	channel.LifecycleID = model.LifecycleId.ValueString()
	channel.Rules = expandChannelRules(model.Rule)
	channel.SpaceID = model.SpaceId.ValueString()
	channel.TenantTags = expandStringList(model.TenantTags)

	return channel
}

func expandChannelRules(rules types.List) []channels.ChannelRule {
	if rules.IsNull() || rules.IsUnknown() {
		return nil
	}

	var rulesMap []map[string]interface{}
	rules.ElementsAs(context.Background(), &rulesMap, false)

	var channelRules []channels.ChannelRule
	for _, ruleMap := range rulesMap {
		channelRule := expandChannelRule(ruleMap)
		channelRules = append(channelRules, channelRule)
	}

	return nil
}

func expandChannelRule(rule map[string]interface{}) channels.ChannelRule {
	var channelRule channels.ChannelRule

	channelRule.ID = rule["id"].(string)
	actionPackage := rule["action_package"].([]map[string]interface{})
	channelRule.ActionPackages = expandChannelRuleDeploymentActionPackages(actionPackage)
	channelRule.Tag = rule["tag"].(string)
	channelRule.VersionRange = rule["version_range"].(string)

	return channelRule
}

func expandChannelRuleDeploymentActionPackages(actionPackages []map[string]interface{}) []packages.DeploymentActionPackage {
	var actionPackagesExpanded []packages.DeploymentActionPackage
	for _, actionPackage := range actionPackages {
		actionPackageExpanded := expandChannelRuleDeploymentActionPackage(actionPackage)
		actionPackagesExpanded = append(actionPackagesExpanded, actionPackageExpanded)
	}
	return actionPackagesExpanded
}

func expandChannelRuleDeploymentActionPackage(actionPackageMap map[string]interface{}) packages.DeploymentActionPackage {
	return packages.DeploymentActionPackage{
		DeploymentAction: actionPackageMap["deployment_action"].(string),
		PackageReference: actionPackageMap["package_reference"].(string),
	}

}

func flattenChannel(ctx context.Context, channel *channels.Channel, model schemas.ChannelModel) schemas.ChannelModel {
	model.ID = types.StringValue(channel.GetID())
	model.Description = types.StringValue(channel.Description)

	if !channel.IsDefault && model.IsDefault.IsNull() {
		model.IsDefault = types.BoolNull()
	} else {
		model.IsDefault = types.BoolValue(channel.IsDefault)
	}

	if channel.LifecycleID == "" && model.LifecycleId.IsNull() {
		model.LifecycleId = types.StringNull()
	} else {
		model.LifecycleId = types.StringValue(channel.LifecycleID)
	}

	model.Name = types.StringValue(channel.Name)
	model.ProjectId = types.StringValue(channel.ProjectID)

	model.Rule = flattenChannelRules(channel.Rules, model.Rule)

	if channel.SpaceID == "" && model.SpaceId.IsNull() {
		model.SpaceId = types.StringNull()
	} else {
		model.SpaceId = types.StringValue(channel.SpaceID)
	}

	model.TenantTags = flattenStringList(channel.TenantTags, model.TenantTags)

	return model
}

func flattenChannelRules(rules []channels.ChannelRule, currentRules types.List) types.List {
	if len(rules) == 0 && currentRules.IsNull() {
		return types.ListNull(types.ObjectType{AttrTypes: getChannelRuleAttrTypes()})
	}
	if rules == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getChannelRuleAttrTypes()})
	}

	var flattenedRules = make([]attr.Value, len(rules))
	for _, rule := range rules {
		obj := flattenChannelRule(&rule)
		flattenedRules = append(flattenedRules, obj)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getChannelRuleAttrTypes()}, flattenedRules)
}

func flattenChannelRule(rule *channels.ChannelRule) types.Object {
	return types.ObjectValueMust(getChannelRuleAttrTypes(), map[string]attr.Value{
		"action_package": flattenChannelRuleDeploymentActionPackages(rule.ActionPackages),
		"id":             types.StringValue(rule.ID),
		"tag":            types.StringValue(rule.Tag),
		"version_range":  types.StringValue(rule.VersionRange),
	})

}

func flattenChannelRuleDeploymentActionPackages(actionPackages []packages.DeploymentActionPackage) types.List {
	var flattenedActionPackages = make([]attr.Value, len(actionPackages))
	for _, actionPackage := range actionPackages {
		obj := flattenChannelRuleDeploymentActionPackage(&actionPackage)
		flattenedActionPackages = append(flattenedActionPackages, obj)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getChannelRuleDeploymentActionPackageAttrTypes()}, flattenedActionPackages)
}

func flattenChannelRuleDeploymentActionPackage(actionPackage *packages.DeploymentActionPackage) types.Object {
	return types.ObjectValueMust(getChannelRuleDeploymentActionPackageAttrTypes(), map[string]attr.Value{
		"deployment_action": types.StringValue(actionPackage.DeploymentAction),
		"package_reference": types.StringValue(actionPackage.PackageReference),
	})
}

func getChannelRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"action_package": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: getChannelRuleDeploymentActionPackageAttrTypes(),
			},
		},
		"id":            types.StringType,
		"tag":           types.StringType,
		"version_range": types.StringType,
	}
}

func getChannelRuleDeploymentActionPackageAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"deployment_action": types.StringType,
		"package_reference": types.StringType,
	}
}
