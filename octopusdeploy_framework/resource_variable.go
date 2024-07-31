package octopusdeploy_framework

import (
	"context"
	"fmt"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type variableTypeResource struct {
	*Config
}

var _ resource.ResourceWithImportState = &variableTypeResource{}

func NewVariableResource() resource.Resource {
	return &variableTypeResource{}
}

func (r *variableTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, fmt.Sprintf("importing variable (%s)", req.ID))

	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"unexpected import identifier",
			fmt.Sprintf("%s_variable import must be in the form of OwnerID:VariableID (e.g. Projects-62:0906031f-68ba-4a15-afaa-657c1564e07b)", util.GetProviderName()))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *variableTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetProviderName() + "_variable"
}

func (r *variableTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetVariableResourceSchema()
}

func (r *variableTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *variableTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	mutex.Lock()

	var data schemas.VariableTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	variableOwnerId, err := getVariableOwnerID(&data)
	if err != nil {
		resp.Diagnostics.AddError("invalid resource configuration", err.Error())
		return
	}

	name := data.Name.ValueString()
	newVariable := variables.NewVariable(name)
	newVariable.Description = data.Description.ValueString()
	newVariable.IsEditable = data.IsEditable.ValueBool()
	newVariable.IsSensitive = data.IsSensitive.ValueBool()
	newVariable.Type = data.Type.ValueString()
	newVariable.Scope = expandVariableScopes(data.Scope)
	newVariable.Prompt = expandPromptedVariableSettings(data.Prompt)
	newVariable.SpaceID = data.SpaceID.ValueString()

	if newVariable.IsSensitive {
		newVariable.Type = schemas.VariableTypeNames.Sensitive
		newVariable.Value = data.SensitiveValue.ValueString()
	} else {
		newVariable.Value = data.Value.ValueString()
	}

	tflog.Info(ctx, fmt.Sprintf("creating variable: %#v", newVariable))

	variableSet, err := variables.AddSingle(r.Config.Client, r.Config.SpaceID, variableOwnerId.ValueString(), newVariable)
	if err != nil {
		resp.Diagnostics.AddError("create variable failed", err.Error())
		return
	}

	err = validateVariable(&variableSet, newVariable, variableOwnerId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("create variable failed", err.Error())
		return
	}

	mapVariableToState(&data, newVariable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	mutex.Unlock()
}

func (r *variableTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	mutex.Lock()

	var data schemas.VariableTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading variable (%s)", data.ID))

	variableOwnerID, err := getVariableOwnerID(&data)
	if err != nil {
		resp.Diagnostics.AddError("invalid resource configuration", err.Error())
		return
	}

	variable, err := variables.GetByID(r.Config.Client, data.SpaceID.ValueString(), variableOwnerID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load variable", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("variable read (%s)", data.ID))
	mapVariableToState(&data, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	mutex.Unlock()
}

func (r *variableTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	mutex.Lock()

	var data, state schemas.VariableTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating variable (%s)", data.ID))

	variableOwnerId, err := getVariableOwnerID(&data)
	if err != nil {
		resp.Diagnostics.AddError("invalid resource configuration", err.Error())
		return
	}

	name := data.Name.ValueString()
	updatedVariable := variables.NewVariable(name)
	updatedVariable.Description = data.Description.ValueString()
	updatedVariable.IsEditable = data.IsEditable.ValueBool()
	updatedVariable.IsSensitive = data.IsSensitive.ValueBool()
	updatedVariable.Type = data.Type.ValueString()
	updatedVariable.Scope = expandVariableScopes(data.Scope)
	updatedVariable.Prompt = expandPromptedVariableSettings(data.Prompt)
	updatedVariable.SpaceID = state.SpaceID.ValueString()

	if updatedVariable.IsSensitive {
		updatedVariable.Type = schemas.VariableTypeNames.Sensitive
		updatedVariable.Value = data.SensitiveValue.ValueString()
	} else {
		updatedVariable.Value = data.Value.ValueString()
	}

	updatedVariable.ID = state.ID.ValueString()

	variableSet, err := variables.UpdateSingle(r.Config.Client, state.SpaceID.ValueString(), variableOwnerId.ValueString(), updatedVariable)
	if err != nil {
		resp.Diagnostics.AddError("update variable failed", err.Error())
		return
	}

	err = validateVariable(&variableSet, updatedVariable, variableOwnerId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("update variable failed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("variable updated (%s)", data.ID))

	mapVariableToState(&data, updatedVariable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	mutex.Unlock()
}

func (r *variableTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	mutex.Lock()

	var data schemas.VariableTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("deleting variable (%s)", data.ID.ValueString()))
	variableOwnerID, err := getVariableOwnerID(&data)
	if err != nil {
		resp.Diagnostics.AddError("invalid resource configuration", err.Error())
		return
	}

	if _, err := variables.DeleteSingle(r.Config.Client, r.Config.SpaceID, variableOwnerID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete variable", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("variable deleted (%s)", data.ID))

	mutex.Unlock()
}

func (r *variableTypeResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data schemas.VariableTypeResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	isSensitive := data.IsSensitive
	variableType := data.Type

	if isSensitive.ValueBool() && variableType.ValueString() != schemas.VariableTypeNames.Sensitive {
		resp.Diagnostics.AddError(
			"invalid resource configuration",
			fmt.Sprintf("when %s is set to true, type needs to be '%s'", schemas.VariableSchemaAttributeNames.IsSensitive, schemas.VariableTypeNames.Sensitive),
		)
		return
	}

	if !isSensitive.ValueBool() && variableType.ValueString() == schemas.VariableTypeNames.Sensitive {
		resp.Diagnostics.AddError(
			"invalid resource configuration",
			fmt.Sprintf("when type is set to '%s', %s needs to be true", schemas.VariableSchemaAttributeNames.IsSensitive, schemas.VariableTypeNames.Sensitive),
		)
	}
}

func getVariableOwnerID(data *schemas.VariableTypeResourceModel) (*basetypes.StringValue, error) {
	if data.ProjectID.IsNull() && data.OwnerID.IsNull() {
		return nil, fmt.Errorf("one of %s or %s must be configured", schemas.VariableSchemaAttributeNames.ProjectID, schemas.VariableSchemaAttributeNames.OwnerID)
	} else if !data.ProjectID.IsNull() {
		return &data.ProjectID, nil
	} else {
		return &data.OwnerID, nil
	}
}

func validateVariable(variableSet *variables.VariableSet, newVariable *variables.Variable, variableOwnerId string) error {
	for _, v := range variableSet.Variables {
		if v.Name == newVariable.Name && v.Type == newVariable.Type && (v.IsSensitive || v.Value == newVariable.Value) && v.Description == newVariable.Description && v.IsSensitive == newVariable.IsSensitive {
			scopeMatches, _, err := variables.MatchesScope(v.Scope, &newVariable.Scope)
			if err != nil || !scopeMatches {
				return err
			}
			if scopeMatches {
				newVariable.ID = v.GetID()
				return nil
			}
		}
	}

	return fmt.Errorf("unable to locate variable for owner ID %s", variableOwnerId)
}

func mapVariableToState(data *schemas.VariableTypeResourceModel, variable *variables.Variable) {
	data.SpaceID = types.StringValue(variable.SpaceID)
	data.Name = types.StringValue(variable.Name)
	data.Description = types.StringValue(variable.Description)
	if !data.IsEditable.IsNull() {
		data.IsEditable = types.BoolValue(variable.IsEditable)
	}
	if !data.IsSensitive.IsNull() {
		data.IsSensitive = types.BoolValue(variable.IsSensitive)
	}
	data.Type = types.StringValue(variable.Type)

	if variable.IsSensitive {
		data.Value = types.StringNull()
	} else {
		if !data.Value.IsNull() {
			data.Value = types.StringValue(variable.Value)
		}
	}

	if !data.Prompt.IsNull() {
		data.Prompt = types.ListValueMust(
			types.ObjectType{AttrTypes: variablePromptOptionsObjectType()},
			[]attr.Value{flattenPromptedVariableSettings(variable.Prompt)},
		)
	}

	if !data.Scope.IsNull() {
		data.Scope = types.ListValueMust(
			types.ObjectType{AttrTypes: variableScopeObjectType()},
			[]attr.Value{flattenVariableScopes(variable.Scope)},
		)
	}

	data.EncryptedValue = types.StringNull()
	data.KeyFingerprint = types.StringNull()

	data.ID = types.StringValue(variable.GetID())
}

func variablePromptOptionsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		schemas.SchemaAttributeNames.Description: types.StringType,
		schemas.VariableSchemaAttributeNames.DisplaySettings: types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: variablePromptOptionsDisplaySettingsObjectType(),
			},
		},
		schemas.VariableSchemaAttributeNames.IsRequired: types.BoolType,
		schemas.VariableSchemaAttributeNames.Label:      types.StringType,
	}
}

func variablePromptOptionsDisplaySettingsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		schemas.VariableSchemaAttributeNames.ControlType: types.StringType,
		schemas.VariableSchemaAttributeNames.SelectOption: types.ListType{
			ElemType: variablePromptoOptionsDisplaySettingsSelectOptionObjectType(),
		},
	}
}

func variablePromptoOptionsDisplaySettingsSelectOptionObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemas.VariableSchemaAttributeNames.Value:       types.StringType,
			schemas.VariableSchemaAttributeNames.DisplayName: types.StringType,
		},
	}
}

func flattenPromptedVariableSettings(variablePromptOptions *variables.VariablePromptOptions) attr.Value {
	if variablePromptOptions == nil {
		return nil
	}

	attrs := map[string]attr.Value{
		schemas.SchemaAttributeNames.Description:             types.StringValue(variablePromptOptions.Description),
		schemas.VariableSchemaAttributeNames.IsRequired:      types.BoolValue(variablePromptOptions.IsRequired),
		schemas.VariableSchemaAttributeNames.Label:           types.StringValue(variablePromptOptions.Label),
		schemas.VariableSchemaAttributeNames.DisplaySettings: types.ListNull(types.ObjectType{AttrTypes: variablePromptOptionsDisplaySettingsObjectType()}),
	}
	if variablePromptOptions.DisplaySettings != nil {
		attrs[schemas.VariableSchemaAttributeNames.DisplaySettings] = types.ListValueMust(
			types.ObjectType{
				AttrTypes: variablePromptOptionsDisplaySettingsObjectType(),
			},
			[]attr.Value{
				flattenDisplaySettings(variablePromptOptions.DisplaySettings),
			},
		)
	}

	return types.ObjectValueMust(variablePromptOptionsObjectType(), attrs)
}

func flattenDisplaySettings(displaySettings *resources.DisplaySettings) attr.Value {
	if displaySettings == nil {
		return nil
	}

	attrs := map[string]attr.Value{
		schemas.VariableSchemaAttributeNames.ControlType: types.StringValue(string(displaySettings.ControlType)),
	}
	if displaySettings.ControlType == resources.ControlTypeSelect {
		if len(displaySettings.SelectOptions) > 0 {
			attrs[schemas.VariableSchemaAttributeNames.SelectOption] = types.ListValueMust(
				variablePromptoOptionsDisplaySettingsSelectOptionObjectType(),
				flattenSelectOptions(displaySettings.SelectOptions),
			)
		}
	} else {
		attrs[schemas.VariableSchemaAttributeNames.SelectOption] = types.ListNull(variablePromptoOptionsDisplaySettingsSelectOptionObjectType())
	}

	return types.ObjectValueMust(
		variablePromptOptionsDisplaySettingsObjectType(),
		attrs,
	)
}

func flattenSelectOptions(selectOptions []*resources.SelectOption) []attr.Value {

	options := make([]attr.Value, len(selectOptions))
	for _, option := range selectOptions {
		options = append(options, types.ObjectValueMust(
			variablePromptOptionsDisplaySettingsObjectType(),
			map[string]attr.Value{
				schemas.VariableSchemaAttributeNames.Value:       types.StringValue(option.Value),
				schemas.VariableSchemaAttributeNames.DisplayName: types.StringValue(option.DisplayName),
			},
		))
	}
	return options
}

func expandPromptedVariableSettings(flattenedVariablePromptOptions types.List) *variables.VariablePromptOptions {
	if flattenedVariablePromptOptions.IsNull() {
		return nil
	}

	obj := flattenedVariablePromptOptions.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	var promptOptions variables.VariablePromptOptions
	if description, ok := attrs[schemas.SchemaAttributeNames.Description].(types.String); ok && !description.IsNull() {
		promptOptions.Description = description.ValueString()
	}

	if isRequired, ok := attrs[schemas.VariableSchemaAttributeNames.IsRequired].(types.Bool); ok && !isRequired.IsNull() {
		promptOptions.IsRequired = isRequired.ValueBool()
	}

	if label, ok := attrs[schemas.VariableSchemaAttributeNames.Label].(types.String); ok && !label.IsNull() {
		promptOptions.Label = label.ValueString()
	}

	if displaySettings, ok := attrs[schemas.VariableSchemaAttributeNames.DisplaySettings].(types.List); ok && !displaySettings.IsNull() {
		promptOptions.DisplaySettings = expandDisplaySettings(displaySettings)
	}

	return &promptOptions
}

func expandDisplaySettings(flattenedDisplaySettings types.List) *resources.DisplaySettings {
	if flattenedDisplaySettings.IsNull() {
		return nil
	}

	obj := flattenedDisplaySettings.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	ct, _ := attrs[schemas.VariableSchemaAttributeNames.ControlType].(types.String)
	controlType := resources.ControlType(ct.ValueString())

	var selectOptions []*resources.SelectOption
	if controlType == resources.ControlTypeSelect {
		selectOptions = expandSelectOptions(attrs[schemas.VariableSchemaAttributeNames.SelectOption].(types.List))
	}

	return resources.NewDisplaySettings(controlType, selectOptions)
}

func expandSelectOptions(flattenedSelectOptions types.List) []*resources.SelectOption {
	if flattenedSelectOptions.IsNull() || flattenedSelectOptions.IsUnknown() {
		return nil
	}

	options := make([]*resources.SelectOption, len(flattenedSelectOptions.Elements()))
	for _, option := range flattenedSelectOptions.Elements() {
		attrs := option.(types.Object).Attributes()
		options = append(options, &resources.SelectOption{
			DisplayName: attrs[schemas.VariableSchemaAttributeNames.DisplayName].(types.String).ValueString(),
			Value:       attrs[schemas.VariableSchemaAttributeNames.Value].(types.String).ValueString(),
		})
	}

	return options
}

func variableScopeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		schemas.VariableScopeFieldNames.Actions:      types.ListType{ElemType: types.StringType},
		schemas.VariableScopeFieldNames.Channels:     types.ListType{ElemType: types.StringType},
		schemas.VariableScopeFieldNames.Environments: types.ListType{ElemType: types.StringType},
		schemas.VariableScopeFieldNames.Machines:     types.ListType{ElemType: types.StringType},
		schemas.VariableScopeFieldNames.Processes:    types.ListType{ElemType: types.StringType},
		schemas.VariableScopeFieldNames.Roles:        types.ListType{ElemType: types.StringType},
		schemas.VariableScopeFieldNames.TenantTags:   types.ListType{ElemType: types.StringType},
	}
}

func flattenVariableScopes(variableScopes variables.VariableScope) attr.Value {
	if variableScopes.IsEmpty() {
		return nil
	}

	flattenedScopes := map[string]attr.Value{}
	flattenedScopes[schemas.VariableScopeFieldNames.Actions] = util.Ternary(variableScopes.Actions != nil, util.FlattenStringList(variableScopes.Actions), types.ListNull(types.StringType))
	flattenedScopes[schemas.VariableScopeFieldNames.Channels] = util.Ternary(variableScopes.Channels != nil, util.FlattenStringList(variableScopes.Channels), types.ListNull(types.StringType))
	flattenedScopes[schemas.VariableScopeFieldNames.Environments] = util.Ternary(variableScopes.Environments != nil, util.FlattenStringList(variableScopes.Environments), types.ListNull(types.StringType))
	flattenedScopes[schemas.VariableScopeFieldNames.Machines] = util.Ternary(variableScopes.Machines != nil, util.FlattenStringList(variableScopes.Machines), types.ListNull(types.StringType))
	flattenedScopes[schemas.VariableScopeFieldNames.Processes] = util.Ternary(variableScopes.ProcessOwners != nil, util.FlattenStringList(variableScopes.ProcessOwners), types.ListNull(types.StringType))
	flattenedScopes[schemas.VariableScopeFieldNames.Roles] = util.Ternary(variableScopes.Roles != nil, util.FlattenStringList(variableScopes.Roles), types.ListNull(types.StringType))
	flattenedScopes[schemas.VariableScopeFieldNames.TenantTags] = util.Ternary(variableScopes.TenantTags != nil, util.FlattenStringList(variableScopes.TenantTags), types.ListNull(types.StringType))

	return types.ObjectValueMust(
		variableScopeObjectType(),
		flattenedScopes,
	)
}

func expandVariableScopes(flattenedVariableScopes types.List) variables.VariableScope {
	if flattenedVariableScopes.IsNull() {
		return variables.VariableScope{}
	}

	obj := flattenedVariableScopes.Elements()[0].(types.Object)
	attrs := obj.Attributes()
	scopes := variables.VariableScope{}

	scopes.Actions = util.ExpandStringList(attrs[schemas.VariableScopeFieldNames.Actions].(types.List))
	scopes.Channels = util.ExpandStringList(attrs[schemas.VariableScopeFieldNames.Channels].(types.List))
	scopes.Environments = util.ExpandStringList(attrs[schemas.VariableScopeFieldNames.Environments].(types.List))
	scopes.Machines = util.ExpandStringList(attrs[schemas.VariableScopeFieldNames.Machines].(types.List))
	scopes.ProcessOwners = util.ExpandStringList(attrs[schemas.VariableScopeFieldNames.Processes].(types.List))
	scopes.Roles = util.ExpandStringList(attrs[schemas.VariableScopeFieldNames.Roles].(types.List))
	scopes.TenantTags = util.ExpandStringList(attrs[schemas.VariableScopeFieldNames.TenantTags].(types.List))

	return scopes
}
